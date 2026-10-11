package font

import (
	"fmt"
)

// bestCmapOffset is where the best Unicode subtable in a cmap begins.
func bestCmapOffset(data []byte) (int, error) {
	count := int(be16(data, 2))
	best, bestScore := -1, -1
	for i := range count {
		rec := 4 + i*8
		if rec+8 > len(data) {
			break
		}
		platform, encoding := be16(data, rec), be16(data, rec+2)
		offset := int(be32(data, rec+4))
		if offset >= len(data) {
			continue
		}
		score := -1
		switch {
		case platform == 3 && encoding == 10: // Windows, full Unicode
			score = 4
		case platform == 0: // Unicode, any version
			score = 3
		case platform == 3 && encoding == 1: // Windows, BMP
			score = 2
		case platform == 1 && encoding == 0: // Macintosh Roman
			score = 1
		}
		if score > bestScore {
			best, bestScore = offset, score
		}
	}
	if best < 0 {
		return 0, fmt.Errorf("antUI: the font has no Unicode cmap")
	}
	return best, nil
}
