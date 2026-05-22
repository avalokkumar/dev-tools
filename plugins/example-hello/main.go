// Command plugin (example-hello) demonstrates the DevForge plugin SDK.
//
// It contributes one Tool ("hello") with one Operation ("say"), which echoes
// a greeting back to the caller.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/devforge/devforge/pkg/pluginsdk"
)

func main() {
	if err := pluginsdk.Serve(pluginsdk.Plugin{
		Name:    "example-hello",
		Version: "0.1.0",
		Operations: []pluginsdk.OpDecl{
			{
				Tool:        "hello",
				Op:          "say",
				Description: "Greet the caller by name.",
				InputSchema: json.RawMessage(`{
  "type": "object",
  "properties": {"name": {"type": "string", "default": "world"}}
}`),
				Handler: handleSay,
			},
		},
	}); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("plugin: %v", err)
	}
}

type sayArgs struct {
	Name string `json:"name"`
}

func handleSay(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
	var a sayArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return nil, err
		}
	}
	if a.Name == "" {
		a.Name = "world"
	}
	// Test hooks (B13). Off in production unless these env vars are set.
	if os.Getenv("DEVFORGE_TEST_CRASH_ON_INVOKE") == "1" {
		os.Exit(1)
	}
	if v := os.Getenv("DEVFORGE_TEST_SLEEP_MS"); v != "" {
		ms, _ := strconv.Atoi(v)
		select {
		case <-time.After(time.Duration(ms) * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return json.Marshal(map[string]string{
		"message": fmt.Sprintf("hello, %s — from example-hello plugin", a.Name),
	})
}
