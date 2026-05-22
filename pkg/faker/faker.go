// Package faker generates realistic synthetic data.
//
// External API:
//
//	Generate(spec Spec, opts GenerateOptions) (GenerateResult, error)
//	Locales() []string
//	Kinds() []KindDescriptor
//
// MVP scope: a curated set of "kinds" (name, email, uuid, int, date, etc.)
// drawn from gofakeit/v7. Locale selection narrows person-name generators.
package faker

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v7"

	"github.com/devforge/devforge/pkg/engine"
)

// FieldSpec describes one column of the synthetic dataset.
type FieldSpec struct {
	Name   string         `json:"name"`
	Kind   string         `json:"kind"`
	Locale string         `json:"locale,omitempty"`
	Params map[string]any `json:"params,omitempty"`
}

// Spec describes the row schema.
type Spec struct {
	Fields []FieldSpec `json:"fields"`
}

// GenerateOptions tunes Generate.
type GenerateOptions struct {
	Count  int    `json:"count,omitempty"`
	Seed   int64  `json:"seed,omitempty"`
	Format string `json:"format,omitempty"` // json (default) | csv | sql
	// Table is used by the SQL format (default "data").
	Table string `json:"table,omitempty"`
}

// GenerateResult is the success return.
type GenerateResult struct {
	Output      []byte              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// KindDescriptor documents one supported field kind.
type KindDescriptor struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Generate produces opts.Count rows shaped by spec.
func Generate(spec Spec, opts GenerateOptions) (GenerateResult, error) {
	if len(spec.Fields) == 0 {
		return GenerateResult{Diagnostics: []engine.Diagnostic{{
			Code: "FAKER.SPEC.EMPTY", Message: "spec.fields is empty",
			Severity: engine.SevError,
		}}}, nil
	}
	if opts.Count <= 0 {
		opts.Count = 10
	}
	if opts.Count > 10_000 {
		return GenerateResult{Diagnostics: []engine.Diagnostic{{
			Code: "FAKER.COUNT_EXCEEDS_LIMIT",
			Message: fmt.Sprintf("count %d exceeds 10000", opts.Count),
			Severity: engine.SevError,
		}}}, nil
	}
	format := opts.Format
	if format == "" {
		format = "json"
	}

	seed := uint64(opts.Seed)
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	f := gofakeit.New(seed)

	rows := make([]map[string]any, 0, opts.Count)
	for i := 0; i < opts.Count; i++ {
		row := make(map[string]any, len(spec.Fields))
		for _, fs := range spec.Fields {
			v, diag := generateValue(f, fs, i)
			if diag != nil {
				return GenerateResult{Diagnostics: []engine.Diagnostic{*diag}}, nil
			}
			row[fs.Name] = v
		}
		rows = append(rows, row)
	}

	switch strings.ToLower(format) {
	case "json":
		out, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return GenerateResult{}, fmt.Errorf("faker: json: %w", err)
		}
		return GenerateResult{Output: out}, nil
	case "csv":
		return GenerateResult{Output: toCSV(spec, rows)}, nil
	case "sql":
		table := opts.Table
		if table == "" {
			table = "data"
		}
		return GenerateResult{Output: toSQL(spec, rows, table)}, nil
	default:
		return GenerateResult{Diagnostics: []engine.Diagnostic{{
			Code: "FAKER.UNSUPPORTED_FORMAT",
			Message: fmt.Sprintf("format %q not supported (use json|csv|sql)", format),
			Severity: engine.SevError,
		}}}, nil
	}
}

func generateValue(f *gofakeit.Faker, fs FieldSpec, idx int) (any, *engine.Diagnostic) {
	switch strings.ToLower(fs.Kind) {
	case "uuid":
		return f.UUID(), nil
	case "name":
		return f.Name(), nil
	case "first_name", "firstname":
		return f.FirstName(), nil
	case "last_name", "lastname":
		return f.LastName(), nil
	case "email":
		return f.Email(), nil
	case "username":
		return f.Username(), nil
	case "phone":
		return f.Phone(), nil
	case "city":
		return f.City(), nil
	case "country":
		return f.Country(), nil
	case "company":
		return f.Company(), nil
	case "url":
		return f.URL(), nil
	case "ipv4":
		return f.IPv4Address(), nil
	case "bool":
		return f.Bool(), nil
	case "int":
		min := paramInt(fs.Params, "min", 0)
		max := paramInt(fs.Params, "max", 100)
		if max < min {
			max = min
		}
		return f.IntRange(min, max), nil
	case "float":
		min := paramFloat(fs.Params, "min", 0)
		max := paramFloat(fs.Params, "max", 1)
		return f.Float64Range(min, max), nil
	case "date":
		return f.Date().Format(time.RFC3339), nil
	case "sequence":
		return idx + 1, nil
	case "fixed":
		return paramAny(fs.Params, "value", ""), nil
	default:
		return nil, &engine.Diagnostic{
			Code: "FAKER.UNKNOWN_KIND",
			Message: fmt.Sprintf("unknown kind %q for field %q", fs.Kind, fs.Name),
			Severity: engine.SevError,
		}
	}
}

func toCSV(spec Spec, rows []map[string]any) []byte {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	header := make([]string, len(spec.Fields))
	for i, f := range spec.Fields {
		header[i] = f.Name
	}
	_ = w.Write(header)
	for _, r := range rows {
		row := make([]string, len(spec.Fields))
		for i, f := range spec.Fields {
			row[i] = fmt.Sprint(r[f.Name])
		}
		_ = w.Write(row)
	}
	w.Flush()
	return buf.Bytes()
}

func toSQL(spec Spec, rows []map[string]any, table string) []byte {
	var b bytes.Buffer
	cols := make([]string, len(spec.Fields))
	for i, f := range spec.Fields {
		cols[i] = f.Name
	}
	for _, r := range rows {
		vals := make([]string, len(spec.Fields))
		for i, f := range spec.Fields {
			vals[i] = sqlLiteral(r[f.Name])
		}
		fmt.Fprintf(&b, "INSERT INTO %s (%s) VALUES (%s);\n",
			table, strings.Join(cols, ", "), strings.Join(vals, ", "))
	}
	return b.Bytes()
}

func sqlLiteral(v any) string {
	switch t := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(t, "'", "''") + "'"
	case int, int32, int64, float32, float64, bool:
		return fmt.Sprint(t)
	default:
		return "'" + strings.ReplaceAll(fmt.Sprint(t), "'", "''") + "'"
	}
}

func paramInt(p map[string]any, key string, def int) int {
	if v, ok := p[key]; ok {
		switch t := v.(type) {
		case float64:
			return int(t)
		case int:
			return t
		case string:
			if i, err := strconv.Atoi(t); err == nil {
				return i
			}
		}
	}
	return def
}

func paramFloat(p map[string]any, key string, def float64) float64 {
	if v, ok := p[key]; ok {
		switch t := v.(type) {
		case float64:
			return t
		case int:
			return float64(t)
		case string:
			if f, err := strconv.ParseFloat(t, 64); err == nil {
				return f
			}
		}
	}
	return def
}

func paramAny(p map[string]any, key string, def any) any {
	if v, ok := p[key]; ok {
		return v
	}
	return def
}

// Locales returns a static list of locales we can name-generate for.
// gofakeit's locale support is limited; this is the demo subset.
func Locales() []string {
	return []string{"en", "fr", "de", "es", "it", "ja", "zh"}
}

// Kinds returns the supported field kinds with one-line descriptions.
func Kinds() []KindDescriptor {
	return []KindDescriptor{
		{Name: "uuid", Description: "RFC 4122 UUID v4"},
		{Name: "name", Description: "full person name"},
		{Name: "first_name", Description: "given name"},
		{Name: "last_name", Description: "family name"},
		{Name: "email", Description: "email address"},
		{Name: "username", Description: "username slug"},
		{Name: "phone", Description: "phone number"},
		{Name: "city", Description: "city name"},
		{Name: "country", Description: "country name"},
		{Name: "company", Description: "company name"},
		{Name: "url", Description: "URL"},
		{Name: "ipv4", Description: "IPv4 address"},
		{Name: "bool", Description: "boolean"},
		{Name: "int", Description: "integer in [min,max]; params: min, max"},
		{Name: "float", Description: "float in [min,max]; params: min, max"},
		{Name: "date", Description: "RFC 3339 timestamp"},
		{Name: "sequence", Description: "1-based row index"},
		{Name: "fixed", Description: "constant; params: value"},
	}
}
