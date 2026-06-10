package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/config"
	"github.com/randomtoy/dbd-ping-overlay/internal/model"
	"github.com/randomtoy/dbd-ping-overlay/internal/netstat"
)

type fakeLister struct {
	pids []int
	err  error
}

func (f fakeLister) ListPIDs(ctx context.Context, processName string) ([]int, error) {
	return f.pids, f.err
}

type fakeReader struct {
	conns []netstat.Connection
	err   error
}

func (f fakeReader) ListConnections(ctx context.Context) ([]netstat.Connection, error) {
	return f.conns, f.err
}

type fakePingRunner struct {
	status model.PingStatus
	err    error
}

func (f fakePingRunner) Ping(ctx context.Context, host string, count int, timeout time.Duration) (model.PingStatus, error) {
	return f.status, f.err
}

func newTestApp() *App {
	return &App{
		cfg:    config.Default(),
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),

		selector: netstat.DefaultSelector{},
	}
}

func TestPollProcessNotRunning(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: nil}

	status := a.poll(context.Background())

	if status.DBDRunning {
		t.Error("DBDRunning = true, want false")
	}
	if status.Message != msgProcessNotRunning {
		t.Errorf("Message = %q, want %q", status.Message, msgProcessNotRunning)
	}
}

func TestPollProcessLookupError(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{err: errors.New("tasklist not found")}

	status := a.poll(context.Background())

	if status.DBDRunning {
		t.Error("DBDRunning = true, want false")
	}
	if status.Message == "" {
		t.Error("Message is empty, want an error description")
	}
}

func TestPollNoServerDetected(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{conns: []netstat.Connection{
		{Proto: "TCP", RemoteAddr: "192.168.1.1:443", State: "ESTABLISHED", PID: 1234},
	}}

	status := a.poll(context.Background())

	if !status.DBDRunning || status.PID != 1234 {
		t.Errorf("DBDRunning/PID = %v/%d, want true/1234", status.DBDRunning, status.PID)
	}
	if status.ServerIP != "" {
		t.Errorf("ServerIP = %q, want empty", status.ServerIP)
	}
	if status.Message != msgNoServerDetected {
		t.Errorf("Message = %q, want %q", status.Message, msgNoServerDetected)
	}
}

func TestPollConnectionLookupError(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{err: errors.New("netstat not found")}

	status := a.poll(context.Background())

	if !status.DBDRunning {
		t.Error("DBDRunning = false, want true")
	}
	if status.Message == "" {
		t.Error("Message is empty, want an error description")
	}
}

func TestPollSuccess(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{conns: []netstat.Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1234},
	}}
	a.pingRunner = fakePingRunner{status: model.PingStatus{
		Available:         true,
		AverageMs:         38,
		PacketLossKnown:   true,
		PacketLossPercent: 0,
	}}

	status := a.poll(context.Background())

	if status.ServerIP != "203.0.113.10" {
		t.Errorf("ServerIP = %q, want %q", status.ServerIP, "203.0.113.10")
	}
	if status.Ping.AverageMs != 38 {
		t.Errorf("Ping.AverageMs = %d, want 38", status.Ping.AverageMs)
	}
	if status.Message != msgOK {
		t.Errorf("Message = %q, want %q", status.Message, msgOK)
	}
}

func TestPollPingUnavailable(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{conns: []netstat.Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1234},
	}}
	const wantMsg = "Ping unavailable / ICMP may be blocked"
	a.pingRunner = fakePingRunner{status: model.PingStatus{
		Available: false,
		Message:   wantMsg,
	}}

	status := a.poll(context.Background())

	if status.Ping.Available {
		t.Error("Ping.Available = true, want false")
	}
	if status.Message != wantMsg {
		t.Errorf("Message = %q, want %q", status.Message, wantMsg)
	}
}

func TestPollPingError(t *testing.T) {
	a := newTestApp()
	a.processLister = fakeLister{pids: []int{1234}}
	a.connReader = fakeReader{conns: []netstat.Connection{
		{Proto: "TCP", RemoteAddr: "203.0.113.10:7777", State: "ESTABLISHED", PID: 1234},
	}}
	a.pingRunner = fakePingRunner{err: errors.New("ping not found")}

	status := a.poll(context.Background())

	if status.ServerIP != "203.0.113.10" {
		t.Errorf("ServerIP = %q, want %q", status.ServerIP, "203.0.113.10")
	}
	if status.Message == "" {
		t.Error("Message is empty, want an error description")
	}
}
