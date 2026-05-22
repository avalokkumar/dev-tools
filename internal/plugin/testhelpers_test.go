package plugin

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func pipePair() (*io.PipeReader, *io.PipeWriter) { return io.Pipe() }

func writeJSONLine(t *testing.T, w io.Writer, msg any) {
	t.Helper()
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, err := w.Write(append(b, '\n')); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readJSONLine(t *testing.T, r io.Reader) map[string]any {
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
