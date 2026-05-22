// Package timex provides time/date utilities (timestamp ↔ date, relative,
// duration parse + breakdown).
//
// External API:
//
//	Convert(ConvertOptions) (ConvertResult, error)
//	Relative(time.Time, time.Time) (RelativeResult, error)
//	Duration(string) (DurationResult, error)
package timex

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devforge/devforge/pkg/engine"
)

// ConvertOptions selects input/output formats for Convert.
type ConvertOptions struct {
	// Input is the value: epoch seconds/ms/us/ns as string, or a date string.
	Input string `json:"input"`
	// InputFormat: "auto" (default), "epoch_s", "epoch_ms", "epoch_us", "epoch_ns",
	// "rfc3339", "iso8601".
	InputFormat string `json:"inputFormat,omitempty"`
	// TZ overrides the location for parsing/printing date strings (default UTC).
	TZ string `json:"tz,omitempty"`
}

// ConvertResult shows the same instant in multiple formats.
type ConvertResult struct {
	EpochS      int64               `json:"epochS"`
	EpochMS     int64               `json:"epochMS"`
	EpochUS     int64               `json:"epochUS"`
	EpochNS     int64               `json:"epochNS"`
	RFC3339     string              `json:"rfc3339"`
	UTC         string              `json:"utc"`
	Local       string              `json:"local"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Convert parses a moment and re-emits it in every supported format.
func Convert(opts ConvertOptions) (ConvertResult, error) {
	loc := time.UTC
	if opts.TZ != "" {
		l, err := time.LoadLocation(opts.TZ)
		if err != nil {
			return ConvertResult{Diagnostics: []engine.Diagnostic{{
				Code: "TIME.CONVERT.UNKNOWN_TZ", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
		loc = l
	}
	t, err := parseInput(opts.Input, opts.InputFormat, loc)
	if err != nil {
		return ConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "TIME.CONVERT.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	utc := t.UTC()
	return ConvertResult{
		EpochS:  utc.Unix(),
		EpochMS: utc.UnixMilli(),
		EpochUS: utc.UnixMicro(),
		EpochNS: utc.UnixNano(),
		RFC3339: utc.Format(time.RFC3339Nano),
		UTC:     utc.Format("2006-01-02 15:04:05.000 MST"),
		Local:   t.In(loc).Format("2006-01-02 15:04:05.000 MST"),
	}, nil
}

func parseInput(in, format string, loc *time.Location) (time.Time, error) {
	in = strings.TrimSpace(in)
	if in == "" {
		return time.Time{}, fmt.Errorf("input is empty")
	}
	if format == "" || format == "auto" {
		// Heuristic: if all digits, treat as epoch and decide on magnitude.
		if isAllDigits(in) {
			return parseEpoch(in)
		}
		// Try common date layouts.
		for _, layout := range []string{
			time.RFC3339Nano, time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02",
		} {
			if t, err := time.ParseInLocation(layout, in, loc); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot auto-detect format for %q", in)
	}
	switch format {
	case "epoch_s", "epoch_ms", "epoch_us", "epoch_ns":
		return parseEpochExplicit(in, format)
	case "rfc3339", "iso8601":
		t, err := time.Parse(time.RFC3339Nano, in)
		if err != nil {
			return time.Parse(time.RFC3339, in)
		}
		return t, err
	default:
		return time.Time{}, fmt.Errorf("unknown inputFormat %q", format)
	}
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func parseEpoch(s string) (time.Time, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	// Decide unit by magnitude:
	//   < 10^11   → seconds (~year 5138 cap)
	//   < 10^14   → milliseconds
	//   < 10^17   → microseconds
	//   else      → nanoseconds
	abs := v
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs < 1e11:
		return time.Unix(v, 0), nil
	case abs < 1e14:
		return time.UnixMilli(v), nil
	case abs < 1e17:
		return time.UnixMicro(v), nil
	default:
		return time.Unix(0, v), nil
	}
}

func parseEpochExplicit(s, format string) (time.Time, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	switch format {
	case "epoch_s":
		return time.Unix(v, 0), nil
	case "epoch_ms":
		return time.UnixMilli(v), nil
	case "epoch_us":
		return time.UnixMicro(v), nil
	case "epoch_ns":
		return time.Unix(0, v), nil
	}
	return time.Time{}, fmt.Errorf("unknown epoch format %q", format)
}

// RelativeResult describes a delta between two instants.
type RelativeResult struct {
	Phrase  string  `json:"phrase"`  // "5 minutes ago" / "in 3 days"
	Seconds float64 `json:"seconds"` // signed delta from->to
}

// Relative produces a human phrase for the gap between from and to.
func Relative(from, to time.Time) (RelativeResult, error) {
	d := to.Sub(from)
	res := RelativeResult{Seconds: d.Seconds()}
	res.Phrase = humanize(d)
	return res, nil
}

func humanize(d time.Duration) string {
	abs := d
	suffix := " from now"
	if d < 0 {
		abs = -d
		suffix = " ago"
	}
	switch {
	case abs < time.Second:
		return "just now"
	case abs < time.Minute:
		return fmt.Sprintf("%d second%s%s", int(abs.Seconds()), plural(int(abs.Seconds())), suffix)
	case abs < time.Hour:
		m := int(abs.Minutes())
		return fmt.Sprintf("%d minute%s%s", m, plural(m), suffix)
	case abs < 24*time.Hour:
		h := int(abs.Hours())
		return fmt.Sprintf("%d hour%s%s", h, plural(h), suffix)
	case abs < 30*24*time.Hour:
		dy := int(abs.Hours() / 24)
		return fmt.Sprintf("%d day%s%s", dy, plural(dy), suffix)
	case abs < 365*24*time.Hour:
		mo := int(abs.Hours() / 24 / 30)
		return fmt.Sprintf("%d month%s%s", mo, plural(mo), suffix)
	}
	yr := int(abs.Hours() / 24 / 365)
	return fmt.Sprintf("%d year%s%s", yr, plural(yr), suffix)
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// DurationResult breaks down a duration into named components.
type DurationResult struct {
	TotalSeconds float64             `json:"totalSeconds"`
	Days         int                 `json:"days"`
	Hours        int                 `json:"hours"`
	Minutes      int                 `json:"minutes"`
	Seconds      int                 `json:"seconds"`
	Milliseconds int                 `json:"milliseconds"`
	Phrase       string              `json:"phrase"`
	Diagnostics  []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Duration parses a Go-style duration string ("1h30m15s") and breaks it down.
// Also accepts plain integer seconds.
func Duration(input string) (DurationResult, error) {
	d, err := time.ParseDuration(strings.TrimSpace(input))
	if err != nil {
		// Try plain integer seconds.
		if v, perr := strconv.ParseInt(strings.TrimSpace(input), 10, 64); perr == nil {
			d = time.Duration(v) * time.Second
		} else {
			return DurationResult{Diagnostics: []engine.Diagnostic{{
				Code: "TIME.DURATION.PARSE", Message: err.Error(), Severity: engine.SevError,
			}}}, nil
		}
	}
	abs := d
	if abs < 0 {
		abs = -abs
	}
	days := int(abs / (24 * time.Hour))
	abs -= time.Duration(days) * 24 * time.Hour
	hrs := int(abs / time.Hour)
	abs -= time.Duration(hrs) * time.Hour
	mins := int(abs / time.Minute)
	abs -= time.Duration(mins) * time.Minute
	secs := int(abs / time.Second)
	abs -= time.Duration(secs) * time.Second
	ms := int(abs / time.Millisecond)
	return DurationResult{
		TotalSeconds: d.Seconds(),
		Days:         days,
		Hours:        hrs,
		Minutes:      mins,
		Seconds:      secs,
		Milliseconds: ms,
		Phrase:       d.String(),
	}, nil
}
