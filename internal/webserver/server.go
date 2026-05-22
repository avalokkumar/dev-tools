// Package webserver hosts the local web UI and the JSON HTTP API.
//
// The SPA is embedded into the binary via embed.FS (see assets.go); the
// HTTP API mirrors the Registry's Operations 1:1. Requests are routed by
// chi; responses are JSON. localhost-only by default.
package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/devforge/devforge/internal/mcpserver"
)

const apiCallTimeout = 30 * time.Second

// Config tunes the server. Zero-value Config yields safe defaults.
type Config struct {
	// Addr is the bind address. Default "127.0.0.1:0" (ephemeral port).
	Addr string
	// Registry is the source of truth for HTTP /api/v1/<tool>/<op> routes.
	// Optional: nil registry = no API routes (only /healthz + SPA).
	Registry *mcpserver.Registry
}

// Server is the HTTP entry point.
type Server struct {
	cfg     Config
	handler http.Handler
	srv     *http.Server
	addr    net.Addr
}

// New builds a Server. Safe to call without ListenAndServe (e.g. for httptest).
func New(cfg Config) (*Server, error) {
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:0"
	}
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	if cfg.Registry != nil {
		mountAPI(r, cfg.Registry)
	}

	dist, err := distFS()
	if err != nil {
		return nil, fmt.Errorf("webserver: load embedded assets: %w", err)
	}
	r.Handle("/*", spaHandler(dist))

	return &Server{cfg: cfg, handler: r}, nil
}

// Handler exposes the underlying http.Handler so tests can use httptest.
func (s *Server) Handler() http.Handler { return s.handler }

// ListenAndServe binds and serves until ctx is cancelled.
func (s *Server) ListenAndServe(ctx context.Context) error {
	l, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("webserver: listen: %w", err)
	}
	s.addr = l.Addr()
	s.srv = &http.Server{
		Handler:           s.handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- s.srv.Serve(l) }()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

// Addr returns the bound address. Empty until ListenAndServe binds.
func (s *Server) Addr() net.Addr { return s.addr }

// spaHandler serves embedded files; on miss, serves index.html inline so SPA routing works.
//
// We deliberately serve index.html as raw bytes (not via http.FileServer) for both
// the root and the SPA fallback. http.FileServer would otherwise issue 301 redirects
// from "/index.html" to "./", and from any directory-like path to its trailing-slash
// form — which the browser resolves relative to the original deep-link URL (e.g.
// /tools/), producing an infinite redirect loop.
func spaHandler(root fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(root))
	indexBytes, _ := fs.ReadFile(root, "index.html")
	serveIndex := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(indexBytes)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := trimLeadingSlash(r.URL.Path)
		if clean == "" || clean == "index.html" {
			serveIndex(w)
			return
		}
		if info, err := fs.Stat(root, clean); err != nil || info.IsDir() {
			serveIndex(w)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func trimLeadingSlash(p string) string {
	if len(p) > 0 && p[0] == '/' {
		return p[1:]
	}
	return p
}
