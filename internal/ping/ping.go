// Package ping measures round-trip latency to a host using the operating
// system's "ping" command. It does not open raw sockets and requires no
// elevated privileges beyond what "ping" itself needs.
package ping

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

// Runner measures round-trip latency to host. Implementations must not
// block indefinitely; ctx should be used to bound the operation.
type Runner interface {
	Ping(ctx context.Context, host string, count int, timeout time.Duration) (model.PingStatus, error)
}

// SystemPing implements Runner using the platform "ping" command.
type SystemPing struct{}

// Ping runs:
//
//	ping -n <count> -w <timeout-ms> <host>
//
// and parses the result. A non-zero exit status from ping (e.g. because all
// requests timed out) is not treated as an error as long as ping produced
// output; in that case the returned status reflects unavailable ICMP rather
// than a Go-level error. An error is only returned if the ping command
// itself could not be run.
func (SystemPing) Ping(ctx context.Context, host string, count int, timeout time.Duration) (model.PingStatus, error) {
	timeoutMs := int(timeout / time.Millisecond)
	if timeoutMs < 1 {
		timeoutMs = 1
	}

	cmd := exec.CommandContext(ctx, "ping",
		"-n", strconv.Itoa(count),
		"-w", strconv.Itoa(timeoutMs),
		host,
	)

	out, runErr := cmd.CombinedOutput()
	if len(out) == 0 && runErr != nil {
		return model.PingStatus{}, fmt.Errorf("run ping: %w", runErr)
	}

	return parsePingOutput(string(out)), nil
}
