package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec concrete distribution types on amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "cosmos-sdk/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorCommission{}, "cosmos-sdk/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "cosmos-sdk/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(PublicTreasuryPoolSpendProposal{}, "cosmos-sdk/PublicTreasuryPoolSpendProposal", nil)
	cdc.RegisterConcrete(MsgFundPublicTreasuryPool{}, "cosmos-sdk/MsgFundPublicTreasuryPool", nil)
	cdc.RegisterConcrete(MsgWithdrawFoundationPool{}, "cosmos-sdk/MsgWithdrawFoundationPool", nil)
	cdc.RegisterConcrete(MsgSetFoundationAllocationRatio{}, "cosmos-sdk/MsgSetFoundationAllocationRatio", nil)
	cdc.RegisterConcrete(MsgLockValidatorRewards{}, "cosmos-sdk/MsgLockValidatorRewards", nil)
	cdc.RegisterConcrete(MsgDisableLockedRewardsAutoRenewal{}, "cosmos-sdk/MsgDisableLockedRewardsAutoRenewal", nil)
	cdc.RegisterConcrete(MsgSetStakingTotalSupplyShift{}, "cosmos-sdk/MsgSetStakingTotalSupplyShift", nil)
}

// ModuleCdc is a generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
