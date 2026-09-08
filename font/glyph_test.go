package font

import (
	"testing"

	"github.com/gabrielluizsf/antui/assert"
)

func TestGlyph(t *testing.T) {
	t.Run("returns correct slice length", func(t *testing.T) {
		glyph := Glyph('A')
		assert.Len(t, Height, glyph)
	})

	t.Run("handles ASCII boundary characters", func(t *testing.T) {
		glyphLo := Glyph(ASCIILo)
		glyphHi := Glyph(ASCIIHi)

		assert.Len(t, Height, glyphLo)
		assert.Len(t, Height, glyphHi)
	})

	t.Run("handles Latin boundary characters", func(t *testing.T) {
		glyphLo := Glyph(LatinLo)
		glyphHi := Glyph(LatinHi)

		assert.Len(t, Height, glyphLo)
		assert.Len(t, Height, glyphHi)
	})

	t.Run("falls back to question mark for out of range runes", func(t *testing.T) {
		expectedFallback := Glyph('?')
		unsupportedGlyph := Glyph('⌘')

		assert.DeepEqual(t, expectedFallback, unsupportedGlyph)
	})

	t.Run("returns expected slice data for ASCII character", func(t *testing.T) {
		glyphA := Glyph('A')

		indexA := int('A') - ASCIILo
		expectedA := font8x16[indexA*Height : indexA*Height+Height]

		assert.DeepEqual(t, expectedA, glyphA)
	})

	t.Run("differentiates between distinct characters", func(t *testing.T) {
		glyphA := Glyph('A')
		glyphB := Glyph('B')

		isSame := true
		for i := range glyphA {
			if glyphA[i] != glyphB[i] {
				isSame = false
				break
			}
		}
		assert.False(t, isSame)
	})
}