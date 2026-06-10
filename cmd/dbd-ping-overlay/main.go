// Command dbd-ping-overlay shows a small always-on-top window with the
// approximate ping to the Dead by Daylight game server, determined entirely
// from outside the game process.
package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/randomtoy/dbd-ping-overlay/internal/app"
	"github.com/randomtoy/dbd-ping-overlay/internal/config"
	"github.com/randomtoy/dbd-ping-overlay/internal/logging"
)

func main() {
	logWriter, closeLog := openLogWriter()
	defer closeLog()

	logger := logging.New(logWriter)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", "value", r, "stack", string(debug.Stack()))
			os.Exit(1)
		}
	}()

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

// openLogWriter returns a writer that duplicates log output to standard
// output and to a log file next to the executable. The release binary is
// built with -H windowsgui and has no console, so the log file is the only
// way to see startup errors. If the file cannot be created, logging falls
// back to standard output only.
func openLogWriter() (io.Writer, func()) {
	exe, err := os.Executable()
	if err != nil {
		return os.Stdout, func() {}
	}

	path := filepath.Join(filepath.Dir(exe), "dbd-ping-overlay.log")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return os.Stdout, func() {}
	}

	return fanoutWriter{os.Stdout, f}, func() { f.Close() }
}

// fanoutWriter writes to every underlying writer, ignoring errors from
// individual writers so that a broken stdout (as in a -H windowsgui binary
// with no console) cannot prevent the log file from receiving output.
type fanoutWriter []io.Writer

func (fw fanoutWriter) Write(p []byte) (int, error) {
	for _, w := range fw {
		w.Write(p)
	}
	return len(p), nil
}
