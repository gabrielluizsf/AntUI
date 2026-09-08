package canvas

// IconSizes are the sizes a window icon is offered at.
//
// Every window system picks the one it wants from what it is given, and the
// picking is nearest-neighbour in most of them — so offering a few sizes and
// letting it choose is better than offering one and letting it scale.
var IconSizes = []int{16, 24, 32, 48, 64, 128, 256}

// IconScaled is a square icon of a size, made from a picture of any size.
//
// It is [Scaled] to a square, and the filter is described there: a box
// filter, which is the one thing a large square picture reduced to a small
// square one wants.
func IconScaled(src *Canvas, size int) *Canvas {
	if src == nil || size <= 0 || src.Width <= 0 || src.Height <= 0 {
		return nil
	}
	out, err := NewCanvas(size, size)
	if err != nil {
		return nil
	}
	scaleInto(out, src)
	return out
}

// IconSet is a picture at every size a window icon is usually asked for,
// largest first — which is the order the window systems here want them in.
func IconSet(src *Canvas) []*Canvas {
	if src == nil {
		return nil
	}
	most := max(src.Width, src.Height)
	var out []*Canvas
	for i := len(IconSizes) - 1; i >= 0; i-- {
		if IconSizes[i] > most {
			continue
		}
		if scaled := IconScaled(src, IconSizes[i]); scaled != nil {
			out = append(out, scaled)
		}
	}
	if len(out) == 0 {
		if scaled := IconScaled(src, IconSizes[0]); scaled != nil {
			out = append(out, scaled)
		}
	}
	return out
}