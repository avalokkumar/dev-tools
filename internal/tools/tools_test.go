package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/devforge/devforge/internal/mcpserver"
)

// TestRegistry_SingleSourceOfTruth — B8: Register populates the registry; the
// same Operation list is what every Surface consumes. We assert presence of
// uuid_generate and uuid_hash, and that calling the handler round-trips.
func TestRegistry_SingleSourceOfTruth(t *testing.T) {
	t.Parallel()
	reg := mcpserver.NewRegistry()
	if err := Register(reg); err != nil {
		t.Fatalf("Register: %v", err)
	}

	gen, ok := reg.Get("uuid_generate")
	if !ok {
		t.Fatalf("uuid_generate not registered")
	}
	if gen.Path() != "/api/v1/uuid/generate" {
		t.Fatalf("Path = %q", gen.Path())
	}

	out, err := gen.Handler(context.Background(), json.RawMessage(`{"version":4,"count":2}`))
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	var got struct {
		Values []string `json:"values"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal: %v; raw=%s", err, out)
	}
	if len(got.Values) != 2 {
		t.Fatalf("Values=%v", got.Values)
	}

	if _, ok := reg.Get("uuid_hash"); !ok {
		t.Fatalf("uuid_hash not registered")
	}
	if reg.Len() < 2 {
		t.Fatalf("Len = %d", reg.Len())
	}
}
