package app

import (
	"context"
	"fmt"
	"time"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
	"github.com/randomtoy/dbd-ping-overlay/internal/netstat"
)

// externalCommandTimeout bounds how long a single "tasklist" or "netstat"
// invocation may take.
const externalCommandTimeout = 5 * time.Second

// pingOverhead is added on top of the configured ping count and per-reply
// timeout to give the "ping" command time to start and print its summary.
const pingOverhead = 5 * time.Second

const (
	msgProcessNotRunning = "Dead by Daylight is not running"
	msgNoServerDetected  = "no game server connection detected"
	msgOK                = "OK"
)

// poll gathers a single status snapshot: whether the configured process is
// running, which remote address looks like the game server, and the
// current ping to that address. Failures in any external command are
// captured in the returned status's Message field rather than returned as
// an error, so the overlay always has something to display.
func (a *App) poll(ctx context.Context) model.Status {
	status := model.Status{LastUpdate: time.Now()}

	pid, ok, err := a.findProcess(ctx)
	if err != nil {
		status.Message = fmt.Sprintf("process lookup failed: %v", err)
		return status
	}
	if !ok {
		status.Message = msgProcessNotRunning
		return status
	}
	status.DBDRunning = true
	status.PID = pid

	serverIP, ok, err := a.findServer(ctx, pid)
	if err != nil {
		status.Message = fmt.Sprintf("connection lookup failed: %v", err)
		return status
	}
	if !ok {
		status.Message = msgNoServerDetected
		return status
	}
	status.ServerIP = serverIP

	pingStatus, err := a.pingServer(ctx, serverIP)
	if err != nil {
		status.Message = fmt.Sprintf("ping failed: %v", err)
		return status
	}
	status.Ping = pingStatus

	if pingStatus.Available {
		status.Message = msgOK
	} else {
		status.Message = pingStatus.Message
	}

	return status
}

// findProcess looks up the configured process name and returns the first
// matching PID, if any.
func (a *App) findProcess(ctx context.Context) (pid int, found bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, externalCommandTimeout)
	defer cancel()

	pids, err := a.processLister.ListPIDs(ctx, a.cfg.ProcessName)
	if err != nil {
		return 0, false, err
	}
	if len(pids) == 0 {
		return 0, false, nil
	}

	return pids[0], true, nil
}

// findServer lists the connections owned by pid and asks the configured
// ServerSelector to pick the most likely game server address.
func (a *App) findServer(ctx context.Context, pid int) (ip string, found bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, externalCommandTimeout)
	defer cancel()

	conns, err := a.connReader.ListConnections(ctx)
	if err != nil {
		return "", false, err
	}

	candidate, ok := a.selector.Select(netstat.FilterByPID(conns, pid))
	if !ok {
		return "", false, nil
	}

	return candidate.IP, true, nil
}

// pingServer measures latency to host, bounding the operation by the
// configured ping count and per-reply timeout.
func (a *App) pingServer(ctx context.Context, host string) (model.PingStatus, error) {
	timeout := time.Duration(a.cfg.PingCount)*a.cfg.PingTimeout + pingOverhead

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return a.pingRunner.Ping(ctx, host, a.cfg.PingCount, a.cfg.PingTimeout)
}
