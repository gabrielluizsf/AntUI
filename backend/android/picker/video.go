//go:build android

package picker

// Video is Image for a film.
func Video() (string, error) { return one("video/*") }
