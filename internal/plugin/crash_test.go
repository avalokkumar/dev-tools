package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
)

// TestPlugin_Crash_ReturnsCleanError — B13: when the plugin exits during a
// call, the host returns a clean error instead of hanging or panicking.
func TestPlugin_Crash_ReturnsCleanError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	pluginDir := buildExamplePlugin(t)
	reg := mcpserver.NewRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p, err := LoadOneWithOptions(ctx, pluginDir, reg, log.New(os.Stderr, "[t] ", 0),
		LoadOptions{ExtraEnv: []string{"DEVFORGE_TEST_CRASH_ON_INVOKE=1"}})
	if err != nil {
		t.Fatalf("LoadOne: %v", err)
	}
	defer p.Close()

	op, _ := reg.Get("hello_say")
	callCtx, callCancel := context.WithTimeout(ctx, 3*time.Second)
	defer callCancel()
	_, err = op.Handler(callCtx, json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Allow the supervisor goroutine to observe the exit and mark dead.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && !p.IsDead() {
		time.Sleep(20 * time.Millisecond)
	}
	if !p.IsDead() {
		t.Fatalf("plugin should be marked dead after process exit")
	}

	// Subsequent calls return a clean error, not a hang.
	_, err = op.Handler(callCtx, json.RawMessage(`{}`))
	if err == nil {
		t.Fatalf("expected error after death")
	}
}

// TestPlugin_Timeout_ContextDeadlineExceeded — B13: a slow plugin call honors
// the caller's context deadline.
func TestPlugin_Timeout_ContextDeadlineExceeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	pluginDir := buildExamplePlugin(t)
	reg := mcpserver.NewRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p, err := LoadOneWithOptions(ctx, pluginDir, reg, log.New(os.Stderr, "[t] ", 0),
		LoadOptions{ExtraEnv: []string{"DEVFORGE_TEST_SLEEP_MS=2000"}})
	if err != nil {
		t.Fatalf("LoadOne: %v", err)
	}
	defer p.Close()

	op, _ := reg.Get("hello_say")
	callCtx, callCancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer callCancel()
	start := time.Now()
	_, err = op.Handler(callCtx, json.RawMessage(`{}`))
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "deadline") {
		t.Fatalf("expected deadline-exceeded, got %v", err)
	}
	if elapsed > 1500*time.Millisecond {
		t.Fatalf("call did not honor deadline (took %s)", elapsed)
	}
}
