package main

// Smoke check that the package compiles and the binary entry point links.
// Functional CLI tests live in internal/cli; e2e tests spawn the real binary.

import "testing"

func TestMain_Compiles(t *testing.T) {
	// Building the test binary already proves main compiles.
	// This test exists so `go test ./...` covers cmd/devforge.
}
