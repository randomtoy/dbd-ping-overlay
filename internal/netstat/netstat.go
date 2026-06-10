// Package netstat inspects the operating system's connection table to find
// the remote endpoints a process is talking to. It only reads information
// that is already exposed by the "netstat" command and never touches the
// target process itself.
package netstat

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Connection is a single row from "netstat -ano".
type Connection struct {
	// Proto is "TCP" or "UDP".
	Proto string

	// LocalAddr is the local endpoint in "host:port" form.
	LocalAddr string

	// RemoteAddr is the remote endpoint in "host:port" form. For UDP
	// sockets that have not connected to a specific peer, this is "*:*".
	RemoteAddr string

	// State is the TCP connection state (e.g. "ESTABLISHED",
	// "LISTENING"). It is always empty for UDP, which is stateless.
	State string

	// PID is the process ID that owns the connection.
	PID int
}

// Reader lists the current network connections known to the operating
// system. It is an interface so that the implementation can later be
// replaced or augmented, e.g. with a passive packet capture backend.
type Reader interface {
	ListConnections(ctx context.Context) ([]Connection, error)
}

// NetstatReader implements Reader using the "netstat -ano" command.
type NetstatReader struct{}

// ListConnections runs "netstat -ano" and parses its output.
func (NetstatReader) ListConnections(ctx context.Context) ([]Connection, error) {
	cmd := exec.CommandContext(ctx, "netstat", "-ano")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run netstat: %w", err)
	}

	return parseNetstatOutput(string(out))
}

// parseNetstatOutput parses the output of "netstat -ano" into a list of
// connections. Lines that do not describe a TCP or UDP connection (headers,
// blank lines, etc.) are skipped.
func parseNetstatOutput(output string) ([]Connection, error) {
	var conns []Connection

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}

		switch strings.ToUpper(fields[0]) {
		case "TCP":
			// Proto, Local Address, Foreign Address, State, PID
			if len(fields) != 5 {
				continue
			}
			pid, err := strconv.Atoi(fields[4])
			if err != nil {
				continue
			}
			conns = append(conns, Connection{
				Proto:      "TCP",
				LocalAddr:  fields[1],
				RemoteAddr: fields[2],
				State:      fields[3],
				PID:        pid,
			})

		case "UDP":
			// Proto, Local Address, Foreign Address, PID (no State column)
			if len(fields) != 4 {
				continue
			}
			pid, err := strconv.Atoi(fields[3])
			if err != nil {
				continue
			}
			conns = append(conns, Connection{
				Proto:      "UDP",
				LocalAddr:  fields[1],
				RemoteAddr: fields[2],
				PID:        pid,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan netstat output: %w", err)
	}

	return conns, nil
}

// FilterByPID returns the subset of conns owned by the given process ID.
func FilterByPID(conns []Connection, pid int) []Connection {
	var result []Connection
	for _, c := range conns {
		if c.PID == pid {
			result = append(result, c)
		}
	}
	return result
}
