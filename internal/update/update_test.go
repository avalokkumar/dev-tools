package update

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

type fakeSource struct {
	current, latest string
}

func (f fakeSource) Name() string { return "fake" }
func (f fakeSource) Latest(_ context.Context, _ string) (Plan, error) {
	return Plan{
		Current:        f.current,
		Latest:         f.latest,
		UpdateRequired: f.latest != f.current,
		Source:         "fake",
	}, nil
}

// TestUpdate_DryRun_LogsPlanned — D1: --dry-run prints "would update".
func TestUpdate_DryRun_LogsPlanned(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	plan, err := Apply(context.Background(), ApplyOptions{
		Source: fakeSource{current: "0.0.0-dev", latest: "0.1.0"},
		DryRun: true,
		Out:    &out,
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !plan.UpdateRequired {
		t.Fatalf("UpdateRequired = false")
	}
	if !strings.Contains(out.String(), "DRY RUN") {
		t.Fatalf("missing DRY RUN: %q", out.String())
	}
	if !strings.Contains(out.String(), "0.0.0-dev → 0.1.0") {
		t.Fatalf("missing version transition: %q", out.String())
	}
}

// TestUpdate_StubSource_NoChannel — D1: default source reports no channel.
func TestUpdate_StubSource_NoChannel(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	plan, err := Apply(context.Background(), ApplyOptions{Out: &out})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if plan.UpdateRequired {
		t.Fatalf("stub should never claim update required")
	}
	if !strings.Contains(out.String(), "no release channel configured") {
		t.Fatalf("missing note: %q", out.String())
	}
}

// TestUpdate_AlreadyCurrent — D1: when latest == current, friendly message.
func TestUpdate_AlreadyCurrent(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	_, err := Apply(context.Background(), ApplyOptions{
		Source: fakeSource{current: "1.0.0", latest: "1.0.0"},
		Out:    &out,
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !strings.Contains(out.String(), "already current") {
		t.Fatalf("expected already-current message: %q", out.String())
	}
}
