package netstat

import "strings"

// Score weights used by DefaultSelector. Established TCP connections are
// the strongest signal that an address is an active game server; UDP
// connections (commonly used for game traffic) are the next best signal.
const (
	scoreTCPEstablished = 3
	scoreUDP            = 2
)

// ServerCandidate is a possible game server address together with a score
// reflecting how confident DefaultSelector is that it is the active server.
// Higher scores indicate stronger candidates.
type ServerCandidate struct {
	IP    string
	Score int
}

// ServerSelector picks the most likely game server address out of a
// process's connections. It is an interface so the selection strategy can
// be improved or replaced later (e.g. once passive packet capture is
// available) without changing how callers use it.
type ServerSelector interface {
	Select(conns []Connection) (ServerCandidate, bool)
}

// DefaultSelector implements a simple frequency and connection-type based
// heuristic:
//
//   - Only public IPv4 remote addresses are considered.
//   - TCP connections are only considered when ESTABLISHED.
//   - UDP connections with a known (non-wildcard) remote address count too.
//   - Each occurrence adds to that address's score; established TCP
//     connections are weighted higher than UDP.
//   - The address with the highest score wins; ties are broken by the
//     address seen first.
type DefaultSelector struct{}

// Select implements ServerSelector.
func (DefaultSelector) Select(conns []Connection) (ServerCandidate, bool) {
	scores := make(map[string]int)
	var order []string

	for _, c := range conns {
		if c.Proto == "TCP" && !strings.EqualFold(c.State, "ESTABLISHED") {
			continue
		}

		host, _, err := ExtractHost(c.RemoteAddr)
		if err != nil {
			// Wildcard ("*:*") or otherwise unparsable address.
			continue
		}

		if !IsPublicIPv4(host) {
			continue
		}

		weight := scoreUDP
		if c.Proto == "TCP" {
			weight = scoreTCPEstablished
		}

		if _, seen := scores[host]; !seen {
			order = append(order, host)
		}
		scores[host] += weight
	}

	var best ServerCandidate
	found := false
	for _, host := range order {
		score := scores[host]
		if !found || score > best.Score {
			best = ServerCandidate{IP: host, Score: score}
			found = true
		}
	}

	return best, found
}
