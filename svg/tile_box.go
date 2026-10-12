package svg

// tileBox is where the tile is and how big, in the drawing's own coordinates:
// the numbers the pattern wrote, taken into the shape's box where the units say
// so. It answers false where the box it would be a fraction of has no area, or
// where the tile itself has none.
func tileBox(pt *pattern, n *Node) (x, y, w, h float64, ok bool) {
	box := n.box
	x, y, w, h = pt.x, pt.y, pt.w, pt.h
	if pt.boxUnits {
		if box[2] <= 0 || box[3] <= 0 {
			return 0, 0, 0, 0, false
		}
		x, y = box[0]+x*box[2], box[1]+y*box[3]
		w, h = w*box[2], h*box[3]
	}
	if pt.contentBox && (box[2] <= 0 || box[3] <= 0) {
		return 0, 0, 0, 0, false
	}
	if w <= 0 || h <= 0 {
		return 0, 0, 0, 0, false
	}
	return x, y, w, h, true
}
