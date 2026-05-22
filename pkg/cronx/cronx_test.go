package cronx

import (
	"testing"
	"time"
)

// TestCron_Parse_FiveField_Description — C6: 5-field unix expression parses + describes.
func TestCron_Parse_FiveField_Description(t *testing.T) {
	t.Parallel()
	res, err := Parse("0 12 * * *", ParseOptions{})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(res.Diagnostics) != 0 {
		t.Fatalf("diags: %+v", res.Diagnostics)
	}
	if len(res.Fields) != 5 || res.Fields[1].Value != "12" {
		t.Fatalf("fields: %+v", res.Fields)
	}
}

// TestCron_Quartz_Seconds — C6: quartz flavor accepts a 6-field expression with seconds.
func TestCron_Quartz_Seconds(t *testing.T) {
	t.Parallel()
	res, err := Parse("*/5 * * * * *", ParseOptions{Flavor: "quartz"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(res.Diagnostics) != 0 {
		t.Fatalf("diags: %+v", res.Diagnostics)
	}
	if res.Fields[0].Name != "second" {
		t.Fatalf("expected leading second field: %+v", res.Fields)
	}
}

// TestCron_NextRuns_5min — C6: */5 * * * * yields five 5-minute spaced runs.
func TestCron_NextRuns_5min(t *testing.T) {
	t.Parallel()
	from := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	res, err := NextRuns("*/5 * * * *", from, 3, time.UTC)
	if err != nil {
		t.Fatalf("NextRuns: %v", err)
	}
	if len(res.Runs) != 3 {
		t.Fatalf("got %d runs", len(res.Runs))
	}
	if res.Runs[1].Sub(res.Runs[0]) != 5*time.Minute {
		t.Fatalf("spacing: %v", res.Runs)
	}
}

// TestCron_AWS_FlavorRejectsDayOfWeekName — C6: AWS does not accept "MON".
func TestCron_AWS_FlavorRejectsDayOfWeekName(t *testing.T) {
	t.Parallel()
	res, err := Parse("0 12 * * MON", ParseOptions{Flavor: "aws"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(res.Diagnostics) == 0 || res.Diagnostics[0].Code != "CRON.AWS_REJECTS_DOW_NAME" {
		t.Fatalf("expected DOW-name diag, got %+v", res.Diagnostics)
	}
}

// TestCron_Invalid — C6: invalid expression surfaces CRON.INVALID.
func TestCron_Invalid(t *testing.T) {
	t.Parallel()
	res, err := Parse("99 99 * * *", ParseOptions{})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if res.Diagnostics[0].Code != "CRON.INVALID" {
		t.Fatalf("code = %q", res.Diagnostics[0].Code)
	}
}
