package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestManifest_Valid — B10: a minimal valid manifest parses.
func TestManifest_Valid(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.toml"), []byte(`
sdk = 1
name = "hello"
version = "0.1.0"
description = "demo plugin"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if m.Name != "hello" {
		t.Fatalf("Name = %q", m.Name)
	}
	if !strings.HasSuffix(m.EntrypointPath(dir), "plugin") {
		t.Fatalf("EntrypointPath = %q", m.EntrypointPath(dir))
	}
}

// TestManifest_RejectVersionMismatch — B10: SDK mismatch fails validation.
func TestManifest_RejectVersionMismatch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.toml"), []byte(`
sdk = 999
name = "future"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadManifest(dir)
	if err == nil || !strings.Contains(err.Error(), "sdk version") {
		t.Fatalf("expected sdk version error, got %v", err)
	}
}

// TestManifest_MissingName — B10: name is required.
func TestManifest_MissingName(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.toml"), []byte(`sdk = 1`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadManifest(dir); err == nil {
		t.Fatalf("expected error")
	}
}
