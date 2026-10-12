package svg

import "strings"

// placeAt reads one of the two names a place may be written under, either as
// one number for the whole of the writing or as one number for each letter of
// it.
func placeAt(runs []TextRun, e *element, name string, warn func(string, ...any)) {
	raw := strings.TrimSpace(e.attr(name))
	if raw == "" {
		return
	}
	args := splitArgs(raw)
	if len(args) == 0 {
		return
	}
	if len(args) == 1 {
		placeOne(runs, name, raw, args[0], warn)
		return
	}
	placeList(runs, name, raw, args, warn)
}
