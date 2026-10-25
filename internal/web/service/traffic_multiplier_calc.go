package service

import (
    "errors"
    "math"
)

var ErrInvalidTrafficMultiplier = errors.New("invalid traffic multiplier")

// WeightedTraffic computes the billable bytes for an inbound.
// A zero multiplier represents the legacy default of 1.
func WeightedTraffic(bytes int64, multiplier float64) (int64, error) {
    if bytes < 0 {
        return 0, ErrInvalidTrafficMultiplier
    }
    if multiplier == 0 {
        multiplier = 1
    }
    if math.IsNaN(multiplier) || math.IsInf(multiplier, 0) || multiplier < 0 {
        return 0, ErrInvalidTrafficMultiplier
    }

    result := float64(bytes) * multiplier
    if math.IsInf(result, 0) || result > float64(math.MaxInt64) {
        return 0, ErrInvalidTrafficMultiplier
    }

    return int64(math.Round(result)), nil
}
