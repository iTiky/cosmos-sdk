package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
)

func NewHandler(k keeper.Keeper, mk mint.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)

		case types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, msg, k)

		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)

		case types.MsgFundPublicTreasuryPool:
			return handleMsgFundPublicTreasuryPool(ctx, msg, k)

		case types.MsgWithdrawFoundationPool:
			return handleMsgWithdrawFoundationPool(ctx, msg, k)

		case types.MsgSetFoundationAllocationRatio:
			return handleMsgSetFoundationAllocationRatio(ctx, msg, k, mk)

		case types.MsgLockValidatorRewards:
			return handleMsgLockValidatorRewards(ctx, msg, k)

		case types.MsgDisableLockedRewardsAutoRenewal:
			return handleMsgDisableLockedRewardsAutoRenewal(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)

func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.Keeper) (*sdk.Result, error) {
	err := k.SetWithdrawAddr(ctx, msg.DelegatorAddress, msg.WithdrawAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.MsgWithdrawDelegatorReward, k keeper.Keeper) (*sdk.Result, error) {
	_, err := k.WithdrawDelegationRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.MsgWithdrawValidatorCommission, k keeper.Keeper) (*sdk.Result, error) {
	_, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgFundPublicTreasuryPool(ctx sdk.Context, msg types.MsgFundPublicTreasuryPool, k keeper.Keeper) (*sdk.Result, error) {
	if err := k.FundPublicTreasuryPool(ctx, msg.Amount, msg.Depositor); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawFoundationPool(ctx sdk.Context, msg types.MsgWithdrawFoundationPool, k keeper.Keeper) (*sdk.Result, error) {
	params := k.GetParams(ctx)

	isNominee := false
	for _, nominee := range params.FoundationNominees {
		if msg.NomineeAddr.Equals(nominee) {
			isNominee = true
			break
		}
	}
	if !isNominee {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "signer is not a nominee account: %s", msg.NomineeAddr)
	}

	if !msg.RecipientAddr.Empty() {
		if err := k.DistributeFromFoundationPoolToWallet(ctx, msg.Amount, msg.RecipientAddr); err != nil {
			return nil, err
		}
		k.Logger(ctx).Info(fmt.Sprintf("transferred %s from the foundation pool to recipient %s (authorized by %s)", msg.Amount, msg.RecipientAddr, msg.NomineeAddr))
		return nil, nil
	}

	if err := k.DistributeFromFoundationPoolToPool(ctx, msg.Amount, msg.RecipientPool); err != nil {
		return nil, err
	}
	k.Logger(ctx).Info(fmt.Sprintf("transferred %s from the foundation pool to recipient %s (authorized by %s)", msg.Amount, msg.RecipientPool, msg.NomineeAddr))

	return nil, nil
}

func handleMsgLockValidatorRewards(ctx sdk.Context, msg types.MsgLockValidatorRewards, k keeper.Keeper) (*sdk.Result, error) {
	unlocksAt, err := k.LockValidatorRewards(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyLockedRewardsState, types.LockedRewardsStateLocked),
		),
	)

	k.Logger(ctx).Info(fmt.Sprintf("Validator %s rewards were locked: unlocksAt: %v", msg.ValidatorAddress, unlocksAt))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetFoundationAllocationRatio(
	ctx sdk.Context,
	msg types.MsgSetFoundationAllocationRatio,
	k keeper.Keeper,
	mk mint.Keeper,
) (*sdk.Result, error) {
	hasPermission := false
	for _, nominee := range k.GetParams(ctx).FoundationNominees {
		if nominee.Equals(msg.FromAddress) {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "operation is allowed only for foundation nominee")
	}

	abpy, err := mk.GetAvgBlocksPerYear(ctx)
	if err != nil {
		return nil, err
	}

	chainAge := sdk.NewInt(ctx.BlockHeight()).QuoRaw(int64(abpy))

	if chainAge.GTE(sdk.NewInt(ChangeFoundationAllocationRatioTTL)) {
		return nil, sdkerrors.Wrapf(ErrExceededTimeLimit, "is not allowed to change after %d year", ChangeFoundationAllocationRatioTTL)
	}

	params := mk.GetParams(ctx)
	params.FoundationAllocationRatio = msg.Ratio
	mk.SetParams(ctx, params)

	return &sdk.Result{}, nil
}

func handleMsgDisableLockedRewardsAutoRenewal(ctx sdk.Context, msg types.MsgDisableLockedRewardsAutoRenewal, k keeper.Keeper) (*sdk.Result, error) {
	if err := k.DisableLockedRewardsAutoRenewal(ctx, msg.ValidatorAddress); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyLockedRewardsState, types.LockedRewardsStateAutoRenewDisabled),
		),
	)

	k.Logger(ctx).Info(fmt.Sprintf("Validator %s rewards lock auto-renewal disabled", msg.ValidatorAddress))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case types.PublicTreasuryPoolSpendProposal:
			return keeper.HandlePublicTreasuryPoolSpendProposal(ctx, k, c)
		case types.TaxParamsUpdateProposal:
			return keeper.HandleTaxParamsUpdateProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
