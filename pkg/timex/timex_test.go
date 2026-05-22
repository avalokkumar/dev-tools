package timex

import (
	"strings"
	"testing"
	"time"
)

func TestConvert_EpochS_Auto(t *testing.T) {
	t.Parallel()
	r, _ := Convert(ConvertOptions{Input: "1715184000"})
	if r.EpochS != 1715184000 {
		t.Fatalf("EpochS = %d", r.EpochS)
	}
	if !strings.Contains(r.RFC3339, "2024-05-08") {
		t.Fatalf("RFC3339 = %s", r.RFC3339)
	}
}

func TestConvert_EpochMS_Auto(t *testing.T) {
	t.Parallel()
	r, _ := Convert(ConvertOptions{Input: "1715184000123"})
	if r.EpochMS != 1715184000123 {
		t.Fatalf("EpochMS = %d", r.EpochMS)
	}
}

func TestConvert_RFC3339Input(t *testing.T) {
	t.Parallel()
	r, _ := Convert(ConvertOptions{Input: "2026-05-09T12:00:00Z"})
	if r.EpochS != time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC).Unix() {
		t.Fatalf("EpochS = %d", r.EpochS)
	}
}

func TestConvert_UnknownTZ(t *testing.T) {
	t.Parallel()
	r, _ := Convert(ConvertOptions{Input: "1", TZ: "Mars/Olympus"})
	if r.Diagnostics[0].Code != "TIME.CONVERT.UNKNOWN_TZ" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestRelative_Past(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	earlier := now.Add(-5 * time.Minute)
	r, _ := Relative(now, earlier)
	if !strings.Contains(r.Phrase, "ago") {
		t.Fatalf("phrase = %q", r.Phrase)
	}
}

func TestRelative_Future(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	later := now.Add(3 * 24 * time.Hour)
	r, _ := Relative(now, later)
	if !strings.Contains(r.Phrase, "from now") {
		t.Fatalf("phrase = %q", r.Phrase)
	}
	if !strings.Contains(r.Phrase, "day") {
		t.Fatalf("phrase = %q", r.Phrase)
	}
}

func TestDuration_Parse(t *testing.T) {
	t.Parallel()
	r, _ := Duration("1h30m")
	if r.Hours != 1 || r.Minutes != 30 {
		t.Fatalf("got %+v", r)
	}
}

func TestDuration_PlainInteger(t *testing.T) {
	t.Parallel()
	r, _ := Duration("90")
	if r.Minutes != 1 || r.Seconds != 30 {
		t.Fatalf("got %+v", r)
	}
}

func TestDuration_Invalid(t *testing.T) {
	t.Parallel()
	r, _ := Duration("bogus")
	if r.Diagnostics[0].Code != "TIME.DURATION.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}
