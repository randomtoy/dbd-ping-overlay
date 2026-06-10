// Package config defines the runtime configuration for the application and
// how it is loaded from command line flags.
package config

import (
	"flag"
	"fmt"
	"time"
)

// Default values used when no flags are supplied.
const (
	DefaultProcessName     = "DeadByDaylight-Win64-Shipping.exe"
	DefaultRefreshInterval = 2 * time.Second
	DefaultPingCount       = 4
	DefaultPingTimeout     = time.Second
)

// Config holds all user-configurable settings for the application.
type Config struct {
	// ProcessName is the executable name to look for, e.g.
	// "DeadByDaylight-Win64-Shipping.exe".
	ProcessName string

	// RefreshInterval controls how often the overlay refreshes its
	// connection and ping data.
	RefreshInterval time.Duration

	// PingCount is the number of ICMP echo requests sent per ping check.
	PingCount int

	// PingTimeout is the per-reply timeout passed to the ping command.
	PingTimeout time.Duration
}

// Default returns a Config populated with sane defaults for an unconfigured
// run.
func Default() Config {
	return Config{
		ProcessName:     DefaultProcessName,
		RefreshInterval: DefaultRefreshInterval,
		PingCount:       DefaultPingCount,
		PingTimeout:     DefaultPingTimeout,
	}
}

// Parse builds a Config from the given command line arguments (typically
// os.Args[1:]), starting from the defaults returned by Default. Unset flags
// keep their default values.
func Parse(args []string) (Config, error) {
	cfg := Default()

	fs := flag.NewFlagSet("dbd-ping-overlay", flag.ContinueOnError)
	fs.StringVar(&cfg.ProcessName, "process-name", cfg.ProcessName,
		"Executable name of the game process to monitor")
	fs.DurationVar(&cfg.RefreshInterval, "refresh-interval", cfg.RefreshInterval,
		"How often to refresh connection and ping data (e.g. 2s)")
	fs.IntVar(&cfg.PingCount, "ping-count", cfg.PingCount,
		"Number of ICMP echo requests sent per ping check")
	fs.DurationVar(&cfg.PingTimeout, "ping-timeout", cfg.PingTimeout,
		"Per-reply timeout passed to the ping command (e.g. 1s)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.ProcessName == "" {
		return fmt.Errorf("process-name must not be empty")
	}
	if c.RefreshInterval <= 0 {
		return fmt.Errorf("refresh-interval must be positive, got %s", c.RefreshInterval)
	}
	if c.PingCount <= 0 {
		return fmt.Errorf("ping-count must be positive, got %d", c.PingCount)
	}
	if c.PingTimeout <= 0 {
		return fmt.Errorf("ping-timeout must be positive, got %s", c.PingTimeout)
	}
	return nil
}
