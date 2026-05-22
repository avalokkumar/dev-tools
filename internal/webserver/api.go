package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devforge/devforge/internal/mcpserver"
)

// maxRequestBytes caps the body size accepted on /api/v1/* routes.
const maxRequestBytes = 16 << 20 // 16 MiB

// mountAPI wires every Operation in the Registry onto the router under
// /api/v1/<tool>/<op>. Each route accepts JSON args, calls the Operation's
// Handler, and returns the JSON Result.
func mountAPI(r chi.Router, reg *mcpserver.Registry) {
	for _, op := range reg.List() {
		op := op
		r.Post(op.Path(), makeAPIHandler(op))
	}
	// /api/v1/operations lists registered operations for client discovery.
	r.Get("/api/v1/operations", func(w http.ResponseWriter, _ *http.Request) {
		type opInfo struct {
			Tool        string          `json:"tool"`
			Op          string          `json:"op"`
			Name        string          `json:"name"`
			Path        string          `json:"path"`
			Description string          `json:"description"`
			InputSchema json.RawMessage `json:"inputSchema,omitempty"`
		}
		out := make([]opInfo, 0)
		for _, op := range reg.List() {
			out = append(out, opInfo{
				Tool: op.Tool, Op: op.Op, Name: op.Name(),
				Path: op.Path(), Description: op.Description, InputSchema: op.InputSchema,
			})
		}
		writeJSON(w, http.StatusOK, out)
	})
}

func makeAPIHandler(op mcpserver.Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(http.MaxBytesReader(w, req.Body, maxRequestBytes))
		if err != nil {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{
				"error": "request body too large",
			})
			return
		}
		if len(body) == 0 {
			body = []byte("{}")
		}
		// Ensure body is valid JSON before invoking the handler.
		if !json.Valid(body) {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "request body must be valid JSON",
			})
			return
		}
		ctx, cancel := context.WithTimeout(req.Context(), apiCallTimeout)
		defer cancel()
		raw, err := op.Handler(ctx, body)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				writeJSON(w, http.StatusGatewayTimeout, map[string]string{"error": "operation timed out"})
				return
			}
			// json.Unmarshal errors from adapter args are user-input issues → 400.
			var typeErr *json.UnmarshalTypeError
			var syntaxErr *json.SyntaxError
			if errors.As(err, &typeErr) || errors.As(err, &syntaxErr) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(raw)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
