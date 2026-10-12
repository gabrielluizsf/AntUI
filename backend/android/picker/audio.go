//go:build android

package picker

// Audio is Image for a sound.
func Audio() (string, error) { return one("audio/*") }
