package filter

import (
	"regexp"
	"strings"

	"github.com/logsnap/internal/snapshot"
)

// Options holds filtering criteria for log entries.
type Options struct {
	Level      string
	ServiceID  string
	MessageRe  string
	FieldKey   string
	FieldValue string
}

// Apply returns a new snapshot containing only entries that match all
// non-empty criteria in opts.
func Apply(snap *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	var re *regexp.Regexp
	if opts.MessageRe != "" {
		var err error
		re, err = regexp.Compile(opts.MessageRe)
		if err != nil {
			return nil, err
		}
	}

	out := snapshot.New(snap.ID+"-filtered", snap.Service)
	for _, e := range snap.Entries {
		if !matchEntry(e, opts, re) {
			continue
		}
		out.AddEntry(e)
	}
	return out, nil
}

func matchEntry(e snapshot.Entry, opts Options, re *regexp.Regexp) bool {
	if opts.Level != "" && !strings.EqualFold(e.Level, opts.Level) {
		return false
	}
	if opts.ServiceID != "" && e.ServiceID != opts.ServiceID {
		return false
	}
	if re != nil && !re.MatchString(e.Message) {
		return false
	}
	if opts.FieldKey != "" {
		v, ok := e.Fields[opts.FieldKey]
		if !ok {
			return false
		}
		if opts.FieldValue != "" && v != opts.FieldValue {
			return false
		}
	}
	return true
}
