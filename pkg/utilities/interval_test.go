package utilities

import (
	"testing"
	"time"
)

func TestIntervalDuration(t *testing.T) {
	tests := []struct {
		interval string
		want    time.Duration
	}{
		{"15m", 15 * time.Minute},
		{"30m", 30 * time.Minute},
		{"1h", time.Hour},
		{"4h", 4 * time.Hour},
		{"1d", 24 * time.Hour},
		{"unknown", time.Hour},
		{"", time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.interval, func(t *testing.T) {
			got := IntervalDuration(tt.interval)
			if got != tt.want {
				t.Errorf("IntervalDuration(%q) = %v, want %v", tt.interval, got, tt.want)
			}
		})
	}
}
