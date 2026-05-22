// Package mcpserver hosts the Registry (single source of truth for Tools and
// Operations) and the thin adapter over the mcp-go SDK.
//
// Both the Web Surface and the MCP Surface consume the Registry; Adapters
// register Operations once at startup and never cross-call each other.
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// Handler is the executable side of an Operation. It receives the raw JSON
// arguments (already validated against the Operation's input schema by the
// caller) and returns the JSON-encoded Result.
type Handler func(ctx context.Context, args json.RawMessage) (json.RawMessage, error)

// Operation describes a single named action of a Tool.
//
// Naming convention:
//   - Tool is a short kebab/lowercase domain (e.g. "uuid", "json").
//   - Op  is the verb in lowercase (e.g. "generate", "format").
//   - Name (the wire id used by MCP and the HTTP path) is "<tool>_<op>".
type Operation struct {
	// Tool is the domain ("uuid", "jwt", ...).
	Tool string

	// Op is the verb ("generate", "decode", ...).
	Op string

	// Description is shown to humans and AI agents.
	Description string

	// InputSchema is the JSON Schema (draft 2020-12) describing the args.
	InputSchema json.RawMessage

	// Handler executes the operation.
	Handler Handler
}

// Name returns the wire identifier "<tool>_<op>".
func (o Operation) Name() string { return o.Tool + "_" + o.Op }

// Path returns the HTTP path "/api/v1/<tool>/<op>".
func (o Operation) Path() string { return "/api/v1/" + o.Tool + "/" + o.Op }

// Registry is an in-process catalog of Operations.
// It is safe for concurrent reads after the registration phase completes.
type Registry struct {
	mu  sync.RWMutex
	ops map[string]Operation
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{ops: make(map[string]Operation)}
}

// Register adds an Operation. Returns an error on duplicate name.
func (r *Registry) Register(op Operation) error {
	if op.Tool == "" || op.Op == "" {
		return fmt.Errorf("registry: Tool and Op are required")
	}
	if op.Handler == nil {
		return fmt.Errorf("registry: %q has no handler", op.Name())
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.ops[op.Name()]; exists {
		return fmt.Errorf("registry: duplicate operation %q", op.Name())
	}
	r.ops[op.Name()] = op
	return nil
}

// List returns Operations sorted by Name. Stable for tests and `tools/list`.
func (r *Registry) List() []Operation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Operation, 0, len(r.ops))
	for _, op := range r.ops {
		out = append(out, op)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}

// Get fetches an Operation by name.
func (r *Registry) Get(name string) (Operation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.ops[name]
	return op, ok
}

// Len returns the count of registered Operations.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.ops)
}
