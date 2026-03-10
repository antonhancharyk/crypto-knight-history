package utilities

import "time"

// IntervalDuration returns the duration for a Binance kline interval string.
func IntervalDuration(interval string) time.Duration {
	switch interval {
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}
