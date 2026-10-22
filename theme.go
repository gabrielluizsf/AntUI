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
	Shadow        canvas.Color // the tint a raised surface casts under itself
	Radius        int          // radius of a control's rounded corners
	RadiusLarge   int          // radius of a container's — a card, a panel
	Padding       int          // inner spacing
}

// LightTheme is the default: dark text on a light background, and what a
// window starts from unless the system it opens on paints its own interface
// dark.
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
		Shadow:        0x1E14141B,
		Radius:        10,
		RadiusLarge:   16,
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
		// Shadows read poorly once the surface they would fall on is already
		// dark; a quieter one than the light theme's still gives a popover or
		// a dialog some separation from the window behind it, and the border
		// above does the rest of the lifting that a shadow would elsewhere.
		Shadow:      0x3C000000,
		Radius:      10,
		RadiusLarge: 16,
		Padding:     8,
	}
}
