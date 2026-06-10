package overlay

import (
	"testing"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

func TestFormatProcessStatus(t *testing.T) {
	if got, want := formatProcessStatus(model.Status{DBDRunning: false}), "DBD: not running"; got != want {
		t.Errorf("formatProcessStatus() = %q, want %q", got, want)
	}

	if got, want := formatProcessStatus(model.Status{DBDRunning: true, PID: 4242}), "DBD: running (PID 4242)"; got != want {
		t.Errorf("formatProcessStatus() = %q, want %q", got, want)
	}
}

func TestFormatServerAddress(t *testing.T) {
	if got, want := formatServerAddress(""), "Server: -"; got != want {
		t.Errorf("formatServerAddress(\"\") = %q, want %q", got, want)
	}

	if got, want := formatServerAddress("203.0.113.10"), "Server: 203.0.113.10"; got != want {
		t.Errorf("formatServerAddress() = %q, want %q", got, want)
	}
}

func TestFormatPing(t *testing.T) {
	if got, want := formatPing(model.PingStatus{Available: false}), "Ping: -"; got != want {
		t.Errorf("formatPing() = %q, want %q", got, want)
	}

	if got, want := formatPing(model.PingStatus{Available: true, AverageMs: 38}), "Ping: 38 ms"; got != want {
		t.Errorf("formatPing() = %q, want %q", got, want)
	}
}

func TestFormatPacketLoss(t *testing.T) {
	if got, want := formatPacketLoss(model.PingStatus{PacketLossKnown: false}), "Loss: -"; got != want {
		t.Errorf("formatPacketLoss() = %q, want %q", got, want)
	}

	if got, want := formatPacketLoss(model.PingStatus{PacketLossKnown: true, PacketLossPercent: 5}), "Loss: 5%"; got != want {
		t.Errorf("formatPacketLoss() = %q, want %q", got, want)
	}
}

func TestFormatLastUpdate(t *testing.T) {
	if got, want := formatLastUpdate(time.Time{}), "Updated: -"; got != want {
		t.Errorf("formatLastUpdate(zero) = %q, want %q", got, want)
	}

	ts := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	if got, want := formatLastUpdate(ts), "Updated: 15:04:05"; got != want {
		t.Errorf("formatLastUpdate() = %q, want %q", got, want)
	}
}
