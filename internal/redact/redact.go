// Package redact provides functionality for masking sensitive fields in log snapshots.
package redact

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourorg/logsnap/internal/snapshot"
)

const defaultMask = "[REDACTED]"

// Options controls redaction behavior.
type Options struct {
	// Fields is a list of field keys whose values will be masked.
	Fields []string
	// Patterns is a list of regex patterns; any value matching a pattern is masked.
	Patterns []string
	// Mask is the replacement string. Defaults to "[REDACTED]".
	Mask string
}

// Apply returns a new snapshot with sensitive data masked according to opts.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, fmt.Errorf("redact: snapshot is nil")
	}

	mask := opts.Mask
	if mask == "" {
		mask = defaultMask
	}

	compiledPatterns, err := compilePatterns(opts.Patterns)
	if err != nil {
		return nil, err
	}

	fieldSet := make(map[string]struct{}, len(opts.Fields))
	for _, f := range opts.Fields {
		fieldSet[strings.ToLower(f)] = struct{}{}
	}

	out := snapshot.New(snap.Label, snap.Source)
	for _, entry := range snap.Entries {
		redacted := redactEntry(entry, fieldSet, compiledPatterns, mask)
		out.AddEntry(redacted)
	}
	return out, nil
}

func redactEntry(e snapshot.Entry, fields map[string]struct{}, patterns []*regexp.Regexp, mask string) snapshot.Entry {
	copy := e
	copy.Fields = make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		if _, matched := fields[strings.ToLower(k)]; matched {
			copy.Fields[k] = mask
			continue
		}
		if matchesAny(v, patterns) {
			copy.Fields[k] = mask
			continue
		}
		copy.Fields[k] = v
	}
	if matchesAny(copy.Message, patterns) {
		copy.Message = mask
	}
	return copy
}

func compilePatterns(raw []string) ([]*regexp.Regexp, error) {
	var out []*regexp.Regexp
	for _, p := range raw {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		out = append(out, re)
	}
	return out, nil
}

func matchesAny(s string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}
