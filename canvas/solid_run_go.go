package canvas

// solidRunGo writes one opaque colour over a whole run. The compiler turns
// this into a store loop of the widest width it can, which is the same thing
// BlendRowSolid's own opaque case did.
func solidRunGo(dst []Color, c Color) {
	solid := c | 0xFF000000
	for i := range dst {
		dst[i] = solid
	}
}
