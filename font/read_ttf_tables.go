package font

import (
	"fmt"
)

// readTTFTables fills the tables a face needs out of a font's directory.
func readTTFTables(f *TTF, tables map[string][]byte) error {
	head, ok := tables["head"]
	if !ok || len(head) < 54 {
		return fmt.Errorf("antUI: the font has no head table")
	}
	f.UnitsPerEm = be16(head, 18)
	if f.UnitsPerEm == 0 {
		return fmt.Errorf("antUI: the font says it has no units per em")
	}
	f.LongLoca = be16(head, 50) == 1

	maxp, ok := tables["maxp"]
	if !ok || len(maxp) < 6 {
		return fmt.Errorf("antUI: the font has no maxp table")
	}
	f.NumGlyphs = int(be16(maxp, 4))

	hhea, ok := tables["hhea"]
	if !ok || len(hhea) < 36 {
		return fmt.Errorf("antUI: the font has no hhea table")
	}
	f.Ascent = int16(be16(hhea, 4))
	f.Descent = int16(be16(hhea, 6))
	f.LineGap = int16(be16(hhea, 8))
	f.NumHMetrics = int(be16(hhea, 34))

	f.Loca, f.Glyf, f.Hmtx = tables["loca"], tables["glyf"], tables["hmtx"]
	if f.Glyf == nil || f.Loca == nil {
		return fmt.Errorf("antUI: the font has no TrueType outlines; a " +
			"CFF font cannot be drawn here")
	}
	cm, ok := tables["cmap"]
	if !ok {
		return fmt.Errorf("antUI: the font has no cmap table")
	}
	c, err := parseCmap(cm)
	if err != nil {
		return err
	}
	f.Cmap = c

	// Colour bitmaps are read when both tables are there and they make
	// sense; a file with one without the other, or with a table that does
	// not read, is a font without colour glyphs rather than a font at all.
	f.cbdt = tables["CBDT"]
	f.colorStrikes = parseColorStrikes(tables["CBLC"], f.cbdt)
	return nil
}
