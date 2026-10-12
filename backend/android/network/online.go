//go:build android

package network

// Online is Read().Online, for the one question most callers have.
func Online() bool {
	s, err := Read()
	return err == nil && s.Online
}
