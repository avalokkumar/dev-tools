package mathx

import (
	"math"
	"testing"
)

func TestEval_Basic(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		in   string
		want float64
	}{
		{"1+2", 3},
		{"2*3+4", 10},
		{"(2+3)*4", 20},
		{"10/4", 2.5},
		{"7%3", 1},
		{"2**10", 1024},
		{"-3+5", 2},
		{"sqrt(9)", 3},
		{"pow(2,8)", 256},
		{"abs(-7)", 7},
		{"min(3,5)", 3},
		{"max(3,5)", 5},
		{"round(2.5)", 3},
		{"pi", math.Pi},
	} {
		r, _ := Eval(c.in)
		if math.Abs(r.Value-c.want) > 1e-9 {
			t.Fatalf("%s = %v, want %v", c.in, r.Value, c.want)
		}
	}
}

func TestEval_DivByZero(t *testing.T) {
	t.Parallel()
	r, _ := Eval("5/0")
	if r.Diagnostics[0].Code != "MATH.EVAL.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestEval_BadSyntax(t *testing.T) {
	t.Parallel()
	r, _ := Eval("2 + (3 ")
	if r.Diagnostics[0].Code != "MATH.EVAL.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestEval_UnknownFn(t *testing.T) {
	t.Parallel()
	r, _ := Eval("hax(1)")
	if r.Diagnostics[0].Code != "MATH.EVAL.PARSE" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

// TestEval_NaN — M-9: sqrt(-1) yields MATH.EVAL.NAN diagnostic.
// Display preserves "NaN" while Value is zeroed for JSON safety.
func TestEval_NaN(t *testing.T) {
	t.Parallel()
	r, _ := Eval("sqrt(-1)")
	if r.Display != "NaN" {
		t.Fatalf("Display = %q, want NaN", r.Display)
	}
	if r.Diagnostics[0].Code != "MATH.EVAL.NAN" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

// TestEval_Inf — M-9: 2**10000 overflows; diagnostic surfaces; Value zeroed.
func TestEval_Inf(t *testing.T) {
	t.Parallel()
	r, _ := Eval("2**10000")
	if r.Display != "+Inf" {
		t.Fatalf("Display = %q, want +Inf", r.Display)
	}
	if r.Diagnostics[0].Code != "MATH.EVAL.INF" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}

func TestUnitConvert_Bytes(t *testing.T) {
	t.Parallel()
	r, _ := UnitConvert(1, "gib", "mib")
	if r.Value != 1024 {
		t.Fatalf("1 GiB = %v MiB, want 1024", r.Value)
	}
}

func TestUnitConvert_Time(t *testing.T) {
	t.Parallel()
	r, _ := UnitConvert(2, "h", "min")
	if r.Value != 120 {
		t.Fatalf("got %v", r.Value)
	}
}

func TestUnitConvert_Throughput_BitsToBytes(t *testing.T) {
	t.Parallel()
	r, _ := UnitConvert(1, "Mbps", "MBps")
	if r.Value != 0.125 {
		t.Fatalf("1 Mbps = %v MBps, want 0.125", r.Value)
	}
}

func TestUnitConvert_TempCtoF(t *testing.T) {
	t.Parallel()
	r, _ := UnitConvert(100, "c", "f")
	if r.Value != 212 {
		t.Fatalf("100C = %v F", r.Value)
	}
}

func TestUnitConvert_DimensionMismatch(t *testing.T) {
	t.Parallel()
	r, _ := UnitConvert(1, "kb", "min")
	if r.Diagnostics[0].Code != "MATH.UNIT.UNSUPPORTED" {
		t.Fatalf("code = %q", r.Diagnostics[0].Code)
	}
}
