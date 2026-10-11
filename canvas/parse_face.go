package canvas

import (
	"fmt"

	"github.com/gabrielluizsf/antui/font"
)

// ParseFace is [LoadFace] for a font already in memory.
func ParseFace(data []byte, pixels float64) (*Face, error) {
	if pixels <= 0 {
		return nil, fmt.Errorf("antui: %v is not a size", pixels)
	}
	f, err := font.ParseTTF(data)
	if err != nil {
		return nil, err
	}
	scale := f.Scaled(pixels)
	return &Face{
		file:     f,
		size:     pixels,
		scale:    scale,
		ascent:   int(float64(f.Ascent)*scale + 0.5),
		descent:  int(-float64(f.Descent)*scale + 0.5),
		height:   f.LineHeight(pixels),
		hasColor: f.HasColorGlyphs(),
		glyphs:   map[rune]*font.Mask{},
	}, nil
}
