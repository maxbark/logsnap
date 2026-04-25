// Package validate provides snapshot integrity checking for logsnap.
package validate

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

// Result holds the outcome of a validation run.
type Result struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

var validLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
	"fatal": true,
}

var serviceIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// Check validates a snapshot for structural and semantic correctness.
func Check(snap *snapshot.Snapshot) (Result, error) {
	if snap == nil {
		return Result{}, errors.New("snapshot is nil")
	}

	result := Result{Valid: true}

	if snap.Meta.Label == "" {
		result.Warnings = append(result.Warnings, "snapshot has no label")
	}

	for i, entry := range snap.Entries {
		prefix := fmt.Sprintf("entry[%d]", i)

		if entry.Timestamp.IsZero() {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: missing timestamp", prefix))
			result.Valid = false
		} else if entry.Timestamp.After(time.Now().Add(5 * time.Minute)) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s: timestamp is in the future", prefix))
		}

		if entry.Message == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: missing message", prefix))
			result.Valid = false
		}

		if entry.Level != "" && !validLevels[entry.Level] {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s: unrecognized level %q", prefix, entry.Level))
		}

		if entry.ServiceID != "" && !serviceIDPattern.MatchString(entry.ServiceID) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s: service_id %q contains unusual characters", prefix, entry.ServiceID))
		}
	}

	return result, nil
}
