package antui

import (
	"errors"
)

// ErrNoAudio is what opening gives back where there is no sound at all: a
// machine with no card, a session with no server, a platform nothing is
// written for yet. It is worth telling apart from a broken device, because a
// game carries on perfectly well in silence and should not stop for it.
var ErrNoAudio = errors.New("antui: no audio device")
