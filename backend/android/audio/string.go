//go:build android

package audio

// String names the change, which is what a log line wants.
func (c Change) String() string {
	switch c {
	case Gained:
		return "gained"
	case Lost:
		return "lost"
	case LostBriefly:
		return "lost briefly"
	case Duck:
		return "duck"
	}
	return "unknown"
}
