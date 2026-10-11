package canvas

import "fmt"

// SetPalette gives a paletted canvas its colours and builds the table that
// maps a colour back to the nearest index.
func (cv *Canvas) SetPalette(entries []Color) error {
	if cv.format != Pal8 {
		return fmt.Errorf("antui: SetPalette on a %s canvas, not a paletted one", cv.format)
	}
	if len(entries) == 0 {
		return fmt.Errorf("antui: SetPalette was given no colours")
	}
	count := min(len(entries), 256)
	if cv.palette == nil {
		cv.palette = make([]Color, 256)
	}
	for i := range 256 {
		cv.palette[i] = entries[min(i, count-1)] | 0xFF000000
	}

	if cv.nearest == nil {
		cv.nearest = make([]byte, 32768)
	}
	for r := range 32 {
		for g := range 32 {
			for b := range 32 {
				red := r<<3 | r>>2
				green := g<<3 | g>>2
				blue := b<<3 | b>>2
				best, bestDist := 0, -1
				for i, p := range cv.palette {
					dr := red - int(p.R())
					dg := green - int(p.G())
					db := blue - int(p.B())
					d := dr*dr + dg*dg + db*db
					if bestDist < 0 || d < bestDist {
						bestDist, best = d, i
					}
				}
				cv.nearest[r<<10|g<<5|b] = byte(best)
			}
		}
	}
	return nil
}
