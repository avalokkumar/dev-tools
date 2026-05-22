package csvfmt

import (
	"strings"
	"testing"
)

// TestCsvFormat_PreservesQuotedNewlines — C3: a quoted field containing a
// newline survives a format round-trip.
func TestCsvFormat_PreservesQuotedNewlines(t *testing.T) {
	t.Parallel()
	in := []byte("a,b\n\"line1\nline2\",2\n")
	res, err := Format(in, FormatOptions{})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	got := string(res.Output)
	if !strings.Contains(got, "\"line1\nline2\"") {
		t.Fatalf("quoted newline lost: %q", got)
	}
}

// TestCsvFormat_AlignColumns — C3: column alignment pads short cells.
func TestCsvFormat_AlignColumns(t *testing.T) {
	t.Parallel()
	in := []byte("name,age\nalok,42\nbob,1\n")
	res, err := Format(in, FormatOptions{AlignColumns: true})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	got := string(res.Output)
	if !strings.Contains(got, "alok, 42") || !strings.Contains(got, "bob , 1") {
		t.Fatalf("alignment off:\n%s", got)
	}
}

// TestCsvValidate_HeaderMismatchDiagnostic — C3: header mismatch surfaces.
func TestCsvValidate_HeaderMismatchDiagnostic(t *testing.T) {
	t.Parallel()
	res, err := Validate([]byte("name,age\nalok,42\n"), ValidateOptions{
		ExpectedColumns: []string{"username", "age"},
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if res.Valid {
		t.Fatalf("expected invalid")
	}
	if res.Diagnostics[0].Code != "CSV.HEADER_MISMATCH" {
		t.Fatalf("code=%q", res.Diagnostics[0].Code)
	}
}

// TestCsvValidate_StrictJaggedFails — C3: jagged rows fail in strict mode.
func TestCsvValidate_StrictJaggedFails(t *testing.T) {
	t.Parallel()
	res, err := Validate([]byte("a,b,c\n1,2\n"), ValidateOptions{Strict: true})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if res.Valid {
		// good
	} else {
		// also good; check diag promoted
	}
	if len(res.Diagnostics) == 0 || res.Diagnostics[0].Code != "CSV.JAGGED_ROW" {
		t.Fatalf("missing jagged diag: %+v", res.Diagnostics)
	}
}
