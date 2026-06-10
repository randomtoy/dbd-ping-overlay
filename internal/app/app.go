// Package app wires together process discovery, connection inspection, and
// ping measurements, and drives the overlay window.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/config"
	"github.com/randomtoy/dbd-ping-overlay/internal/model"
	"github.com/randomtoy/dbd-ping-overlay/internal/netstat"
	"github.com/randomtoy/dbd-ping-overlay/internal/overlay"
	"github.com/randomtoy/dbd-ping-overlay/internal/ping"
	"github.com/randomtoy/dbd-ping-overlay/internal/process"
)

// overlayWindow is the subset of *overlay.Window used by App. It exists so
// tests can run the polling loop against a fake window instead of creating
// a real one.
type overlayWindow interface {
	Run() int
	Update(model.Status)
	Close()
}

// App polls the system for the game process and its server connection on a
// fixed interval and reports the result through an overlay window.
type App struct {
	cfg    config.Config
	logger *slog.Logger

	processLister process.Lister
	connReader    netstat.Reader
	selector      netstat.ServerSelector
	pingRunner    ping.Runner

	newWindow func(onClose func()) (overlayWindow, error)
}

// New creates an App using the default backends: "tasklist" for process
// discovery, "netstat" for connections, the system "ping" command, and the
// native overlay window.
func New(cfg config.Config, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,

		processLister: process.TasklistLister{},
		connReader:    netstat.NetstatReader{},
		selector:      netstat.DefaultSelector{},
		pingRunner:    ping.SystemPing{},

		newWindow: func(onClose func()) (overlayWindow, error) {
			return overlay.New(onClose)
		},
	}
}

// Run creates the overlay window, starts the background polling loop, and
// blocks until the window is closed.
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	win, err := a.newWindow(cancel)
	if err != nil {
		return fmt.Errorf("create overlay window: %w", err)
	}

	go a.pollLoop(ctx, win)

	win.Run()
	return nil
}

// pollLoop refreshes the status immediately and then on every
// RefreshInterval until ctx is canceled.
func (a *App) pollLoop(ctx context.Context, win overlayWindow) {
	var prev model.Status
	first := true

	refresh := func() {
		current := a.poll(ctx)
		a.logTransition(prev, current, first)
		win.Update(current)
		prev = current
		first = false
	}

	refresh()

	ticker := time.NewTicker(a.cfg.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			refresh()
		}
	}
}
