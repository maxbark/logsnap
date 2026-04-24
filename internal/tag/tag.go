// Package tag provides functionality for tagging snapshots with
// arbitrary key-value metadata labels.
package tag

import (
	"fmt"
	"regexp"

	"github.com/user/logsnap/internal/snapshot"
)

var validKey = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Apply attaches the provided tags to the snapshot's metadata.
// Tags are expressed as "key=value" strings. Existing tags with
// the same key are overwritten. Returns an error if any tag is
// malformed or if the snapshot is nil.
func Apply(snap *snapshot.Snapshot, tags []string) error {
	if snap == nil {
		return fmt.Errorf("tag: snapshot must not be nil")
	}

	parsed, err := parseTags(tags)
	if err != nil {
		return err
	}

	if snap.Meta.Tags == nil {
		snap.Meta.Tags = make(map[string]string)
	}

	for k, v := range parsed {
		snap.Meta.Tags[k] = v
	}

	return nil
}

// Remove deletes the named tag keys from the snapshot's metadata.
// Keys that do not exist are silently ignored.
func Remove(snap *snapshot.Snapshot, keys []string) error {
	if snap == nil {
		return fmt.Errorf("tag: snapshot must not be nil")
	}

	if snap.Meta.Tags == nil {
		return nil
	}

	for _, k := range keys {
		delete(snap.Meta.Tags, k)
	}

	return nil
}

// parseTags converts a slice of "key=value" strings into a map.
func parseTags(raw []string) (map[string]string, error) {
	out := make(map[string]string, len(raw))
	for _, t := range raw {
		k, v, ok := splitTag(t)
		if !ok {
			return nil, fmt.Errorf("tag: invalid format %q, expected key=value", t)
		}
		if !validKey.MatchString(k) {
			return nil, fmt.Errorf("tag: invalid key %q, must match [a-zA-Z][a-zA-Z0-9_-]*", k)
		}
		out[k] = v
	}
	return out, nil
}

// splitTag splits "key=value" into its components.
func splitTag(s string) (key, value string, ok bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
