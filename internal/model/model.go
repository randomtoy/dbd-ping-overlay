// Package model contains the shared data types passed between the
// process, netstat, ping, and overlay packages.
package model

import "time"

// PingStatus describes the outcome of a single round of ping measurements
// against a candidate server address.
type PingStatus struct {
	// Available reports whether an average latency value could be parsed
	// from the ping output.
	Available bool

	// AverageMs is the average round-trip time in milliseconds. It is only
	// meaningful when Available is true.
	AverageMs int

	// PacketLossKnown reports whether a packet loss percentage could be
	// parsed from the ping output.
	PacketLossKnown bool

	// PacketLossPercent is the percentage of lost packets. It is only
	// meaningful when PacketLossKnown is true.
	PacketLossPercent int

	// Message holds a human readable explanation when ping data is not
	// available, e.g. because ICMP traffic is blocked.
	Message string
}

// Status is a snapshot of the application state at a point in time. It is
// produced by the polling loop in internal/app and consumed by the overlay
// window.
type Status struct {
	// DBDRunning reports whether a Dead by Daylight process was found.
	DBDRunning bool

	// PID is the process ID of the detected Dead by Daylight process. It is
	// only meaningful when DBDRunning is true.
	PID int

	// ServerIP is the public IP address of the most likely game server, or
	// empty if none could be determined.
	ServerIP string

	// Ping holds the result of pinging ServerIP.
	Ping PingStatus

	// LastUpdate is the time at which this status was produced.
	LastUpdate time.Time

	// Message is a short human readable summary of the current state, used
	// to surface errors and informational messages in the overlay.
	Message string
}
