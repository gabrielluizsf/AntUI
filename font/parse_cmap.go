package font

import (
	"fmt"
)

// parseCmap picks the best subtable in a cmap and prepares it.
//
// "Best" is a Unicode one, and among those the widest: a format 12 table can
// say things a format 4 table cannot, and a font that has both means the
// same by them for the characters they share.
func parseCmap(data []byte) (CmapTable, error) {
	if len(data) < 4 {
		return CmapTable{}, fmt.Errorf("antUI: a truncated cmap")
	}
	best, err := bestCmapOffset(data)
	if err != nil {
		return CmapTable{}, err
	}
	sub := data[best:]
	if len(sub) < 4 {
		return CmapTable{}, fmt.Errorf("antUI: a truncated cmap subtable")
	}
	return readCmapSubtable(sub)
}
