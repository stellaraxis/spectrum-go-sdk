package timefmt

import (
	"testing"
	"time"
)

func TestFormatUsesFixedLocalTimezone(t *testing.T) {
	input := time.Date(2026, 4, 20, 7, 8, 9, 123456000, time.UTC)

	got := Format(input)
	want := "2026-04-20 15:08:09.123"

	if got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}
