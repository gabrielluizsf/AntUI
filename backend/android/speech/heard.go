//go:build android

package speech

// Heard is what the recogniser made of what was said.
type Heard struct {
	// Text is its best guess.
	Text string
	// Alternatives are all of them, best first, with Text at the front.
	Alternatives []string
	// Confidence is 0 to 1 for Text, and 0 when the engine will not say —
	// which many will not, so it cannot be used as a threshold.
	Confidence float32
	// Partial says this is a guess so far and not the answer. More will
	// follow, and the last one will not be partial.
	Partial bool
}
