package cli

import "io"

// Stdio bundles the streams a Surface invocation reads/writes.
// Tests inject buffers; production wires os.Stdin/Stdout/Stderr.
type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
