// Package cronx parses cron expressions and computes upcoming run times.
//
// External API:
//
//	Parse(expr string, opts ParseOptions) (ParseResult, error)
//	NextRuns(expr string, from time.Time, n int, tz *time.Location) (NextRunsResult, error)
//
// Supported flavors:
//   - "unix"   (default): minute hour dom month dow         — robfig parser
//   - "quartz": second minute hour dom month dow            — six-field parser
//   - "aws":    minute hour dom month dow year              — strict, declines DOW names
package cronx

import (
	"fmt"
	"strings"
	"time"

	cronpkg "github.com/robfig/cron/v3"

	"github.com/devforge/devforge/pkg/engine"
)

// ParseOptions tunes Parse.
type ParseOptions struct {
	// Flavor: "unix" (default), "quartz", "aws".
	Flavor string `json:"flavor,omitempty"`
}

// Field is one component of the cron expression.
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ParseResult is the success return.
type ParseResult struct {
	Description string              `json:"description"`
	Fields      []Field             `json:"fields"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Parse validates and decomposes a cron expression.
func Parse(expr string, opts ParseOptions) (ParseResult, error) {
	flavor := opts.Flavor
	if flavor == "" {
		flavor = "unix"
	}
	parts := strings.Fields(expr)
	switch flavor {
	case "unix":
		if len(parts) != 5 {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.WRONG_FIELD_COUNT",
				Message: fmt.Sprintf("unix flavor expects 5 fields, got %d", len(parts)),
				Severity: engine.SevError,
			}}}, nil
		}
		if _, err := cronpkg.ParseStandard(expr); err != nil {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.INVALID", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		return ParseResult{
			Description: describeUnix(parts),
			Fields: []Field{
				{Name: "minute", Value: parts[0]},
				{Name: "hour", Value: parts[1]},
				{Name: "dayOfMonth", Value: parts[2]},
				{Name: "month", Value: parts[3]},
				{Name: "dayOfWeek", Value: parts[4]},
			},
		}, nil
	case "quartz":
		if len(parts) != 6 {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.WRONG_FIELD_COUNT",
				Message: fmt.Sprintf("quartz flavor expects 6 fields, got %d", len(parts)),
				Severity: engine.SevError,
			}}}, nil
		}
		// Use a parser with seconds.
		p := cronpkg.NewParser(cronpkg.SecondOptional | cronpkg.Minute | cronpkg.Hour |
			cronpkg.Dom | cronpkg.Month | cronpkg.Dow)
		if _, err := p.Parse(expr); err != nil {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.INVALID", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		return ParseResult{
			Description: "every " + parts[0] + "s within " + describeUnix(parts[1:]),
			Fields: []Field{
				{Name: "second", Value: parts[0]},
				{Name: "minute", Value: parts[1]},
				{Name: "hour", Value: parts[2]},
				{Name: "dayOfMonth", Value: parts[3]},
				{Name: "month", Value: parts[4]},
				{Name: "dayOfWeek", Value: parts[5]},
			},
		}, nil
	case "aws":
		if len(parts) != 5 && len(parts) != 6 {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.WRONG_FIELD_COUNT",
				Message: fmt.Sprintf("aws flavor expects 5 or 6 fields, got %d", len(parts)),
				Severity: engine.SevError,
			}}}, nil
		}
		// AWS cron does not allow day-of-week alpha names.
		if len(parts) >= 5 && containsAlpha(parts[4]) {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.AWS_REJECTS_DOW_NAME",
				Message: "AWS flavor does not accept day-of-week names; use numbers (1-7)",
				Severity: engine.SevError,
			}}}, nil
		}
		// Otherwise piggyback on unix validation for the first 5 fields.
		joined := strings.Join(parts[:5], " ")
		if _, err := cronpkg.ParseStandard(joined); err != nil {
			return ParseResult{Diagnostics: []engine.Diagnostic{{
				Code: "CRON.INVALID", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		return ParseResult{
			Description: describeUnix(parts[:5]),
			Fields:      buildAWSFields(parts),
		}, nil
	default:
		return ParseResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRON.UNSUPPORTED_FLAVOR", Message: "flavor must be unix, quartz, or aws",
			Severity: engine.SevError,
		}}}, nil
	}
}

func describeUnix(parts []string) string {
	if len(parts) < 5 {
		return ""
	}
	return fmt.Sprintf(
		"min=%s hour=%s dom=%s month=%s dow=%s",
		parts[0], parts[1], parts[2], parts[3], parts[4])
}

func containsAlpha(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

func buildAWSFields(parts []string) []Field {
	out := []Field{
		{Name: "minute", Value: parts[0]},
		{Name: "hour", Value: parts[1]},
		{Name: "dayOfMonth", Value: parts[2]},
		{Name: "month", Value: parts[3]},
		{Name: "dayOfWeek", Value: parts[4]},
	}
	if len(parts) == 6 {
		out = append(out, Field{Name: "year", Value: parts[5]})
	}
	return out
}

// NextRunsResult is the success return.
type NextRunsResult struct {
	Runs        []time.Time         `json:"runs"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// NextRuns returns the next n run times after `from`. tz may be nil (UTC).
func NextRuns(expr string, from time.Time, n int, tz *time.Location) (NextRunsResult, error) {
	if tz == nil {
		tz = time.UTC
	}
	if n <= 0 {
		n = 5
	}
	if n > 1000 {
		n = 1000
	}
	sched, err := cronpkg.ParseStandard(expr)
	if err != nil {
		return NextRunsResult{Diagnostics: []engine.Diagnostic{{
			Code: "CRON.INVALID", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	t := from.In(tz)
	runs := make([]time.Time, 0, n)
	for i := 0; i < n; i++ {
		t = sched.Next(t)
		if t.IsZero() {
			break
		}
		runs = append(runs, t)
	}
	return NextRunsResult{Runs: runs}, nil
}
