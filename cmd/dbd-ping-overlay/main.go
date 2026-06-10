// Command dbd-ping-overlay shows a small always-on-top window with the
// approximate ping to the Dead by Daylight game server, determined entirely
// from outside the game process.
package main

import (
	"os"

	"github.com/randomtoy/dbd-ping-overlay/internal/app"
	"github.com/randomtoy/dbd-ping-overlay/internal/config"
	"github.com/randomtoy/dbd-ping-overlay/internal/logging"
)

func main() {
	logger := logging.New(os.Stdout)

	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("starting dbd-ping-overlay",
		"process_name", cfg.ProcessName,
		"refresh_interval", cfg.RefreshInterval,
		"ping_count", cfg.PingCount,
		"ping_timeout", cfg.PingTimeout,
	)

	if err := app.New(cfg, logger).Run(); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}
