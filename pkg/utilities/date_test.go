package utilities

import (
	"testing"
	"time"
)

func TestFromUnixToUTC(t *testing.T) {
	// 1000 ms = 1 second
	got := FromUnixToUTC(1000)
	want := time.Unix(1, 0).UTC()
	if !got.Equal(want) {
		t.Errorf("FromUnixToUTC(1000) = %v, want %v", got, want)
	}

	// 0
	got = FromUnixToUTC(0)
	want = time.Unix(0, 0).UTC()
	if !got.Equal(want) {
		t.Errorf("FromUnixToUTC(0) = %v, want %v", got, want)
	}

	// 1234567890123 ms = 1234567890.123 s
	got = FromUnixToUTC(1234567890123)
	want = time.Unix(1234567890, 123000000).UTC()
	if !got.Equal(want) {
		t.Errorf("FromUnixToUTC(1234567890123) = %v, want %v", got, want)
	}
}
