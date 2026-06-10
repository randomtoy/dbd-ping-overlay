// Package process locates running processes by image name without reading
// their memory or attaching to them in any way. It currently shells out to
// the built-in "tasklist" command.
package process

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Lister returns the process IDs of all running processes that match a
// given executable name. Implementations must not read process memory or
// otherwise interact with the target process beyond listing it.
type Lister interface {
	ListPIDs(ctx context.Context, processName string) ([]int, error)
}

// TasklistLister implements Lister using the Windows "tasklist" command.
type TasklistLister struct{}

// ListPIDs runs "tasklist" filtered by image name and returns the PIDs of
// every matching process. If no process matches, it returns an empty slice
// and a nil error.
func (TasklistLister) ListPIDs(ctx context.Context, processName string) ([]int, error) {
	cmd := exec.CommandContext(ctx, "tasklist",
		"/FI", fmt.Sprintf("IMAGENAME eq %s", processName),
		"/FO", "CSV",
		"/NH",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run tasklist: %w", err)
	}

	return parseTasklistCSV(string(out))
}

// parseTasklistCSV parses the headerless CSV output of
//
//	tasklist /FI "IMAGENAME eq <name>" /FO CSV /NH
//
// and returns the PID column of every row. When the process is not running,
// tasklist prints an informational line instead of CSV data; in that case
// parseTasklistCSV returns an empty, nil-error result.
func parseTasklistCSV(output string) ([]int, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return nil, nil
	}

	// When no process matches the filter, tasklist prints a localized
	// informational message (e.g. "INFO: No tasks are running which match
	// the specified criteria.") instead of a CSV row.
	if !strings.HasPrefix(trimmed, `"`) {
		return nil, nil
	}

	reader := csv.NewReader(strings.NewReader(trimmed))

	var pids []int
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("parse tasklist output: %w", err)
		}

		// Expected columns: Image Name, PID, Session Name, Session#, Mem Usage
		if len(record) < 2 {
			continue
		}

		pid, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			continue
		}

		pids = append(pids, pid)
	}

	return pids, nil
}
