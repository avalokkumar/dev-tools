// Package colorx converts colors between hex, rgb, and hsl.
//
// External API:
//
//	Convert(input string, opts Options) (Result, error)
package colorx

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
)

// Options tunes Convert.
type Options struct {
	// To is the target representation: "hex", "rgb", "hsl", or "" for "all".
	To string `json:"to,omitempty"`
}

// Result reports the color in every standard representation.
type Result struct {
	Hex         string              `json:"hex"`
	RGB         string              `json:"rgb"`
	HSL         string              `json:"hsl"`
	R           int                 `json:"r"`
	G           int                 `json:"g"`
	B           int                 `json:"b"`
	H           int                 `json:"h"`
	S           int                 `json:"s"`
	L           int                 `json:"l"`
	Output      string              `json:"output,omitempty"` // populated when To != ""
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Convert parses input and re-emits in every representation. Accepts:
//   - "#abc"
//   - "#aabbcc"
//   - "rgb(170, 187, 204)"
//   - "hsl(210, 25%, 73%)"
//   - "red", "blue" (CSS named colors)
func Convert(input string, opts Options) (Result, error) {
	r, g, b, ok := parseAny(input)
	if !ok {
		return Result{Diagnostics: []engine.Diagnostic{{
			Code: "COLOR.PARSE", Message: fmt.Sprintf("cannot parse color %q", input),
			Severity: engine.SevError,
		}}}, nil
	}
	res := Result{R: r, G: g, B: b}
	res.Hex = fmt.Sprintf("#%02x%02x%02x", r, g, b)
	res.RGB = fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
	h, s, l := rgbToHSL(r, g, b)
	res.H, res.S, res.L = h, s, l
	res.HSL = fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l)
	switch strings.ToLower(opts.To) {
	case "hex":
		res.Output = res.Hex
	case "rgb":
		res.Output = res.RGB
	case "hsl":
		res.Output = res.HSL
	}
	return res, nil
}

func parseAny(s string) (int, int, int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 0, false
	}
	// Named colors first.
	if r, g, b, ok := named(strings.ToLower(s)); ok {
		return r, g, b, true
	}
	// Hex.
	if strings.HasPrefix(s, "#") {
		return parseHex(s[1:])
	}
	// rgb()
	if strings.HasPrefix(strings.ToLower(s), "rgb") {
		return parseRGB(s)
	}
	// hsl()
	if strings.HasPrefix(strings.ToLower(s), "hsl") {
		return parseHSL(s)
	}
	// Bare hex without #.
	return parseHex(s)
}

func parseHex(s string) (int, int, int, bool) {
	s = strings.TrimSpace(s)
	if len(s) == 3 {
		// Expand "abc" → "aabbcc".
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 {
		return 0, 0, 0, false
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(v >> 16 & 0xff), int(v >> 8 & 0xff), int(v & 0xff), true
}

func parseRGB(s string) (int, int, int, bool) {
	open := strings.Index(s, "(")
	close := strings.LastIndex(s, ")")
	if open < 0 || close < 0 || close <= open {
		return 0, 0, 0, false
	}
	parts := strings.Split(s[open+1:close], ",")
	if len(parts) < 3 {
		return 0, 0, 0, false
	}
	r, ok1 := parseChannel(parts[0])
	g, ok2 := parseChannel(parts[1])
	b, ok3 := parseChannel(parts[2])
	return r, g, b, ok1 && ok2 && ok3
}

func parseChannel(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, false
		}
		return int(v / 100 * 255), true
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	if v < 0 || v > 255 {
		return 0, false
	}
	return v, true
}

func parseHSL(s string) (int, int, int, bool) {
	open := strings.Index(s, "(")
	close := strings.LastIndex(s, ")")
	if open < 0 || close < 0 {
		return 0, 0, 0, false
	}
	parts := strings.Split(s[open+1:close], ",")
	if len(parts) < 3 {
		return 0, 0, 0, false
	}
	h, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, false
	}
	sat, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"), 64)
	if err != nil {
		return 0, 0, 0, false
	}
	li, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"), 64)
	if err != nil {
		return 0, 0, 0, false
	}
	r, g, b := hslToRGB(h, sat/100, li/100)
	return r, g, b, true
}

func rgbToHSL(r, g, b int) (int, int, int) {
	rf, gf, bf := float64(r)/255, float64(g)/255, float64(b)/255
	maxv := math.Max(rf, math.Max(gf, bf))
	minv := math.Min(rf, math.Min(gf, bf))
	l := (maxv + minv) / 2
	var h, s float64
	if maxv == minv {
		h, s = 0, 0
	} else {
		d := maxv - minv
		if l > 0.5 {
			s = d / (2 - maxv - minv)
		} else {
			s = d / (maxv + minv)
		}
		switch maxv {
		case rf:
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/d + 2
		case bf:
			h = (rf-gf)/d + 4
		}
		h *= 60
	}
	return int(math.Round(h)), int(math.Round(s * 100)), int(math.Round(l * 100))
}

func hslToRGB(h, s, l float64) (int, int, int) {
	if s == 0 {
		v := int(math.Round(l * 255))
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	conv := func(t float64) int {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		var v float64
		switch {
		case t < 1.0/6:
			v = p + (q-p)*6*t
		case t < 0.5:
			v = q
		case t < 2.0/3:
			v = p + (q-p)*(2.0/3-t)*6
		default:
			v = p
		}
		return int(math.Round(v * 255))
	}
	hp := h / 360
	return conv(hp + 1.0/3), conv(hp), conv(hp - 1.0/3)
}

// Common CSS named colors. Subset; covers the highest-frequency names.
func named(s string) (int, int, int, bool) {
	t, ok := cssNames[s]
	if !ok {
		return 0, 0, 0, false
	}
	return t[0], t[1], t[2], true
}

var cssNames = map[string][3]int{
	"black":   {0, 0, 0},
	"white":   {255, 255, 255},
	"red":     {255, 0, 0},
	"green":   {0, 128, 0},
	"blue":    {0, 0, 255},
	"yellow":  {255, 255, 0},
	"cyan":    {0, 255, 255},
	"magenta": {255, 0, 255},
	"gray":    {128, 128, 128},
	"grey":    {128, 128, 128},
	"silver":  {192, 192, 192},
	"orange":  {255, 165, 0},
	"purple":  {128, 0, 128},
	"pink":    {255, 192, 203},
	"brown":   {165, 42, 42},
	"navy":    {0, 0, 128},
	"teal":    {0, 128, 128},
	"olive":   {128, 128, 0},
	"maroon":  {128, 0, 0},
	"lime":    {0, 255, 0},
	"aqua":    {0, 255, 255},
	"fuchsia": {255, 0, 255},
}
