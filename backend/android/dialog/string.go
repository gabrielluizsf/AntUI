//go:build android

package dialog

// String names the button, which is what a log line wants.
func (b Button) String() string {
	switch b {
	case Positive:
		return "positive"
	case Negative:
		return "negative"
	case Neutral:
		return "neutral"
	}
	return "cancelled"
}
