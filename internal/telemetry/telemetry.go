// Package telemetry buffers anonymized usage events to a local file.
//
// Off by default. Opt-in via either:
//   - environment: DEVFORGE_TELEMETRY=1
//   - config:      ~/.devforge/telemetry.enabled (any non-empty file)
//
// Events are written as newline-delimited JSON; never auto-uploaded by this
// package. A future Phase D enhancement will add a flush-to-remote endpoint.
package telemetry

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Event is one telemetry record.
type Event struct {
	Timestamp time.Time      `json:"ts"`
	Name      string         `json:"name"`
	Props     map[string]any `json:"props,omitempty"`
	OS        string         `json:"os"`
	Arch      string         `json:"arch"`
}

// Recorder appends events to disk when enabled, otherwise no-ops.
type Recorder struct {
	enabled bool
	path    string

	mu  sync.Mutex
	buf []Event
}

// New returns a Recorder. Call Enabled() to decide whether to bother
// constructing event Props.
func New() *Recorder {
	return &Recorder{
		enabled: detectEnabled(),
		path:    defaultPath(),
	}
}

// NewWithPath is the test seam.
func NewWithPath(enabled bool, path string) *Recorder {
	return &Recorder{enabled: enabled, path: path}
}

// Enabled reports whether telemetry is on.
func (r *Recorder) Enabled() bool { return r.enabled }

// Track records an event. Cheap when disabled.
func (r *Recorder) Track(name string, props map[string]any) {
	if !r.enabled {
		return
	}
	r.mu.Lock()
	r.buf = append(r.buf, Event{
		Timestamp: time.Now().UTC(),
		Name:      name,
		Props:     props,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	})
	r.mu.Unlock()
}

// Flush writes the buffered events to the configured path and clears the
// buffer. Safe to call when disabled (no-op).
func (r *Recorder) Flush() error {
	if !r.enabled {
		return nil
	}
	r.mu.Lock()
	pending := r.buf
	r.buf = nil
	r.mu.Unlock()
	if len(pending) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range pending {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

// Read returns all events stored on disk. Useful for tests and `devforge
// telemetry show`.
func Read(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []Event
	dec := json.NewDecoder(f)
	for {
		var e Event
		if err := dec.Decode(&e); err != nil {
			if errors.Is(err, io.EOF) {
				return out, nil
			}
			return out, err
		}
		out = append(out, e)
	}
}

func detectEnabled() bool {
	if os.Getenv("DEVFORGE_TELEMETRY") == "1" {
		return true
	}
	flag := flagPath()
	if flag == "" {
		return false
	}
	if info, err := os.Stat(flag); err == nil && info.Size() > 0 {
		return true
	}
	return false
}

func flagPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".devforge", "telemetry.enabled")
}

func defaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "devforge-telemetry.jsonl"
	}
	return filepath.Join(home, ".devforge", "events.jsonl")
}
