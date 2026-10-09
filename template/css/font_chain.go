package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// fontChain is one answer to a Font query, with the faces of the names the
// list holds chained behind each other, or nil when it holds none of them.
func (sh *Sheet) fontChain(family string, weight uint16, slanted bool) *FontFace {
	var chain []*FontFace
	for _, name := range splitTopLevel(family) {
		if ff := sh.fontNamed(unquote(name), weight, slanted); ff != nil {
			chain = append(chain, ff)
		}
	}
	if len(chain) == 0 {
		return nil
	}
	// Every face in the chain is drawn at the size the program draws the
	// rest of its text at, so naming a family does not change how big the
	// text comes out — the stylesheet's own font-size is what sizes it.
	base := canvas.DefaultFace().Size()
	if base <= 0 {
		base = float64(DefaultFontSize)
	}
	// The chain is built from the end back: each face takes the one after
	// it as its fallback, so the list runs out at the face the program
	// draws with by default.
	face := canvas.DefaultFace()
	for i := len(chain) - 1; i >= 0; i-- {
		face = chain[i].face.AtSize(base).WithFallback(face)
	}
	head := *chain[0]
	head.face = face
	return &head
}
