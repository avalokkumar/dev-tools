// Package pluginsdk is the public SDK that plugin authors import.
//
// Wire format: newline-delimited JSON over stdio (one JSON object per line).
// Each frame is a JSON-RPC 2.0 message. Plugins are spawned as child
// processes; the loader speaks RPC to them via stdin/stdout. Plugin stderr is
// forwarded to the loader's logger for diagnostics.
//
// Plugin authors normally use Serve(...) which handles the initialize/invoke
// methods automatically. Conn is exposed for advanced cases.
package pluginsdk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Frame is a single JSON-RPC 2.0 message.
type Frame struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *FrameError     `json:"error,omitempty"`
}

// FrameError carries a JSON-RPC error.
type FrameError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *FrameError) Error() string { return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message) }

// Conn is a bidirectional JSON-RPC framing wrapper. It is symmetric: both
// the loader (parent) and the plugin (child) use it.
type Conn struct {
	w  io.Writer
	wm sync.Mutex

	r       *bufio.Reader
	pending sync.Map // id-string → chan *Frame

	nextID atomic.Int64

	// inbound dispatches "request" frames (those carrying a Method) to a handler.
	handler func(ctx context.Context, f Frame) (json.RawMessage, *FrameError)

	closed atomic.Bool
}

// NewConn constructs a Conn over arbitrary streams.
// Pass nil handler if this side never accepts requests.
func NewConn(r io.Reader, w io.Writer, handler func(context.Context, Frame) (json.RawMessage, *FrameError)) *Conn {
	return &Conn{
		w:       w,
		r:       bufio.NewReader(r),
		handler: handler,
	}
}

// Serve reads frames until ctx is cancelled or the reader EOFs.
// Responses to outbound calls are routed to their pending channel; inbound
// requests are dispatched to the handler.
func (c *Conn) Serve(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		line, err := c.r.ReadBytes('\n')
		if len(line) == 0 && err == io.EOF {
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}
		var f Frame
		if uerr := json.Unmarshal(line, &f); uerr != nil {
			// Skip malformed frame; protocol-level error.
			continue
		}
		// Response (carries a result/error and an id we sent).
		if f.Method == "" && len(f.ID) > 0 {
			if chAny, ok := c.pending.LoadAndDelete(string(f.ID)); ok {
				ch := chAny.(chan *Frame)
				select {
				case ch <- &f:
				default:
				}
			}
			if err == io.EOF {
				return nil
			}
			continue
		}
		// Request from the peer.
		if c.handler != nil && f.Method != "" {
			go c.dispatch(ctx, f)
		}
		if err == io.EOF {
			return nil
		}
	}
}

func (c *Conn) dispatch(ctx context.Context, req Frame) {
	res, rerr := c.handler(ctx, req)
	if len(req.ID) == 0 {
		// Notification — no reply.
		return
	}
	resp := Frame{JSONRPC: "2.0", ID: req.ID}
	if rerr != nil {
		resp.Error = rerr
	} else {
		resp.Result = res
	}
	_ = c.write(resp)
}

// Call performs a request/response RPC and waits for the matching reply.
func (c *Conn) Call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	id := c.nextID.Add(1)
	idRaw, _ := json.Marshal(id)
	var paramsRaw json.RawMessage
	if params != nil {
		var err error
		paramsRaw, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("plugin: marshal params: %w", err)
		}
	}
	ch := make(chan *Frame, 1)
	c.pending.Store(string(idRaw), ch)
	defer c.pending.Delete(string(idRaw))

	if err := c.write(Frame{JSONRPC: "2.0", ID: idRaw, Method: method, Params: paramsRaw}); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case f := <-ch:
		if f.Error != nil {
			return nil, f.Error
		}
		return f.Result, nil
	}
}

// Notify sends a one-way notification (no id, no reply).
func (c *Conn) Notify(method string, params any) error {
	var paramsRaw json.RawMessage
	if params != nil {
		var err error
		paramsRaw, err = json.Marshal(params)
		if err != nil {
			return err
		}
	}
	return c.write(Frame{JSONRPC: "2.0", Method: method, Params: paramsRaw})
}

func (c *Conn) write(f Frame) error {
	if c.closed.Load() {
		return io.ErrClosedPipe
	}
	if f.JSONRPC == "" {
		f.JSONRPC = "2.0"
	}
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	c.wm.Lock()
	defer c.wm.Unlock()
	if _, err := c.w.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

// Close marks the conn as closed; further writes return ErrClosedPipe.
func (c *Conn) Close() { c.closed.Store(true) }
