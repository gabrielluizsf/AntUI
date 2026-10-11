package svg

import "github.com/gabrielluizsf/antui/canvas"

// paintedBox is the box a node paints over, in the drawing's own coordinates:
// every shape it draws and every shape inside it draws, each where it was put
// — the shapes carry their transforms already, so this is the box as it lands
// on the page rather than as it was written. A picture has no transform put
// into it — the transform meets it while it is painted — so its box is taken
// through the element's own here for the same result. Writing has no box of
// its own until a font is asked, which is what [walkRuns] and [writingBox] do
// where this walks, and the letters land through the element's transform the
// same way a picture's box does. A node with nothing in it that paints has no
// box, and says so rather than answering an empty one.
func paintedBox(n *Node) ([4]float64, bool) {
	var minMax [4]float64
	found := false
	add := func(x0, y0, x1, y1 float64) {
		if !found {
			minMax = [4]float64{x0, y0, x1, y1}
			found = true
			return
		}
		minMax[0] = min(minMax[0], x0)
		minMax[1] = min(minMax[1], y0)
		minMax[2] = max(minMax[2], x1)
		minMax[3] = max(minMax[3], y1)
	}
	var walk func(n *Node)
	walk = func(n *Node) {
		if n.Path != nil && !n.Path.Empty() {
			if x0, y0, x1, y1, ok := n.Path.Bounds(); ok {
				add(x0, y0, x1, y1)
			}
		}
		if n.Pic != nil {
			x, y, w, h := n.Pic.Box[0], n.Pic.Box[1], n.Pic.Box[2], n.Pic.Box[3]
			if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
				x, y, w, h = movedBox(x, y, w, h, n.Style.Transform)
			}
			add(x, y, x+w, y+h)
		}
		if len(n.Runs) > 0 {
			pieces := walkRuns(n.Runs, canvas.Identity())
			// Every run hidden and every run empty leaves no pieces, and a box
			// of nothing is not one to measure a region against: the same no
			// box at all as a node with no writing in it.
			if len(pieces) > 0 {
				wb := writingBox(pieces)
				x0, y0, x1, y1 := wb[0], wb[1], wb[2], wb[3]
				if x1 > x0 || y1 > y0 {
					if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
						x, y, w, h := movedBox(x0, y0, x1-x0, y1-y0, n.Style.Transform)
						x0, y0, x1, y1 = x, y, x+w, y+h
					}
					add(x0, y0, x1, y1)
				}
			}
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	walk(n)
	if !found {
		return [4]float64{}, false
	}
	return [4]float64{minMax[0], minMax[1], minMax[2] - minMax[0], minMax[3] - minMax[1]}, true
}
