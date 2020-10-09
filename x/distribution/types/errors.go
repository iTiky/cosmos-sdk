package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/distribution module sentinel errors
var (
	ErrInternal                       = sdkerrors.Register(ModuleName, 0, "internal error")
	ErrEmptyDelegatorAddr             = sdkerrors.Register(ModuleName, 1, "delegator address is empty")
	ErrEmptyWithdrawAddr              = sdkerrors.Register(ModuleName, 2, "withdraw address is empty")
	ErrEmptyValidatorAddr             = sdkerrors.Register(ModuleName, 3, "validator address is empty")
	ErrEmptyDelegationDistInfo        = sdkerrors.Register(ModuleName, 4, "no delegation distribution info")
	ErrNoValidatorDistInfo            = sdkerrors.Register(ModuleName, 5, "no validator distribution info")
	ErrNoValidatorCommission          = sdkerrors.Register(ModuleName, 6, "no validator commission to withdraw")
	ErrSetWithdrawAddrDisabled        = sdkerrors.Register(ModuleName, 7, "set withdraw address disabled")
	ErrBadDistribution                = sdkerrors.Register(ModuleName, 8, "pool does not have sufficient coins to distribute")
	ErrInvalidProposalAmount          = sdkerrors.Register(ModuleName, 9, "invalid amount")
	ErrEmptyProposalRecipient         = sdkerrors.Register(ModuleName, 10, "invalid proposal recipient")
	ErrNoValidatorExists              = sdkerrors.Register(ModuleName, 11, "validator does not exist")
	ErrNoDelegationExists             = sdkerrors.Register(ModuleName, 12, "delegation does not exist")
	ErrWrongFoundationAllocationRatio = sdkerrors.Register(ModuleName, 13, "foundation allocation ratio is wrong")
	ErrExceededTimeLimit              = sdkerrors.Register(ModuleName, 14, "exceeded time limit")
	ErrWithdrawLocked                 = sdkerrors.Register(ModuleName, 15, "rewards withdraw is locked")
	ErrInvalidLockOperation           = sdkerrors.Register(ModuleName, 16, "invalid rewards lock operation")
)
