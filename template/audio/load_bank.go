package audio

import (
	"fmt"
	"io/fs"
	"path"
)

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
