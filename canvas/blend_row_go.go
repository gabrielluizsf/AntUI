package canvas

// blendRowGo is the reference implementation, and the one that runs where
// there is no assembly for the architecture.
func blendRowGo(dst, src []Color) {
	for i, c := range src {
		switch c.A() {
		case 0:
			// Nothing to add, and dst keeps its own alpha.
		case 255:
			dst[i] = c
		default:
			dst[i] = Blend(dst[i], c)
		}
	}
}
