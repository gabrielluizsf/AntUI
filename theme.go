package antui

import "github.com/gabrielluizsf/antui/canvas"

// Theme is the set of colours and metrics the ready-made widgets draw with.
// Every field is public and may be changed at will — Window.Theme returns a
// pointer to the live one.
type Theme struct {
	Background    canvas.Color
	Surface       canvas.Color
	SurfaceHover  canvas.Color
	SurfaceActive canvas.Color
	Border        canvas.Color
	Accent        canvas.Color
	AccentHover   canvas.Color
	Text          canvas.Color
	TextMuted     canvas.Color
	TextOnAccent  canvas.Color
	Radius        int // radius of the rounded corners
	Padding       int // inner spacing
}

// LightTheme is the default: dark text on a light background.
func LightTheme() Theme {
	return Theme{
		Background:    0xFFF4F4F6,
		Surface:       0xFFFFFFFF,
		SurfaceHover:  0xFFECECF0,
		SurfaceActive: 0xFFDDDDE3,
		Border:        0xFFC9C9D0,
		Accent:        0xFF3E63DD,
		AccentHover:   0xFF5472E4,
		Text:          0xFF1C1C1F,
		TextMuted:     0xFF6F6F7B,
		TextOnAccent:  0xFFFFFFFF,
		Radius:        6,
		Padding:       8,
	}
}

// DarkTheme is the same set of roles in reverse: light text on a dark
// background, with the accent lifted so it still reads against it.
func DarkTheme() Theme {
	return Theme{
		Background:    0xFF1A1A1E,
		Surface:       0xFF26262B,
		SurfaceHover:  0xFF303038,
		SurfaceActive: 0xFF3A3A44,
		Border:        0xFF44444E,
		Accent:        0xFF5472E4,
		AccentHover:   0xFF6B85E8,
		Text:          0xFFF0F0F3,
		TextMuted:     0xFF9B9BA6,
		TextOnAccent:  0xFFFFFFFF,
		Radius:        6,
		Padding:       8,
	}
}
