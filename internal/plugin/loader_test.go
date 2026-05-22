package plugin

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
)

// buildExamplePlugin compiles plugins/example-hello into the given dir as
// "plugin" and returns its parent dir (the plugin directory).
//
// The repository's plugins/example-hello uses a `replace` directive in its
// go.mod, so we build it in place and return its directory directly.
func buildExamplePlugin(t *testing.T) string {
	t.Helper()
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	pluginDir := filepath.Join(repoRoot, "plugins", "example-hello")
	bin := filepath.Join(pluginDir, "plugin")
	if _, err := os.Stat(bin); err == nil {
		return pluginDir
	}
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = pluginDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build example plugin: %v\n%s", err, out)
	}
	return pluginDir
}

// TestPluginLoader_Hello_Registered — B11: example-hello plugin loads and
// contributes one Operation to the Registry.
func TestPluginLoader_Hello_Registered(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	pluginDir := buildExamplePlugin(t)

	reg := mcpserver.NewRegistry()
	logger := log.New(os.Stderr, "[test] ", 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p, err := LoadOne(ctx, pluginDir, reg, logger)
	if err != nil {
		t.Fatalf("LoadOne: %v", err)
	}
	defer p.Close()

	if _, ok := reg.Get("hello_say"); !ok {
		ops := reg.List()
		names := make([]string, 0, len(ops))
		for _, o := range ops {
			names = append(names, o.Name())
		}
		t.Fatalf("hello_say not registered; registered=%v", names)
	}
}

// TestPlugin_Hello_ReachableViaMCP — B12: plugin Operation appears in MCP
// tools/list and tools/call returns its handler output.
func TestPlugin_Hello_ReachableViaMCP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	pluginDir := buildExamplePlugin(t)
	reg := mcpserver.NewRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	p, err := LoadOne(ctx, pluginDir, reg, log.New(os.Stderr, "[test] ", 0))
	if err != nil {
		t.Fatalf("LoadOne: %v", err)
	}
	defer p.Close()

	if _, ok := reg.Get("hello_say"); !ok {
		t.Fatalf("hello_say not in Registry")
	}
	// Build MCP server from this Registry.
	mcpsrv := mcpserver.New(reg)

	cliToSrv, srvIn := pipePair()
	srvOut, cliFromSrv := pipePair()
	go mcpsrv.Listen(ctx, cliToSrv, cliFromSrv)

	// initialize
	writeJSONLine(t, srvIn, map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "t", "version": "0"},
		},
	})
	_ = readJSONLine(t, srvOut)

	// tools/call
	writeJSONLine(t, srvIn, map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "tools/call",
		"params": map[string]any{
			"name":      "hello_say",
			"arguments": map[string]any{"name": "alok"},
		},
	})
	resp := readJSONLine(t, srvOut)
	result, _ := resp["result"].(map[string]any)
	content, _ := result["content"].([]any)
	if len(content) == 0 {
		t.Fatalf("no content: %+v", result)
	}
	first, _ := content[0].(map[string]any)
	text, _ := first["text"].(string)
	if !strings.Contains(text, "hello, alok") {
		t.Fatalf("missing greeting: %q", text)
	}

	_ = srvIn.Close()
	_ = cliFromSrv.Close()
}

// TestPlugin_Hello_Invoke — B11/B12: invoking hello_say through the registered
// Handler routes to the plugin process and returns its result.
func TestPlugin_Hello_Invoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	pluginDir := buildExamplePlugin(t)
	reg := mcpserver.NewRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	p, err := LoadOne(ctx, pluginDir, reg, log.New(os.Stderr, "[test] ", 0))
	if err != nil {
		t.Fatalf("LoadOne: %v", err)
	}
	defer p.Close()

	op, _ := reg.Get("hello_say")
	out, err := op.Handler(ctx, json.RawMessage(`{"name":"devforge"}`))
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}
	if !strings.Contains(string(out), "hello, devforge") {
		t.Fatalf("unexpected output: %s", out)
	}
}
