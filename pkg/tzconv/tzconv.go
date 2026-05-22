// Package tzconv converts times between IANA time zones with DST awareness.
//
// External API:
//
//	Convert(t time.Time, fromTZ, toTZ string) (ConvertResult, error)
//	ListZones(filter string) ([]Zone, error)
package tzconv

import (
	"strings"
	"time"

	"github.com/devforge/devforge/pkg/engine"
)

// ConvertResult is the success return.
type ConvertResult struct {
	Original    time.Time           `json:"original"`
	Converted   time.Time           `json:"converted"`
	DSTNote     string              `json:"dstNote,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Convert reinterprets t (assumed to be in fromTZ if it has zero zone info,
// otherwise in its own zone) as a moment in toTZ. DSTNote flags gap/overlap.
func Convert(t time.Time, fromTZ, toTZ string) (ConvertResult, error) {
	from, err := time.LoadLocation(fromTZ)
	if err != nil {
		return ConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "TZ.UNKNOWN_FROM", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	to, err := time.LoadLocation(toTZ)
	if err != nil {
		return ConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "TZ.UNKNOWN_TO", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	// Reinterpret wall components in fromTZ. Capture request fields first so
	// we can detect normalization (i.e. the wall time fell into a DST gap).
	yy, mm, dd := t.Year(), t.Month(), t.Day()
	hh, mi, ss, ns := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
	wall := time.Date(yy, mm, dd, hh, mi, ss, ns, from)
	note := dstNote(wall, from, hh, dd)
	converted := wall.In(to)
	return ConvertResult{Original: wall, Converted: converted, DSTNote: note}, nil
}

// dstNote classifies a wall time relative to its zone.
//   - "gap":     the requested wall components do not exist (spring-forward).
//   - "overlap": they exist twice (fall-back).
//   - "":        unambiguous.
func dstNote(wall time.Time, loc *time.Location, reqHour, reqDay int) string {
	// Go's time.Date normalises non-existent wall times by pushing forward by
	// the DST step, so the resulting hour/day will differ from the request.
	if wall.Hour() != reqHour || wall.Day() != reqDay {
		return "gap"
	}
	// Overlap detection: compare offsets one hour before and at the wall time.
	// In a fall-back overlap, the offset before is greater than at wall.
	prev := wall.Add(-time.Hour)
	_, offBefore := prev.In(loc).Zone()
	_, offAt := wall.Zone()
	if offBefore > offAt {
		return "overlap"
	}
	return ""
}

// Zone is a single IANA zone snapshot at the current moment.
type Zone struct {
	Name          string `json:"name"`
	Abbrev        string `json:"abbrev"`
	OffsetSeconds int    `json:"offsetSeconds"`
	IsDST         bool   `json:"isDst"`
}

// ListZones returns a curated set of well-known IANA zones, optionally
// filtered by case-insensitive substring of the name. The Go runtime does
// not enumerate the zoneinfo database, so we keep a popular allow-list.
func ListZones(filter string) ([]Zone, error) {
	now := time.Now()
	out := make([]Zone, 0, len(commonZones))
	f := strings.ToLower(filter)
	for _, name := range commonZones {
		if f != "" && !strings.Contains(strings.ToLower(name), f) {
			continue
		}
		loc, err := time.LoadLocation(name)
		if err != nil {
			continue
		}
		t := now.In(loc)
		ab, off := t.Zone()
		out = append(out, Zone{
			Name:          name,
			Abbrev:        ab,
			OffsetSeconds: off,
			IsDST:         t.IsDST(),
		})
	}
	return out, nil
}

// commonZones is a hand-picked subset that covers most user demand. Full
// zoneinfo enumeration is left for a Phase D enhancement.
var commonZones = []string{
	"UTC",
	"Africa/Cairo",
	"Africa/Johannesburg",
	"Africa/Lagos",
	"America/Anchorage",
	"America/Argentina/Buenos_Aires",
	"America/Chicago",
	"America/Denver",
	"America/Los_Angeles",
	"America/Mexico_City",
	"America/New_York",
	"America/Sao_Paulo",
	"America/Toronto",
	"America/Vancouver",
	"Asia/Bangkok",
	"Asia/Dubai",
	"Asia/Hong_Kong",
	"Asia/Jakarta",
	"Asia/Karachi",
	"Asia/Kolkata",
	"Asia/Manila",
	"Asia/Seoul",
	"Asia/Shanghai",
	"Asia/Singapore",
	"Asia/Tokyo",
	"Australia/Melbourne",
	"Australia/Sydney",
	"Europe/Amsterdam",
	"Europe/Berlin",
	"Europe/Dublin",
	"Europe/Istanbul",
	"Europe/London",
	"Europe/Madrid",
	"Europe/Moscow",
	"Europe/Paris",
	"Europe/Rome",
	"Europe/Stockholm",
	"Europe/Warsaw",
	"Europe/Zurich",
	"Pacific/Auckland",
	"Pacific/Honolulu",
}
