//go:build android

package speech

// NewSpeaker starts the user's own engine and waits for it to be ready.
//
// It waits because there is nothing useful to do with an engine that is not:
// asking it to speak before it has bound is refused, silently, and the words
// are simply never said.
//
// A device with an engine installed but none chosen as the default has no
// default, and this fails with [ErrNoEngine]. [Engines] lists what is there
// and [NewSpeakerWith] names one.
func NewSpeaker() (*Speaker, error) { return newSpeaker("") }
