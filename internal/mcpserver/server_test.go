package mcpserver

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

// TestMCP_UuidGenerate_Roundtrip — B7: register a uuid_generate Operation
// directly with the Registry and exercise the full initialize → tools/list →
// tools/call sequence over the in-memory pipe pair.
func TestMCP_UuidGenerate_Roundtrip(t *testing.T) {
	t.Parallel()
	reg := NewRegistry()
	if err := reg.Register(Operation{
		Tool:        "uuid",
		Op:          "generate",
		Description: "test",
		InputSchema: []byte(`{"type":"object"}`),
		Handler: func(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{"values":["fixed-value"]}`), nil
		},
	}); err != nil {
		t.Fatalf("register: %v", err)
	}
	srv := New(reg)

	cliToSrv, srvIn := io.Pipe()
	srvOut, cliFromSrv := io.Pipe()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Listen(ctx, cliToSrv, cliFromSrv) }()

	// initialize
	mustWriteJSON(t, srvIn, map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "t", "version": "0"},
		},
	})
	_ = mustReadJSON(t, srvOut)

	// tools/list
	mustWriteJSON(t, srvIn, map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "tools/list",
	})
	listResp := mustReadJSON(t, srvOut)
	listResult, _ := listResp["result"].(map[string]any)
	toolsArr, _ := listResult["tools"].([]any)
	if len(toolsArr) != 1 {
		t.Fatalf("tools/list count = %d, want 1", len(toolsArr))
	}

	// tools/call
	mustWriteJSON(t, srvIn, map[string]any{
		"jsonrpc": "2.0", "id": 3, "method": "tools/call",
		"params": map[string]any{
			"name":      "uuid_generate",
			"arguments": map[string]any{},
		},
	})
	callResp := mustReadJSON(t, srvOut)
	callResult, _ := callResp["result"].(map[string]any)
	content, _ := callResult["content"].([]any)
	if len(content) == 0 {
		t.Fatalf("empty tools/call content: %+v", callResult)
	}
	first, _ := content[0].(map[string]any)
	text, _ := first["text"].(string)
	if !strings.Contains(text, "fixed-value") {
		t.Fatalf("missing handler output in: %q", text)
	}

	_ = srvIn.Close()
	_ = cliFromSrv.Close()
	cancel()
	<-errCh
}

// TestMcpServe_HandshakeOnly — A4: an empty-registry server completes the
// MCP `initialize` handshake and reports its serverInfo.
func TestMcpServe_HandshakeOnly(t *testing.T) {
	t.Parallel()

	reg := NewRegistry()
	s := New(reg)

	clientToServer, serverIn := io.Pipe()
	serverOut, clientFromServer := io.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Listen(ctx, clientToServer, clientFromServer)
	}()

	// Send `initialize` request.
	initReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "devforge-test",
				"version": "0.0.0",
			},
		},
	}
	mustWriteJSON(t, serverIn, initReq)

	// Read response line.
	resp := mustReadJSON(t, serverOut)
	result, ok := resp["result"].(map[string]any)
	if !ok {
		t.Fatalf("missing result field: %+v", resp)
	}
	srvInfo, ok := result["serverInfo"].(map[string]any)
	if !ok {
		t.Fatalf("missing serverInfo: %+v", result)
	}
	if name, _ := srvInfo["name"].(string); name != "devforge" {
		t.Fatalf("serverInfo.name = %q, want \"devforge\"", name)
	}

	// Cleanly close streams to let Listen return.
	_ = serverIn.Close()
	_ = clientFromServer.Close()
	cancel()
	select {
	case <-errCh:
	case <-time.After(2 * time.Second):
		t.Fatalf("server did not return after close")
	}
}

func mustWriteJSON(t *testing.T, w io.Writer, msg any) {
	t.Helper()
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, err := w.Write(append(b, '\n')); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func mustReadJSON(t *testing.T, r io.Reader) map[string]any {
	t.Helper()
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Fatalf("read: %v", err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		t.Fatalf("empty response")
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("unmarshal %q: %v", line, err)
	}
	return m
}
