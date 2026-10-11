package template

import "github.com/gabrielluizsf/antui/template/event"

// New joins a painting and a hearing into a working template: everything the
// [Style] does not say — how a press feels, how focus moves, what event an
// interaction counts as — is the same for every template, and the package
// does it. The sound bank hears every event the components report.
//
//	myStyle := MyStyle{}
//	mySound := MySoundBank{}
//	tpl := template.New(myStyle, mySound)
func New(style Style, sound event.SoundBank) Template {
	return &uiTemplate{style: style, sound: sound, ld: newLoading()}
}
