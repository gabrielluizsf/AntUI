package svg

import (
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// keepTile keeps a picture drawn under this key and chain for whatever asks
// the same of the pattern later — and where a picture drawn in parallel got
// there first, lets go of this one so that only one of the two is held. The
// chain is copied: the slice it comes in belongs to the caller, which goes on
// appending to it as the drawing goes deeper.
func (pt *pattern) keepTile(key tileKey, chain []*pattern, pic *canvas.Canvas) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	for _, t := range pt.tiles {
		if t.key == key && slices.Equal(t.chain, chain) {
			return
		}
	}
	pt.tiles = append(pt.tiles, tilePicture{key: key, chain: slices.Clone(chain), pic: pic})
}
