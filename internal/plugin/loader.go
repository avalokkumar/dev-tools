package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/pkg/pluginsdk"
)

// initializeTimeout caps how long the loader waits for a plugin to register.
const initializeTimeout = 5 * time.Second

// invokeTimeout caps a single Operation invocation.
const invokeTimeout = 30 * time.Second

// PluginOpDecl is the wire shape a plugin returns from the "initialize" RPC.
type PluginOpDecl struct {
	Tool        string          `json:"tool"`
	Op          string          `json:"op"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// initializeResult is the full reply to the initialize RPC.
type initializeResult struct {
	Plugin     string         `json:"plugin"`
	Operations []PluginOpDecl `json:"operations"`
}

// invokeParams are the parameters to the "invoke" RPC.
type invokeParams struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

// Plugin is a live, running plugin process attached to the Registry.
type Plugin struct {
	Manifest *Manifest
	Dir      string

	logger *log.Logger
	cmd    *exec.Cmd
	conn   *pluginsdk.Conn

	// exited closes when the supervisor goroutine has finished cmd.Wait().
	// Close() blocks on this channel rather than calling Wait() itself,
	// which would race the supervisor.
	exited chan struct{}

	mu      sync.Mutex
	opNames []string // names contributed to the registry
	dead    bool
	closed  bool
}

// LoadAll scans dir for subdirectories with a manifest.toml and starts a
// Plugin for each, registering its Operations with reg. Failures on a single
// plugin are logged but do not abort the others.
func LoadAll(ctx context.Context, dir string, reg *mcpserver.Registry, logger *log.Logger) ([]*Plugin, error) {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("plugin: scan %s: %w", dir, err)
	}
	var plugins []*Plugin
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pluginDir := filepath.Join(dir, e.Name())
		p, err := LoadOne(ctx, pluginDir, reg, logger)
		if err != nil {
			logger.Printf("plugin %s: load failed: %v", e.Name(), err)
			continue
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// LoadOptions tunes plugin loading. Zero-value is the production default.
type LoadOptions struct {
	// ExtraEnv is appended to the process environment.
	ExtraEnv []string
}

// LoadOne starts the plugin in pluginDir and registers its Operations.
func LoadOne(ctx context.Context, pluginDir string, reg *mcpserver.Registry, logger *log.Logger) (*Plugin, error) {
	return LoadOneWithOptions(ctx, pluginDir, reg, logger, LoadOptions{})
}

// LoadOneWithOptions is the options-explicit variant.
func LoadOneWithOptions(ctx context.Context, pluginDir string, reg *mcpserver.Registry, logger *log.Logger, opts LoadOptions) (*Plugin, error) {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	manifest, err := LoadManifest(pluginDir)
	if err != nil {
		return nil, err
	}
	bin := manifest.EntrypointPath(pluginDir)
	if _, err := os.Stat(bin); err != nil {
		return nil, fmt.Errorf("entrypoint %s: %w", bin, err)
	}

	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = pluginDir
	cmd.Env = append(os.Environ(), opts.ExtraEnv...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	// Forward plugin stderr verbatim.
	go forwardStderr(stderr, logger, manifest.Name)

	conn := pluginsdk.NewConn(stdout, stdin, nil)
	p := &Plugin{
		Manifest: manifest,
		Dir:      pluginDir,
		logger:   logger,
		cmd:      cmd,
		conn:     conn,
		exited:   make(chan struct{}),
	}

	// Run conn.Serve in the background; when it returns we wait on the
	// process and signal exit. This is the ONLY goroutine that calls
	// cmd.Wait() — Close() waits on p.exited instead.
	servCtx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = conn.Serve(servCtx)
		_ = cmd.Wait()
		p.markDead()
		close(p.exited)
		cancel()
	}()

	// Initialize.
	initCtx, initCancel := context.WithTimeout(ctx, initializeTimeout)
	defer initCancel()
	resRaw, err := conn.Call(initCtx, "initialize", map[string]any{"sdk": SDKVersion})
	if err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("initialize: %w", err)
	}
	var initRes initializeResult
	if err := json.Unmarshal(resRaw, &initRes); err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("initialize response: %w", err)
	}

	// Register every declared Operation; the handler routes to the plugin
	// over the same conn. Track op names for graceful shutdown.
	for _, decl := range initRes.Operations {
		decl := decl
		op := mcpserver.Operation{
			Tool:        decl.Tool,
			Op:          decl.Op,
			Description: decl.Description,
			InputSchema: decl.InputSchema,
			Handler:     p.makeHandler(decl),
		}
		if err := reg.Register(op); err != nil {
			logger.Printf("plugin %s: register %s: %v", manifest.Name, op.Name(), err)
			continue
		}
		p.opNames = append(p.opNames, op.Name())
	}
	logger.Printf("plugin %s: registered %d operations", manifest.Name, len(p.opNames))
	return p, nil
}

func (p *Plugin) makeHandler(decl PluginOpDecl) mcpserver.Handler {
	name := decl.Tool + "_" + decl.Op
	return func(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
		if p.IsDead() {
			return nil, fmt.Errorf("plugin %s is no longer running", p.Manifest.Name)
		}
		callCtx, cancel := context.WithTimeout(ctx, invokeTimeout)
		defer cancel()
		raw, err := p.conn.Call(callCtx, "invoke", invokeParams{Name: name, Args: args})
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
}

// IsDead reports whether the plugin process has exited.
func (p *Plugin) IsDead() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.dead
}

func (p *Plugin) markDead() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dead = true
}

// Close terminates the plugin process. Idempotent.
func (p *Plugin) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	p.conn.Close()
	_ = p.cmd.Process.Signal(os.Interrupt)

	// Wait for the supervisor goroutine to finish cmd.Wait(); fall back to
	// SIGKILL after a grace period.
	select {
	case <-p.exited:
	case <-time.After(2 * time.Second):
		_ = p.cmd.Process.Kill()
		<-p.exited
	}
	return nil
}

func forwardStderr(r io.Reader, logger *log.Logger, name string) {
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			logger.Printf("[plugin %s stderr] %s", name, buf[:n])
		}
		if err != nil {
			return
		}
	}
}
