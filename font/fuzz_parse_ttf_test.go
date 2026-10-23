package font

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzParseTTF feeds the font reader files nobody vetted: whatever a font
// says about itself — table offsets, glyph counts, contour sizes, component
// lists, bitmap strikes — must end in an error or a font, never a panic.
func FuzzParseTTF(f *testing.F) {
	seeds, _ := filepath.Glob("../testdata/fonts/*.ttf")
	for _, path := range seeds {
		if data, err := os.ReadFile(path); err == nil {
			f.Add(data)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		ttf, err := ParseTTF(data)
		if err != nil {
			return
		}
		for _, r := range []rune{'A', 'g', 0x00E9, 0x1F600} {
			ttf.Cmap.Glyph(r)
		}
		for g := range min(ttf.NumGlyphs, 64) {
			ttf.Advance(g)
			Render(ttf.GlyphContours(g, 0), ttf.Scaled(16))
			ttf.ColorGlyph(g, 16)
		}
	})
}
