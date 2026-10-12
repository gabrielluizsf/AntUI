package audio

// Simple is the cartoonish sound bank: bouncy, playful clicks, from the
// "arcade" pack of the UI SFX library (CC0). It is the voice of the Builtin
// look.
var Simple = loadBank("Simple", "simple", map[string]string{
	"button":   "press.wav",
	"checkbox": "toggle-off.wav",
	"radio":    "select.wav",
	"text":     "typing.wav",
})
