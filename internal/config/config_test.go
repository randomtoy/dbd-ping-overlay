package config

import (
	"testing"
	"time"
)

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse(nil)
	if err != nil {
		t.Fatalf("Parse(nil) returned error: %v", err)
	}

	if cfg.ProcessName != DefaultProcessName {
		t.Errorf("ProcessName = %q, want %q", cfg.ProcessName, DefaultProcessName)
	}
	if cfg.RefreshInterval != DefaultRefreshInterval {
		t.Errorf("RefreshInterval = %s, want %s", cfg.RefreshInterval, DefaultRefreshInterval)
	}
	if cfg.PingCount != DefaultPingCount {
		t.Errorf("PingCount = %d, want %d", cfg.PingCount, DefaultPingCount)
	}
	if cfg.PingTimeout != DefaultPingTimeout {
		t.Errorf("PingTimeout = %s, want %s", cfg.PingTimeout, DefaultPingTimeout)
	}
}

func TestParseOverrides(t *testing.T) {
	args := []string{
		"--process-name", "custom.exe",
		"--refresh-interval", "5s",
		"--ping-count", "2",
		"--ping-timeout", "500ms",
	}

	cfg, err := Parse(args)
	if err != nil {
		t.Fatalf("Parse(%v) returned error: %v", args, err)
	}

	if cfg.ProcessName != "custom.exe" {
		t.Errorf("ProcessName = %q, want %q", cfg.ProcessName, "custom.exe")
	}
	if cfg.RefreshInterval != 5*time.Second {
		t.Errorf("RefreshInterval = %s, want %s", cfg.RefreshInterval, 5*time.Second)
	}
	if cfg.PingCount != 2 {
		t.Errorf("PingCount = %d, want %d", cfg.PingCount, 2)
	}
	if cfg.PingTimeout != 500*time.Millisecond {
		t.Errorf("PingTimeout = %s, want %s", cfg.PingTimeout, 500*time.Millisecond)
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"empty process name", []string{"--process-name", ""}},
		{"zero refresh interval", []string{"--refresh-interval", "0s"}},
		{"negative ping count", []string{"--ping-count", "-1"}},
		{"zero ping timeout", []string{"--ping-timeout", "0s"}},
		{"unknown flag", []string{"--does-not-exist", "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args); err == nil {
				t.Errorf("Parse(%v) returned nil error, want error", tt.args)
			}
		})
	}
}
