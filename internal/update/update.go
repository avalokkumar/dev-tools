// Package update implements the `devforge update` subcommand.
//
// The actual binary swap is delegated to a pluggable Source. The default
// source is a stub that reports "no release channel configured" — real
// channels (GitHub releases, internal mirror) plug in via SetSource. The
// command itself is exercised by tests via the dry-run path.
package update

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/devforge/devforge/internal/version"
)

// Plan describes what an update would do without performing it.
type Plan struct {
	Current        string `json:"current"`
	Latest         string `json:"latest"`
	UpdateRequired bool   `json:"updateRequired"`
	Source         string `json:"source"`
	Asset          string `json:"asset,omitempty"`
	Note           string `json:"note,omitempty"`
}

// Source resolves the latest available release. Implementations must be
// safe for concurrent use; they are called once per invocation.
type Source interface {
	Name() string
	Latest(ctx context.Context, currentVersion string) (Plan, error)
}

// stubSource is the default. It always returns "no channel configured".
type stubSource struct{}

func (stubSource) Name() string { return "stub" }
func (stubSource) Latest(_ context.Context, current string) (Plan, error) {
	return Plan{
		Current:        current,
		Latest:         current,
		UpdateRequired: false,
		Source:         "stub",
		Note:           "no release channel configured; pass --source github:<owner>/<repo> when one ships",
	}, nil
}

var defaultSource Source = stubSource{}

// SetSource overrides the package-level default. Production wiring may use
// this in main() to attach a real GitHub release source once available.
func SetSource(s Source) { defaultSource = s }

// CheckOptions controls Check.
type CheckOptions struct {
	Source Source // nil = use defaultSource
}

// Check resolves the latest release without applying it.
func Check(ctx context.Context, opts CheckOptions) (Plan, error) {
	src := opts.Source
	if src == nil {
		src = defaultSource
	}
	return src.Latest(ctx, version.Get().Version)
}

// Apply runs the update plan. The stub source has no asset; real sources
// download + atomically swap the binary.
type ApplyOptions struct {
	Source Source
	DryRun bool
	Out    io.Writer
}

// Apply (or simulate) an update. Returns the Plan describing what happened.
func Apply(ctx context.Context, opts ApplyOptions) (Plan, error) {
	if opts.Out == nil {
		opts.Out = io.Discard
	}
	plan, err := Check(ctx, CheckOptions{Source: opts.Source})
	if err != nil {
		return plan, err
	}
	plan.Asset = fmt.Sprintf("devforge_%s_%s", runtime.GOOS, runtime.GOARCH)

	if !plan.UpdateRequired {
		fmt.Fprintf(opts.Out, "devforge %s is already current (source=%s).\n",
			plan.Current, plan.Source)
		if plan.Note != "" {
			fmt.Fprintf(opts.Out, "note: %s\n", plan.Note)
		}
		return plan, nil
	}

	if opts.DryRun {
		fmt.Fprintf(opts.Out,
			"DRY RUN: would update %s → %s from %s (asset=%s).\n",
			plan.Current, plan.Latest, plan.Source, plan.Asset)
		return plan, nil
	}

	// Real apply path is delegated to the Source implementation in Phase D+.
	// The stub never reaches here.
	fmt.Fprintf(opts.Out, "update path for source %s is not yet implemented.\n", plan.Source)
	return plan, nil
}
