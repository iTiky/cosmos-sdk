//nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Verify interface at compile time
var _, _, _ sdk.Msg = &MsgSetWithdrawAddress{}, &MsgWithdrawDelegatorReward{}, &MsgWithdrawValidatorCommission{}

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgSetWithdrawAddress struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	WithdrawAddress  sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
}

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		DelegatorAddress: delAddr,
		WithdrawAddress:  withdrawAddr,
	}
}

func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }
func (msg MsgSetWithdrawAddress) Type() string  { return "set_withdraw_address" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.WithdrawAddress.Empty() {
		return ErrEmptyWithdrawAddr
	}

	return nil
}

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegatorReward struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) MsgWithdrawDelegatorReward {
	return MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
}

func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }
func (msg MsgWithdrawDelegatorReward) Type() string  { return "withdraw_delegator_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	return nil
}

// msg struct for validator withdraw
type MsgWithdrawValidatorCommission struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func NewMsgWithdrawValidatorCommission(valAddr sdk.ValAddress) MsgWithdrawValidatorCommission {
	return MsgWithdrawValidatorCommission{
		ValidatorAddress: valAddr,
	}
}

func (msg MsgWithdrawValidatorCommission) Route() string { return ModuleName }
func (msg MsgWithdrawValidatorCommission) Type() string  { return "withdraw_validator_commission" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorCommission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawValidatorCommission) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	return nil
}

const TypeMsgFundPublicTreasuryPool = "fund_public_treasury_pool"

// MsgFundPublicTreasuryPool defines a Msg type that allows an account to directly fund the public treasury pool.
type MsgFundPublicTreasuryPool struct {
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"`
}

// NewMsgFundPublicTreasuryPool returns a new MsgFundPublicTreasuryPool with a sender and a funding amount.
func NewMsgFundPublicTreasuryPool(amount sdk.Coins, depositor sdk.AccAddress) MsgFundPublicTreasuryPool {
	return MsgFundPublicTreasuryPool{
		Amount:    amount,
		Depositor: depositor,
	}
}

// Route returns the MsgFundPublicTreasuryPool message route.
func (msg MsgFundPublicTreasuryPool) Route() string { return ModuleName }

// Type returns the MsgFundPublicTreasuryPool message type.
func (msg MsgFundPublicTreasuryPool) Type() string { return TypeMsgFundPublicTreasuryPool }

// GetSigners returns the signer addresses that are expected to sign the result of GetSignBytes.
func (msg MsgFundPublicTreasuryPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

// GetSignBytes returns the raw bytes for a MsgFundPublicTreasuryPool message that the expected signer needs to sign.
func (msg MsgFundPublicTreasuryPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic MsgFundPublicTreasuryPool message validation.
func (msg MsgFundPublicTreasuryPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Depositor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor.String())
	}

	return nil
}

const TypeMsgWithdrawFoundationPool = "withdraw_foundation_pool"

// MsgWithdrawFoundationPool defines a Msg type that allows transfer from the FoundationPool
// to target account waller / module account.
// Nominee limits number of accounts allowed to do the transfer.
type MsgWithdrawFoundationPool struct {
	NomineeAddr   sdk.AccAddress `json:"nominee_address" yaml:"nominee_address"`
	RecipientAddr sdk.AccAddress `json:"recipient_address" yaml:"recipient_address"`
	RecipientPool RewardPoolName `json:"recipient_pool" yaml:"recipient_pool"`
	Amount        sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewMsgFundPublicTreasuryPool returns a new MsgFundPublicTreasuryPool with a sender and a funding amount.
func NewMsgWithdrawFoundationPool(nominee, recipient sdk.AccAddress, pool RewardPoolName, amount sdk.Coins) MsgWithdrawFoundationPool {
	return MsgWithdrawFoundationPool{
		NomineeAddr:   nominee,
		RecipientAddr: recipient,
		RecipientPool: pool,
		Amount:        amount,
	}
}

// Route returns the MsgWithdrawFoundationPool message route.
func (msg MsgWithdrawFoundationPool) Route() string { return ModuleName }

// Type returns the MsgWithdrawFoundationPool message type.
func (msg MsgWithdrawFoundationPool) Type() string { return TypeMsgWithdrawFoundationPool }

// GetSigners returns the signer addresses that are expected to sign the result of GetSignBytes.
func (msg MsgWithdrawFoundationPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.NomineeAddr}
}

// GetSignBytes returns the raw bytes for a MsgWithdrawFoundationPool message that the expected signer needs to sign.
func (msg MsgWithdrawFoundationPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic MsgWithdrawFoundationPool message validation.
func (msg MsgWithdrawFoundationPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.NomineeAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "nominee address: empty")
	}
	if msg.RecipientAddr.Empty() && msg.RecipientPool == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "recipient: not defined")
	}
	if !msg.RecipientAddr.Empty() && msg.RecipientPool != "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "recipient: only one is allowed (wallet/pool): %s / %s", msg.RecipientAddr, msg.RecipientPool)
	}
	if msg.RecipientPool != "" {
		if !msg.RecipientPool.IsValid() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "recipient: invalid pool name: %s", msg.RecipientPool)
		}
		if msg.RecipientPool == FoundationPoolName {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "recipient: target pool can not be %s", FoundationPoolName)
		}
	}

	return nil
}

const TypeMsgSetFoundationAllocationRatio = "set_foundation_allocation_ratio"

// MsgSetFoundationAllocationRatio - high level transaction of the coin module
type MsgSetFoundationAllocationRatio struct {
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Ratio       sdk.Dec        `json:"ratio" yaml:"ratio"`
}

// NewMsgSetFoundationAllocationRatio - construct msg for change FoundationAllocationRatio.
func NewMsgSetFoundationAllocationRatio(fromAddr sdk.AccAddress, ratio sdk.Dec) MsgSetFoundationAllocationRatio {
	return MsgSetFoundationAllocationRatio{FromAddress: fromAddr, Ratio: ratio}
}

// Route Implements Msg.
func (msg MsgSetFoundationAllocationRatio) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetFoundationAllocationRatio) Type() string { return TypeMsgSetFoundationAllocationRatio }

// ValidateBasic Implements Msg.
func (msg MsgSetFoundationAllocationRatio) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.Ratio.GT(sdk.NewDec(ChangeFoundationAllocationRatioMaxValue)) {
		return sdkerrors.Wrap(ErrWrongFoundationAllocationRatio, "ratio is greater than the maximum value for FoundationAllocationRatio")
	}
	if msg.Ratio.LT(sdk.NewDec(ChangeFoundationAllocationRatioMinValue)) {
		return sdkerrors.Wrap(ErrWrongFoundationAllocationRatio, "ratio is lower than the maximum value for FoundationAllocationRatio")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetFoundationAllocationRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetFoundationAllocationRatio) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
