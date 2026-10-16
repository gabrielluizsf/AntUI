package svg

import "github.com/gabrielluizsf/antui/canvas"

// walk takes in the box of every shape n draws and every shape inside it draws.
func (b *boxBounds) walk(n *Node) {
	if n.Path != nil && !n.Path.Empty() {
		if x0, y0, x1, y1, ok := n.Path.Bounds(); ok {
			b.add(x0, y0, x1, y1)
		}
	}
	if n.Pic != nil {
		box := n.Pic.fittedBox()
		x, y, w, h := box[0], box[1], box[2], box[3]
		if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
			x, y, w, h = movedBox(x, y, w, h, n.Style.Transform)
		}
		b.add(x, y, x+w, y+h)
	}
	if len(n.Runs) > 0 {
		b.addWriting(n)
	}
	for _, k := range n.Kids {
		b.walk(k)
	}
}

// addWriting takes in the box of the writing, which needs a font to measure.
// Every run hidden and every run empty leaves no pieces, and a box of nothing
// is not one to measure a region against: the same no box at all as a node with
// no writing in it.
func (b *boxBounds) addWriting(n *Node) {
	pieces := walkRuns(n.Runs, canvas.Identity())
	if len(pieces) == 0 {
		return
	}
	wb := writingBox(pieces)
	x0, y0, x1, y1 := wb[0], wb[1], wb[2], wb[3]
	if x1 <= x0 && y1 <= y0 {
		return
	}
	if n.Style.HasTransform && n.Style.Transform != (canvas.Matrix{}) {
		x, y, w, h := movedBox(x0, y0, x1-x0, y1-y0, n.Style.Transform)
		x0, y0, x1, y1 = x, y, x+w, y+h
	}
	b.add(x0, y0, x1, y1)
}
