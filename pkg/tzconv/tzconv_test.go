package tzconv

import (
	"testing"
	"time"
)

// TestTz_Convert_NYToTokyo — C8: NY noon → Tokyo gives expected wall hour.
func TestTz_Convert_NYToTokyo(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	res, err := Convert(t0, "America/New_York", "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if len(res.Diagnostics) != 0 {
		t.Fatalf("diags: %+v", res.Diagnostics)
	}
	// In May 2026 NY = EDT (UTC-4); Tokyo = JST (UTC+9). Diff = 13 hours.
	wantHour := 1 // 12 + 13 = 25 → 1 next day
	if res.Converted.Hour() != wantHour {
		t.Fatalf("Tokyo hour = %d, want %d", res.Converted.Hour(), wantHour)
	}
}

// TestTz_Convert_DSTGap_NoteSet — C8: 02:30 on US spring-forward Sunday is a gap.
func TestTz_Convert_DSTGap_NoteSet(t *testing.T) {
	t.Parallel()
	// 2026-03-08 02:30 in America/New_York is in the DST gap (skipped).
	t0 := time.Date(2026, 3, 8, 2, 30, 0, 0, time.UTC)
	res, err := Convert(t0, "America/New_York", "UTC")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.DSTNote != "gap" {
		t.Fatalf("expected gap, got %q", res.DSTNote)
	}
}

// TestTz_ListZones_FilterEurope — C8: filter narrows results.
func TestTz_ListZones_FilterEurope(t *testing.T) {
	t.Parallel()
	zs, err := ListZones("europe")
	if err != nil {
		t.Fatalf("ListZones: %v", err)
	}
	if len(zs) < 5 {
		t.Fatalf("expected several Europe zones, got %d", len(zs))
	}
	for _, z := range zs {
		if z.Name == "" {
			t.Fatalf("empty name")
		}
	}
}

// TestTz_Convert_UnknownZone — C8: bogus zone surfaces diagnostic.
func TestTz_Convert_UnknownZone(t *testing.T) {
	t.Parallel()
	res, err := Convert(time.Now(), "Mars/Olympus_Mons", "UTC")
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if res.Diagnostics[0].Code != "TZ.UNKNOWN_FROM" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}
