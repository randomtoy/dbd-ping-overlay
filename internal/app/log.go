package app

import "github.com/randomtoy/dbd-ping-overlay/internal/model"

// logTransition emits log records only when something noteworthy changed
// since the previous poll, so the application does not spam the log on
// every refresh tick.
func (a *App) logTransition(prev, current model.Status, first bool) {
	if first || prev.DBDRunning != current.DBDRunning {
		if current.DBDRunning {
			a.logger.Info("game process found", "process", a.cfg.ProcessName, "pid", current.PID)
		} else {
			a.logger.Info("game process not running", "process", a.cfg.ProcessName)
		}
	}

	if current.ServerIP != "" && current.ServerIP != prev.ServerIP {
		a.logger.Info("server candidate selected", "ip", current.ServerIP)
	}

	if current.Message != prev.Message {
		a.logger.Info("status changed", "message", current.Message)
	}
}
