package webserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/tools"
)

// TestWebAPI_Uuid_Generate_200 — B5: POST /api/v1/uuid/generate returns 200
// + JSON body with `values` of requested length.
func TestWebAPI_Uuid_Generate_200(t *testing.T) {
	t.Parallel()
	reg := mcpserver.NewRegistry()
	if err := tools.Register(reg); err != nil {
		t.Fatalf("register tools: %v", err)
	}
	s, err := New(Config{Registry: reg})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := strings.NewReader(`{"version":4,"count":2}`)
	resp, err := http.Post(srv.URL+"/api/v1/uuid/generate", "application/json", body)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, raw)
	}
	var got struct {
		Values []string `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Values) != 2 {
		t.Fatalf("Values=%v", got.Values)
	}
}

// TestWebAPI_OperationsList — B5: /api/v1/operations enumerates registered ops.
func TestWebAPI_OperationsList(t *testing.T) {
	t.Parallel()
	reg := mcpserver.NewRegistry()
	if err := tools.Register(reg); err != nil {
		t.Fatalf("register tools: %v", err)
	}
	s, err := New(Config{Registry: reg})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/operations")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte(`"name":"uuid_generate"`)) {
		t.Fatalf("missing uuid_generate in: %s", body)
	}
}

// TestWebAPI_RejectsInvalidJSON — B5: malformed body → 400.
func TestWebAPI_RejectsInvalidJSON(t *testing.T) {
	t.Parallel()
	reg := mcpserver.NewRegistry()
	if err := tools.Register(reg); err != nil {
		t.Fatalf("register tools: %v", err)
	}
	s, _ := New(Config{Registry: reg})
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/v1/uuid/generate", "application/json",
		strings.NewReader("{not-json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status=%d", resp.StatusCode)
	}
}
