package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// backgroundLayer is the per-layer state the token scan fills before the
// pieces are folded into the style.
type backgroundLayer struct {
	img      BackImage
	posToks  []string
	repToks  []string
	att      uint8
	attSet   bool
	boxes    []uint8
	color    canvas.Color
	colorSet bool
}

// scanBackgroundLayer walks one layer's tokens, splitting out the image, the
// repeat keywords, the attachment, the box keywords and — on the last layer
// only — a trailing colour.
func scanBackgroundLayer(head string, last bool, ctx Units) backgroundLayer {
	var lay backgroundLayer
	for _, tk := range splitTokens(head) {
		switch tk {
		case "border-box":
			lay.boxes = append(lay.boxes, BackBorder)
		case "padding-box":
			lay.boxes = append(lay.boxes, BackPadding)
		case "content-box":
			lay.boxes = append(lay.boxes, BackContent)
		case "scroll":
			lay.att, lay.attSet = BackAttachScroll, true
		case "fixed":
			lay.att, lay.attSet = BackAttachFixed, true
		case "local":
			lay.att, lay.attSet = BackAttachLocal, true
		default:
			if bi, ok := parseBackImage(tk, ctx); ok {
				lay.img = bi
				continue
			}
			if last {
				if c, err := ParseColor(tk); err == nil {
					lay.color, lay.colorSet = c, true
					continue
				}
			}
			if isRepeatKeyword(tk) {
				lay.repToks = append(lay.repToks, tk)
				continue
			}
			lay.posToks = append(lay.posToks, tk)
		}
	}
	return lay
}
