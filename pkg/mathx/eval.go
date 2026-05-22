// Package mathx provides a safe arithmetic expression evaluator + a
// dev-focused unit converter. No I/O, no shell, no third-party libs.
package mathx

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/devforge/devforge/pkg/engine"
)

// EvalResult is the value of a successfully evaluated expression.
type EvalResult struct {
	Value       float64             `json:"value"`
	Display     string              `json:"display"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Eval evaluates an arithmetic expression. Supports + - * / % ** parens,
// unary minus, and a few stdlib functions (sqrt, pow, abs, log, log2, exp).
func Eval(expr string) (EvalResult, error) {
	p := &parser{src: expr}
	v, err := p.parseExpression()
	if err != nil {
		return EvalResult{Diagnostics: []engine.Diagnostic{{
			Code: "MATH.EVAL.PARSE", Message: err.Error(), Severity: engine.SevError,
		}}}, nil
	}
	if p.pos < len(p.src) {
		p.skipWS()
		if p.pos < len(p.src) {
			return EvalResult{Diagnostics: []engine.Diagnostic{{
				Code: "MATH.EVAL.TRAILING",
				Message: fmt.Sprintf("unexpected input after expression: %q", p.src[p.pos:]),
				Severity: engine.SevError,
			}}}, nil
		}
	}
	res := EvalResult{Value: v, Display: format(v)}
	// JSON cannot represent NaN/Inf; surface diagnostic and zero the value
	// so json.Marshal succeeds. Display already shows "NaN"/"+Inf"/"-Inf".
	if math.IsNaN(v) {
		res.Value = 0
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "MATH.EVAL.NAN",
			Message: "expression produced NaN (e.g. 0/0, sqrt(-1), log(-x))",
			Severity: engine.SevError,
		})
	}
	if math.IsInf(v, 0) {
		res.Value = 0
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "MATH.EVAL.INF",
			Message: "expression overflowed to infinity",
			Severity: engine.SevError,
		})
	}
	return res, nil
}

func format(v float64) string {
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 0) {
		if v > 0 {
			return "+Inf"
		}
		return "-Inf"
	}
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// ---- recursive-descent parser ----

type parser struct {
	src string
	pos int
}

func (p *parser) skipWS() {
	for p.pos < len(p.src) && unicode.IsSpace(rune(p.src[p.pos])) {
		p.pos++
	}
}

func (p *parser) parseExpression() (float64, error) {
	v, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		p.skipWS()
		if p.pos >= len(p.src) {
			return v, nil
		}
		c := p.src[p.pos]
		if c != '+' && c != '-' {
			return v, nil
		}
		p.pos++
		rhs, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if c == '+' {
			v += rhs
		} else {
			v -= rhs
		}
	}
}

func (p *parser) parseTerm() (float64, error) {
	v, err := p.parsePower()
	if err != nil {
		return 0, err
	}
	for {
		p.skipWS()
		if p.pos >= len(p.src) {
			return v, nil
		}
		c := p.src[p.pos]
		if c != '*' && c != '/' && c != '%' {
			return v, nil
		}
		p.pos++
		rhs, err := p.parsePower()
		if err != nil {
			return 0, err
		}
		switch c {
		case '*':
			v *= rhs
		case '/':
			if rhs == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			v /= rhs
		case '%':
			if rhs == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			v = math.Mod(v, rhs)
		}
	}
}

func (p *parser) parsePower() (float64, error) {
	v, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	p.skipWS()
	if p.pos+1 < len(p.src) && p.src[p.pos] == '*' && p.src[p.pos+1] == '*' {
		p.pos += 2
		rhs, err := p.parsePower() // right-associative
		if err != nil {
			return 0, err
		}
		return math.Pow(v, rhs), nil
	}
	return v, nil
}

func (p *parser) parseUnary() (float64, error) {
	p.skipWS()
	if p.pos >= len(p.src) {
		return 0, fmt.Errorf("unexpected end of input")
	}
	switch p.src[p.pos] {
	case '+':
		p.pos++
		return p.parseUnary()
	case '-':
		p.pos++
		v, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (float64, error) {
	p.skipWS()
	if p.pos >= len(p.src) {
		return 0, fmt.Errorf("unexpected end of input")
	}
	c := p.src[p.pos]
	if c == '(' {
		p.pos++
		v, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWS()
		if p.pos >= len(p.src) || p.src[p.pos] != ')' {
			return 0, fmt.Errorf("expected ')'")
		}
		p.pos++
		return v, nil
	}
	if unicode.IsLetter(rune(c)) {
		start := p.pos
		for p.pos < len(p.src) && (unicode.IsLetter(rune(p.src[p.pos])) || unicode.IsDigit(rune(p.src[p.pos]))) {
			p.pos++
		}
		name := p.src[start:p.pos]
		// Constants:
		switch strings.ToLower(name) {
		case "pi":
			return math.Pi, nil
		case "e":
			return math.E, nil
		}
		// Functions: name(args)
		p.skipWS()
		if p.pos >= len(p.src) || p.src[p.pos] != '(' {
			return 0, fmt.Errorf("unknown identifier %q", name)
		}
		p.pos++
		args, err := p.parseArgs()
		if err != nil {
			return 0, err
		}
		return callFunc(strings.ToLower(name), args)
	}
	if unicode.IsDigit(rune(c)) || c == '.' {
		return p.parseNumber()
	}
	return 0, fmt.Errorf("unexpected character %q", c)
}

func (p *parser) parseArgs() ([]float64, error) {
	var out []float64
	p.skipWS()
	if p.pos < len(p.src) && p.src[p.pos] == ')' {
		p.pos++
		return out, nil
	}
	for {
		v, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
		p.skipWS()
		if p.pos >= len(p.src) {
			return nil, fmt.Errorf("expected ')'")
		}
		if p.src[p.pos] == ',' {
			p.pos++
			continue
		}
		if p.src[p.pos] == ')' {
			p.pos++
			return out, nil
		}
		return nil, fmt.Errorf("expected ',' or ')'")
	}
}

func (p *parser) parseNumber() (float64, error) {
	start := p.pos
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if !(unicode.IsDigit(rune(c)) || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-') {
			break
		}
		// '+' / '-' only valid as part of an exponent.
		if (c == '+' || c == '-') && (p.pos == start || (p.src[p.pos-1] != 'e' && p.src[p.pos-1] != 'E')) {
			break
		}
		p.pos++
	}
	v, err := strconv.ParseFloat(p.src[start:p.pos], 64)
	if err != nil {
		return 0, fmt.Errorf("bad number %q", p.src[start:p.pos])
	}
	return v, nil
}

func callFunc(name string, args []float64) (float64, error) {
	check := func(want int) error {
		if len(args) != want {
			return fmt.Errorf("%s wants %d arg(s), got %d", name, want, len(args))
		}
		return nil
	}
	switch name {
	case "sqrt":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Sqrt(args[0]), nil
	case "pow":
		if err := check(2); err != nil {
			return 0, err
		}
		return math.Pow(args[0], args[1]), nil
	case "abs":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Abs(args[0]), nil
	case "log":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Log(args[0]), nil
	case "log2":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Log2(args[0]), nil
	case "log10":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Log10(args[0]), nil
	case "exp":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Exp(args[0]), nil
	case "sin":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Sin(args[0]), nil
	case "cos":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Cos(args[0]), nil
	case "tan":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Tan(args[0]), nil
	case "min":
		if err := check(2); err != nil {
			return 0, err
		}
		return math.Min(args[0], args[1]), nil
	case "max":
		if err := check(2); err != nil {
			return 0, err
		}
		return math.Max(args[0], args[1]), nil
	case "floor":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Floor(args[0]), nil
	case "ceil":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Ceil(args[0]), nil
	case "round":
		if err := check(1); err != nil {
			return 0, err
		}
		return math.Round(args[0]), nil
	}
	return 0, fmt.Errorf("unknown function %q", name)
}
