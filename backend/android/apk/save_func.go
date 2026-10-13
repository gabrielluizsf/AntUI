package apk

import "os"

// save copies a file into the cache, atomically. Every failure is ignored:
// the caller has what it needs and the cache is only ever an optimisation.
func save(from, to string) {
	body, err := os.ReadFile(from)
	if err != nil {
		return
	}
	tmp := to + ".tmp"
	if os.WriteFile(tmp, body, 0o644) != nil {
		return
	}
	if os.Rename(tmp, to) != nil {
		os.Remove(tmp)
	}
}
