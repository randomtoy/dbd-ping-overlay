//go:build windows

package overlay

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

const (
	windowWidth  = 300
	windowHeight = 140
	windowTitle  = "DBD Ping Overlay"
)

// Window is a small always-on-top status window showing whether Dead by
// Daylight is running and the latest ping measurement to its server.
type Window struct {
	mw *walk.MainWindow

	processLabel *walk.Label
	serverLabel  *walk.Label
	pingLabel    *walk.Label
	lossLabel    *walk.Label
	updatedLabel *walk.Label
	messageLabel *walk.Label
}

// New creates and lays out the overlay window without showing it. Call Run
// to display it and start the message loop.
//
// onClose, if non-nil, is invoked once when the user closes the window.
// Callers typically use it to cancel a context and stop background work.
func New(onClose func()) (*Window, error) {
	w := &Window{}

	err := MainWindow{
		AssignTo: &w.mw,
		Title:    windowTitle,
		Bounds:   Rectangle{Width: windowWidth, Height: windowHeight},
		Layout:   VBox{},
		Children: []Widget{
			Label{AssignTo: &w.processLabel, Text: formatProcessStatus(model.Status{})},
			Label{AssignTo: &w.serverLabel, Text: formatServerAddress("")},
			Label{AssignTo: &w.pingLabel, Text: formatPing(model.PingStatus{})},
			Label{AssignTo: &w.lossLabel, Text: formatPacketLoss(model.PingStatus{})},
			Label{AssignTo: &w.updatedLabel, Text: formatLastUpdate(model.Status{}.LastUpdate)},
			Label{AssignTo: &w.messageLabel, Text: "Starting..."},
		},
	}.Create()
	if err != nil {
		return nil, fmt.Errorf("create overlay window: %w", err)
	}

	// Best effort: keep the window above other applications, including
	// Dead by Daylight running in borderless/windowed mode. Failure here is
	// not fatal, the window is still fully usable.
	win.SetWindowPos(w.mw.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)

	if onClose != nil {
		w.mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			onClose()
		})
	}

	return w, nil
}

// Run displays the window and blocks until it is closed.
func (w *Window) Run() int {
	return w.mw.Run()
}

// Update refreshes the displayed status. It is safe to call from any
// goroutine.
func (w *Window) Update(status model.Status) {
	w.mw.Synchronize(func() {
		w.processLabel.SetText(formatProcessStatus(status))
		w.serverLabel.SetText(formatServerAddress(status.ServerIP))
		w.pingLabel.SetText(formatPing(status.Ping))
		w.lossLabel.SetText(formatPacketLoss(status.Ping))
		w.updatedLabel.SetText(formatLastUpdate(status.LastUpdate))
		w.messageLabel.SetText(status.Message)
	})
}

// Close closes the window programmatically, e.g. during shutdown triggered
// from outside the UI.
func (w *Window) Close() {
	w.mw.Synchronize(func() {
		w.mw.Close()
	})
}
