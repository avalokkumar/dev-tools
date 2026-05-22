package tools

import (
	"errors"

	"github.com/devforge/devforge/pkg/engine"
)

// errExitOne signals "command finished, but report a non-zero exit code".
// We carry the exit code through cobra's RunE error path so tests see exit=1
// without messy os.Exit calls.
var errExitOne = exitError(1)

type exitError int

func (e exitError) Error() string  { return "command failed" }
func (e exitError) ExitCode() int  { return int(e) }
func (e exitError) Is(target error) bool {
	var t exitError
	if !errors.As(target, &t) {
		return false
	}
	return int(e) == int(t)
}

func hasErrorDiag(d []engine.Diagnostic) bool { return engine.HasError(d) }
