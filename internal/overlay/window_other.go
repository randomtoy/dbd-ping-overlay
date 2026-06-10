//go:build !windows

package overlay

import (
	"errors"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

// Window is a stub used when building for platforms other than Windows,
// where the native overlay window cannot be created. It exists so the rest
// of the codebase builds and can be tested cross-platform.
type Window struct{}

// New always fails on non-Windows platforms.
func New(onClose func()) (*Window, error) {
	return nil, errors.New("overlay window is only supported on Windows")
}

// Run is unreachable on non-Windows platforms because New always errors.
func (w *Window) Run() int { return 1 }

// Update is unreachable on non-Windows platforms because New always errors.
func (w *Window) Update(status model.Status) {}

// Close is unreachable on non-Windows platforms because New always errors.
func (w *Window) Close() {}
