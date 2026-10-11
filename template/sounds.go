package template

import "github.com/gabrielluizsf/antui/template/event"

// Sounds is the bank the components' events ring through right now.
func (c *CSS) Sounds() event.SoundBank { return c.ui.sound }
