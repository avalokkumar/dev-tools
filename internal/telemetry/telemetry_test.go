package telemetry

import (
	"path/filepath"
	"testing"
)

// TestTelemetry_DisabledByDefault — D2: a fresh Recorder is off.
func TestTelemetry_DisabledByDefault(t *testing.T) {
	t.Parallel()
	r := NewWithPath(false, "")
	if r.Enabled() {
		t.Fatalf("expected disabled")
	}
	r.Track("a", nil)
	if err := r.Flush(); err != nil {
		t.Fatalf("Flush should be no-op: %v", err)
	}
}

// TestTelemetry_Optin_FlushesEvents — D2: when enabled, events persist.
func TestTelemetry_Optin_FlushesEvents(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	r := NewWithPath(true, path)
	if !r.Enabled() {
		t.Fatalf("expected enabled")
	}
	r.Track("uuid_generate", map[string]any{"version": 7})
	r.Track("json_format", nil)
	if err := r.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d events", len(got))
	}
	if got[0].Name != "uuid_generate" {
		t.Fatalf("name = %q", got[0].Name)
	}
	if v, _ := got[0].Props["version"].(float64); v != 7 {
		t.Fatalf("props = %v", got[0].Props)
	}
}

// TestTelemetry_FlushNoopWhenEmpty — D2: empty buffer + enabled = no file.
func TestTelemetry_FlushNoopWhenEmpty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	r := NewWithPath(true, path)
	if err := r.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 events")
	}
}
