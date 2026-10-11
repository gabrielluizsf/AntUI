package template

import "github.com/gabrielluizsf/antui/template/event"

// SetSound swaps the sound bank the components ring their events through. The
// default is the [github.com/gabrielluizsf/antui/template/audio.Simple]
// bank; a nil bank makes the template silent.
func (c *CSS) SetSound(s event.SoundBank) { c.ui.sound = s }
