package font

import (
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

func TestParseTTF(t *testing.T) {
	// A font must have at least 12 bytes.
	_, err := ParseTTF([]byte{1, 2, 3})
	assert.Error(t, err)

	// A CFF font (tag "OTTO") is unsupported and should return an error.
	ottoFont := []byte{0x4F, 0x54, 0x54, 0x4F, 0, 0, 0, 0, 0, 0, 0, 0}
	_, err = ParseTTF(ottoFont)
	assert.Error(t, err)

	// An invalid tag should return an error.
	invalidFont := []byte{0x11, 0x22, 0x33, 0x44, 0, 0, 0, 0, 0, 0, 0, 0}
	_, err = ParseTTF(invalidFont)
	assert.Error(t, err)

	// A truncated font collection should trigger an error.
	truncatedCollection := []byte{0x74, 0x74, 0x63, 0x66, 0, 0, 0, 0, 0, 0, 0, 0}
	_, err = ParseTTF(truncatedCollection)
	assert.Error(t, err)
}

func TestTTF_Advance(t *testing.T) {
	// When there are no horizontal metrics, it defaults to half the UnitsPerEm.
	fDefault := &TTF{NumHMetrics: 0, UnitsPerEm: 1000}
	assert.Equal(t, int16(500), fDefault.Advance(0))

	// When horizontal metrics exist, it parses the Hmtx table (4 bytes per metric).
	// 0x0190 = 400 in decimal.
	fMetrics := &TTF{
		NumHMetrics: 1,
		UnitsPerEm:  1000,
		Hmtx:        []byte{0x01, 0x90, 0x00, 0x00},
	}
	assert.Equal(t, int16(400), fMetrics.Advance(0))
}

func TestTTF_Outline(t *testing.T) {
	f := &TTF{
		NumGlyphs: 2,
		LongLoca:  false,
		// Offsets: 0x0000 -> 0, 0x0002 -> 4, 0x0004 -> 8 (multiplied by 2)
		Loca: []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x04},
		Glyf: []byte{10, 20, 30, 40, 50, 60, 70, 80},
	}

	// Test first glyph.
	outline0 := f.Outline(0)
	assert.Len(t, 4, outline0)
	assert.DeepEqual(t, []byte{10, 20, 30, 40}, outline0)

	// Test second glyph.
	outline1 := f.Outline(1)
	assert.Len(t, 4, outline1)
	assert.DeepEqual(t, []byte{50, 60, 70, 80}, outline1)

	// Test out of bounds glyph index (should return nothing).
	outlineInvalid := f.Outline(99)
	assert.Equal(t, 0, len(outlineInvalid))
}

func TestTTF_GlyphContours(t *testing.T) {
	f := &TTF{
		NumGlyphs: 1,
		LongLoca:  false,
		Loca:      []byte{0x00, 0x00, 0x00, 0x01},
		Glyf:      []byte{1, 2}, // Truncated glyph data (less than 10 bytes).
	}

	// Maximum depth exceeded.
	resDepth := f.GlyphContours(0, 6)
	assert.Equal(t, 0, len(resDepth))

	// Truncated data fallback.
	resTruncated := f.GlyphContours(0, 0)
	assert.Equal(t, 0, len(resTruncated))
}

func TestTTF_Scaled(t *testing.T) {
	f := &TTF{UnitsPerEm: 1000}
	
	// 16.0 pixels / 1000 UnitsPerEm = 0.016
	assert.Equal(t, 0.016, f.Scaled(16.0))
}

func TestTTF_LineHeight(t *testing.T) {
	f := &TTF{
		UnitsPerEm: 1000,
		Ascent:     800,
		Descent:    -200,
		LineGap:    90,
	}

	// Total font units: 800 - (-200) + 90 = 1090
	// Scaled by 16px: 1090 * (16 / 1000) = 17.44
	// Ceiled value: 18
	assert.Equal(t, 18, f.LineHeight(16.0))
}