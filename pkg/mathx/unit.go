package mathx

import (
	"fmt"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// UnitConvertResult holds the converted value.
type UnitConvertResult struct {
	Value       float64             `json:"value"`
	From        string              `json:"from"`
	To          string              `json:"to"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// UnitConvert converts a numeric value between supported units of the same
// dimension. Supported categories:
//
//	bytes:        b, kb, mb, gb, tb, pb (decimal); kib, mib, gib, tib, pib (binary)
//	time:         ns, us, ms, s, min, h, day, week
//	throughput:   bps, kbps, mbps, gbps, Bps, KBps, MBps, GBps  (b lower = bits, B upper = bytes)
//	temperature:  c, f, k
//	length:       mm, cm, m, km, in, ft, yd, mi
func UnitConvert(value float64, from, to string) (UnitConvertResult, error) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from == "" || to == "" {
		return UnitConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "MATH.UNIT.MISSING", Message: "both from and to are required", Severity: engine.SevError,
		}}}, nil
	}
	out, err := convertSameDimension(value, from, to)
	if err != nil {
		return UnitConvertResult{Diagnostics: []engine.Diagnostic{{
			Code: "MATH.UNIT.UNSUPPORTED", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	return UnitConvertResult{Value: out, From: from, To: to}, nil
}

func convertSameDimension(value float64, from, to string) (float64, error) {
	if dim, ok := classify(from); ok {
		if otherDim, ok2 := classify(to); !ok2 || otherDim != dim {
			return 0, fmt.Errorf("from %q (%s) and to %q (%s) are different dimensions",
				from, dim, to, otherDim)
		}
		switch dim {
		case "bytes", "throughput", "time", "length":
			factor := unitFactor(from)
			result := value * factor / unitFactor(to)
			return result, nil
		case "temperature":
			return convertTemperature(value, from, to)
		}
	}
	return 0, fmt.Errorf("unknown unit %q", from)
}

func classify(u string) (string, bool) {
	if _, ok := byteFactors[strings.ToLower(u)]; ok {
		return "bytes", true
	}
	if _, ok := throughputFactors[u]; ok {
		return "throughput", true
	}
	if _, ok := timeFactors[strings.ToLower(u)]; ok {
		return "time", true
	}
	if _, ok := lengthFactors[strings.ToLower(u)]; ok {
		return "length", true
	}
	if isTempUnit(u) {
		return "temperature", true
	}
	return "", false
}

func unitFactor(u string) float64 {
	if v, ok := byteFactors[strings.ToLower(u)]; ok {
		return v
	}
	if v, ok := throughputFactors[u]; ok {
		return v
	}
	if v, ok := timeFactors[strings.ToLower(u)]; ok {
		return v
	}
	if v, ok := lengthFactors[strings.ToLower(u)]; ok {
		return v
	}
	return 1
}

func isTempUnit(u string) bool {
	switch strings.ToLower(u) {
	case "c", "f", "k", "celsius", "fahrenheit", "kelvin":
		return true
	}
	return false
}

func convertTemperature(v float64, from, to string) (float64, error) {
	c := func() float64 {
		switch strings.ToLower(from) {
		case "c", "celsius":
			return v
		case "f", "fahrenheit":
			return (v - 32) * 5 / 9
		case "k", "kelvin":
			return v - 273.15
		}
		return v
	}()
	switch strings.ToLower(to) {
	case "c", "celsius":
		return c, nil
	case "f", "fahrenheit":
		return c*9/5 + 32, nil
	case "k", "kelvin":
		return c + 273.15, nil
	}
	return 0, fmt.Errorf("unknown temperature unit %q", to)
}

// All values are factors that convert "from this unit" → SI base.
var byteFactors = map[string]float64{
	"b":   1,
	"kb":  1e3,
	"mb":  1e6,
	"gb":  1e9,
	"tb":  1e12,
	"pb":  1e15,
	"kib": 1024,
	"mib": 1024 * 1024,
	"gib": 1024 * 1024 * 1024,
	"tib": 1024 * 1024 * 1024 * 1024,
	"pib": 1024 * 1024 * 1024 * 1024 * 1024,
}

// Throughput is dimension "bits-per-second". Bits vs bytes differ by the
// case of the unit letter ("b" = bits, "B" = bytes); the SI prefix (k/M/G)
// is matched case-insensitively for ergonomics.
var throughputFactors = map[string]float64{
	"bps":  1,
	"kbps": 1e3, "Kbps": 1e3,
	"mbps": 1e6, "Mbps": 1e6,
	"gbps": 1e9, "Gbps": 1e9,
	"Bps":  8,
	"kBps": 8 * 1e3, "KBps": 8 * 1e3,
	"mBps": 8 * 1e6, "MBps": 8 * 1e6,
	"gBps": 8 * 1e9, "GBps": 8 * 1e9,
}

var timeFactors = map[string]float64{
	"ns":   1e-9,
	"us":   1e-6,
	"ms":   1e-3,
	"s":    1,
	"min":  60,
	"h":    3600,
	"day":  86400,
	"week": 7 * 86400,
}

var lengthFactors = map[string]float64{
	"mm": 0.001,
	"cm": 0.01,
	"m":  1,
	"km": 1000,
	"in": 0.0254,
	"ft": 0.3048,
	"yd": 0.9144,
	"mi": 1609.344,
}
