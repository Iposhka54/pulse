package utils

import "time"

func DurationToSeconds(duration time.Duration) int64 {
	return int64(duration / time.Second)
}

func DurationToSecondsFloat(duration time.Duration) float64 {
	return float64(DurationToSeconds(duration))
}
