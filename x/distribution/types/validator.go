package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorHistoricalRewards keeps historical rewards for a validator.
// Height is implicit within the store key.
// ReferenceCount =
//    number of outstanding delegations which ended the associated period (and might need to read that record)
//  + number of slashes which ended the associated period (and might need to read that record)
//  + one per validator for the zeroeth period, set on initialization
type ValidatorHistoricalRewards struct {
	// Sum from the zeroeth period until this period of rewards / tokens, per the spec
	CumulativeRewardRatio sdk.DecCoins `json:"cumulative_reward_ratio" yaml:"cumulative_reward_ratio"`
	// Indicates the number of objects which might need to reference this historical entry at any point
	ReferenceCount uint16 `json:"reference_count" yaml:"reference_count"`
}

// NewValidatorHistoricalRewards creates a new ValidatorHistoricalRewards.
func NewValidatorHistoricalRewards(cumulativeRewardRatio sdk.DecCoins, referenceCount uint16) ValidatorHistoricalRewards {
	return ValidatorHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

// ValidatorCurrentRewards keeps current rewards and current period for a validator.
// Kept as a running counter and incremented each block as long as the validator's tokens remain constant.
type ValidatorCurrentRewards struct {
	// Current rewards
	Rewards sdk.DecCoins `json:"rewards" yaml:"rewards"`
	// Current period
	Period uint64 `json:"period" yaml:"period"`
}

// NewValidatorCurrentRewards creates a new ValidatorCurrentRewards.
func NewValidatorCurrentRewards(rewards sdk.DecCoins, period uint64) ValidatorCurrentRewards {
	return ValidatorCurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// ValidatorAccumulatedCommission keeps accumulated commission for a validator.
// Kept as a running counter, can be withdrawn at any time.
type ValidatorAccumulatedCommission = sdk.DecCoins

// InitialValidatorAccumulatedCommission returns the initial accumulated commission (zero).
func InitialValidatorAccumulatedCommission() ValidatorAccumulatedCommission {
	return ValidatorAccumulatedCommission{}
}

// ValidatorSlashEvent needed to calculate appropriate amounts of staking token
// for delegations which withdraw after a slash has occurred.
// Height is implicit within the store key.
type ValidatorSlashEvent struct {
	// Period when the slash occurred
	ValidatorPeriod uint64 `json:"validator_period" yaml:"validator_period"`
	// Slash fraction
	Fraction sdk.Dec `json:"fraction" yaml:"fraction"`
}

// NewValidatorSlashEvent creates a new ValidatorSlashEvent.
func NewValidatorSlashEvent(validatorPeriod uint64, fraction sdk.Dec) ValidatorSlashEvent {
	return ValidatorSlashEvent{
		ValidatorPeriod: validatorPeriod,
		Fraction:        fraction,
	}
}

func (vs ValidatorSlashEvent) String() string {
	return fmt.Sprintf(`Period:   %d
Fraction: %s`, vs.ValidatorPeriod, vs.Fraction)
}

// ValidatorSlashEvents is a collection of ValidatorSlashEvent.
type ValidatorSlashEvents []ValidatorSlashEvent

func (vs ValidatorSlashEvents) String() string {
	out := "Validator Slash Events:\n"
	for i, sl := range vs {
		out += fmt.Sprintf(`  Slash %d:
    Period:   %d
    Fraction: %s
`, i, sl.ValidatorPeriod, sl.Fraction)
	}
	return strings.TrimSpace(out)
}

// ValidatorOutstandingRewards keeps outstanding (un-withdrawn) rewards for a validator.
// It is inexpensive to track, allows simple sanity checks.
type ValidatorOutstandingRewards = sdk.DecCoins

// ValidatorLockedRewards contains locked rewards data.
type ValidatorLockedRewards struct {
	// Locked shares to all shares relation (zero if there is no locking)
	LockedRatio sdk.Dec `json:"locked_ratio" yaml:"locked_ratio"`
}

// NewValidatorLockedRewards creates a new ValidatorLockedRewards.
func NewValidatorLockedRewards(lockedRatio sdk.Dec) ValidatorLockedRewards {
	return ValidatorLockedRewards{
		LockedRatio: lockedRatio,
	}
}
