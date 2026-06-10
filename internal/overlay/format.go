// Package overlay implements the small always-on-top status window. The
// window itself can only be built on Windows (see window.go), but the text
// formatting logic lives in this file so it can be unit tested on any
// platform.
package overlay

import (
	"fmt"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

const timeFormat = "15:04:05"

func formatProcessStatus(s model.Status) string {
	if !s.DBDRunning {
		return "DBD: not running"
	}
	return fmt.Sprintf("DBD: running (PID %d)", s.PID)
}

func formatServerAddress(serverIP string) string {
	if serverIP == "" {
		return "Server: -"
	}
	return "Server: " + serverIP
}

func formatPing(p model.PingStatus) string {
	if !p.Available {
		return "Ping: -"
	}
	return fmt.Sprintf("Ping: %d ms", p.AverageMs)
}

func formatPacketLoss(p model.PingStatus) string {
	if !p.PacketLossKnown {
		return "Loss: -"
	}
	return fmt.Sprintf("Loss: %d%%", p.PacketLossPercent)
}

func formatLastUpdate(t time.Time) string {
	if t.IsZero() {
		return "Updated: -"
	}
	return "Updated: " + t.Format(timeFormat)
}
