package audio

import (
	"embed"
	"fmt"
	"io/fs"
	"path"

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

// Tech is the futuristic sound bank: crisp electronic clicks, from the
// "scifi" pack of the UI SFX library (CC0). It is the voice of the Cyberpunk
// look.
var Tech = loadBank("Tech", "tech", map[string]string{
	"button":   "press.wav",
	"checkbox": "toggle-off.wav",
	"radio":    "select.wav",
	"text":     "typing.wav",
})

// Simple is the cartoonish sound bank: bouncy, playful clicks, from the
// "arcade" pack of the UI SFX library (CC0). It is the voice of the Builtin
// look.
var Simple = loadBank("Simple", "simple", map[string]string{
	"button":   "press.wav",
	"checkbox": "toggle-off.wav",
	"radio":    "select.wav",
	"text":     "typing.wav",
})

// loadBank decodes one pack's files into a ready bank. The role names line
// up with the components: button, checkbox, radio, slider, text.
func loadBank(name, pack string, files map[string]string) *Template {
	t := &Template{Name: name}
	load := func(role string, dst **Sound) {
		fn, ok := files[role]
		if !ok {
			panic(fmt.Sprintf("audio: %s has no %s sound", name, role))
		}
		data, err := fs.ReadFile(sfxFS, path.Join("sfx", pack, fn))
		if err != nil {
			panic(fmt.Sprintf("audio: %s %s: %v", name, role, err))
		}
		samples, rate, err := decodeWAV(data)
		if err != nil {
			panic(fmt.Sprintf("audio: %s/%s: %v", pack, fn, err))
		}
		// The decoder already swings to ninety-two percent; a small trim
		// keeps several sounds together from clipping the stream.
		*dst = &Sound{sample: samples, rate: rate, gain: 0.85}
	}
	load("button", &t.Button)
	load("checkbox", &t.Checkbox)
	load("radio", &t.Radio)
	load("text", &t.TextInput)
	return t
}
