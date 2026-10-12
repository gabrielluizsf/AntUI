//go:build android

package speech

// Options is how to listen.
type Options struct {
	// Language is a tag — "pt-BR", "en-GB". Empty uses the device's own.
	Language string
	// Partial reports guesses as they are made, not only the final answer.
	// They arrive with [Heard.Partial] set and are worth showing: a person
	// watching words appear knows they are being heard.
	Partial bool
	// Continuous starts listening again as soon as an answer is given, so
	// that a conversation can carry on. Without it, listening stops after
	// one utterance.
	//
	// It is not free — the microphone stays open — and on most devices the
	// engine is a network service, so it is also not private. Say so.
	Continuous bool
}
