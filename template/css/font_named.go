package css

import (
	"strings"
)

// fontNamed is the face the sheet holds for one family name at a weight and a
// slant, or nil when it holds none.
func (sh *Sheet) fontNamed(name string, weight uint16, slanted bool) *FontFace {
	want := normWeight(weight)
	for _, lean := range [2]bool{slanted, !slanted} {
		var pick *FontFace
		best := -1
		for _, ff := range sh.fonts {
			if !strings.EqualFold(ff.Family, name) || ff.Slanted() != lean {
				continue
			}
			if r := weightRank(want, ff.Weight); pick == nil || r < best {
				pick, best = ff, r
			}
		}
		if pick != nil {
			return pick
		}
	}
	return nil
}
