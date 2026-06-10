package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
	"github.com/randomtoy/dbd-ping-overlay/internal/netstat"
)

// fakeWindow records every status it is asked to display and lets the test
// trigger window close on demand.
type fakeWindow struct {
	mu      sync.Mutex
	updates []model.Status
	onClose func()
	closed  chan struct{}
}

func newFakeWindow(onClose func()) *fakeWindow {
	return &fakeWindow{onClose: onClose, closed: make(chan struct{})}
}

func (w *fakeWindow) Run() int {
	<-w.closed
	return 0
}

func (w *fakeWindow) Update(status model.Status) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.updates = append(w.updates, status)
}

func (w *fakeWindow) Close() {
	select {
	case <-w.closed:
	default:
		close(w.closed)
	}
	if w.onClose != nil {
		w.onClose()
	}
}

func (w *fakeWindow) updateCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.updates)
}

func TestAppRunUpdatesWindowAndStopsOnClose(t *testing.T) {
	a := newTestApp()
	a.cfg.RefreshInterval = 10 * time.Millisecond
	a.processLister = fakeLister{pids: nil}
	a.connReader = fakeReader{}
	a.pingRunner = fakePingRunner{}

	var win *fakeWindow
	a.newWindow = func(onClose func()) (overlayWindow, error) {
		win = newFakeWindow(onClose)
		return win, nil
	}

	done := make(chan struct{})
	go func() {
		if err := a.Run(); err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
		close(done)
	}()

	// Wait for at least one refresh, then close the window as the user
	// would.
	deadline := time.After(time.Second)
	for win == nil || win.updateCount() == 0 {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for first window update")
		case <-time.After(time.Millisecond):
		}
	}

	win.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run() did not return after window close")
	}

	if win.updateCount() == 0 {
		t.Error("window was never updated")
	}
}

func TestPollLoopStopsOnContextCancel(t *testing.T) {
	a := newTestApp()
	a.cfg.RefreshInterval = 5 * time.Millisecond
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{conns: []netstat.Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1234},
	}}
	a.pingRunner = fakePingRunner{status: model.PingStatus{Available: true, AverageMs: 20}}

	ctx, cancel := context.WithCancel(context.Background())
	win := newFakeWindow(nil)

	loopDone := make(chan struct{})
	go func() {
		a.pollLoop(ctx, win)
		close(loopDone)
	}()

	deadline := time.After(time.Second)
	for win.updateCount() < 2 {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for repeated updates")
		case <-time.After(time.Millisecond):
		}
	}

	cancel()

	select {
	case <-loopDone:
	case <-time.After(time.Second):
		t.Fatal("pollLoop did not stop after context cancel")
	}

	last := win.updates[len(win.updates)-1]
	if last.ServerIP != "203.0.113.10" {
		t.Errorf("last update ServerIP = %q, want %q", last.ServerIP, "203.0.113.10")
	}
}
