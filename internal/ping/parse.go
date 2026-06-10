package ping

import (
	"regexp"
	"strconv"

	"github.com/randomtoy/dbd-ping-overlay/internal/model"
)

// unavailableMessage is shown when no average latency could be parsed,
// typically because every request timed out and ICMP is being blocked or
// filtered somewhere along the path.
const unavailableMessage = "Ping unavailable / ICMP may be blocked"

// Average latency, e.g.:
//
//	Average = 38ms                (English)
//	Среднее = 38мсек               (Russian, no space)
//	Среднее = 38 мсек              (Russian, with space)
var (
	averageEnglishRe = regexp.MustCompile(`Average\s*=\s*(\d+)\s*ms`)
	averageRussianRe = regexp.MustCompile(`Среднее\s*=\s*(\d+)\s*мсек`)
)

// Packet loss, e.g.:
//
//	(0% loss)    (English)
//	(0% потерь)  (Russian)
var (
	lossEnglishRe = regexp.MustCompile(`\((\d+)%\s*loss\)`)
	lossRussianRe = regexp.MustCompile(`\((\d+)%\s*потерь\)`)
)

// parsePingOutput extracts average latency and packet loss from the
// combined output of the "ping" command. It supports both English and
// Russian Windows locales. If no average latency can be found, the returned
// status has Available set to false and Message explains why.
func parsePingOutput(output string) model.PingStatus {
	var status model.PingStatus

	if avg, ok := matchInt(averageEnglishRe, output); ok {
		status.Available = true
		status.AverageMs = avg
	} else if avg, ok := matchInt(averageRussianRe, output); ok {
		status.Available = true
		status.AverageMs = avg
	}

	if loss, ok := matchInt(lossEnglishRe, output); ok {
		status.PacketLossKnown = true
		status.PacketLossPercent = loss
	} else if loss, ok := matchInt(lossRussianRe, output); ok {
		status.PacketLossKnown = true
		status.PacketLossPercent = loss
	}

	if !status.Available {
		status.Message = unavailableMessage
	}

	return status
}

// matchInt applies re to s and returns the first capture group parsed as an
// int. ok is false if re does not match or the capture group is not a
// valid integer.
func matchInt(re *regexp.Regexp, s string) (value int, ok bool) {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return 0, false
	}

	n, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, false
	}

	return n, true
}
