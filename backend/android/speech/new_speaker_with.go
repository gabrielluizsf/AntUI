//go:build android

package speech

// NewSpeakerWith starts a particular engine, named by its package.
//
// Prefer [NewSpeaker]: which engine speaks is the user's choice, made in the
// settings, and an app that overrides it overrides an accessibility setting.
// This is for the case where there is no choice to respect — no default set
// — and for a test that has to be sure which engine answered.
func NewSpeakerWith(engine string) (*Speaker, error) { return newSpeaker(engine) }
