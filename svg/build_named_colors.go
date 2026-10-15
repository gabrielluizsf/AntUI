package svg

import "github.com/gabrielluizsf/antui/canvas"

// buildNamedColors turns the table of hex values above into colours, checking
// each of them as it goes so a typo in the table is a bug that stops the program
// here rather than a colour that is quietly wrong in a drawing.
func buildNamedColors() map[string]canvas.Color {
	out := make(map[string]canvas.Color, len(namedColorHex)+1)
	for name, hex := range namedColorHex {
		c, err := canvas.ParseColor(hex)
		if err != nil {
			// The table above is written by hand, so a typo in it is a bug
			// that must not reach a caller as a colour that is quietly wrong.
			panic("antui/svg: the named colour " + name + " is not a colour: " + err.Error())
		}
		out[name] = c
	}
	out["transparent"] = canvas.RGBA(canvas.RGBAParams{R: 0, G: 0, B: 0, A: 0})
	return out
}
