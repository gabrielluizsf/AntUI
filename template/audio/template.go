package audio

import (
	"embed"

	"github.com/gabrielluizsf/antui/template/event"
)

//go:embed sfx
var sfxFS embed.FS

// A Template is one set of sounds, one for each component a template draws.
// It is also an [event.SoundBank]: hand it to a visual template and it plays
// the sound that answers each component event, so a game never touches audio
// itself. Two come ready-made — [Tech] for the futuristic look, [Simple] for
// the cartoonish one — and building a third is nothing more than loading
// sounds into the fields.
type Template struct {
	Name      string // what the bank calls itself, for listings and save files
	Button    *Sound // a button pressed
	Checkbox  *Sound // a checkbox flipped
	Radio     *Sound // a radio picked
	TextInput *Sound // text typed into a field
}

// On plays the sound that answers a component event. A nil bank — or a bank
// with a nil sound — answers in silence.
func (t *Template) On(e event.Event) {
	if t == nil {
		return
	}
	switch {
	case e.Is(event.Button, event.Click):
		t.Button.Play()
	case e.Is(event.Checkbox, event.Toggle):
		t.Checkbox.Play()
	case e.Is(event.Switch, event.Toggle):
		t.Checkbox.Play()
	case e.Is(event.Radio, event.Pick):
		t.Radio.Play()
	case e.Is(event.Select, event.Pick):
		t.Radio.Play()
	case e.Is(event.DatePicker, event.Change):
		t.Radio.Play()
	case e.Is(event.DatePicker, event.Open):
		t.Button.Play()
	case e.Is(event.TextInput, event.Focus):
		t.TextInput.Play()
	case e.Is(event.TextInput, event.Type):
		t.TextInput.PlayStep(e.Step)
	case e.Is(event.TextArea, event.Type):
		t.TextInput.PlayStep(e.Step)
	case e.Is(event.Select, event.Open):
		t.Button.Play()
	}
}
