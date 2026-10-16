package svg

// animationTick is the frame every animated picture in the drawing is showing
// at a moment on the clock, folded into one number: two paints that read the
// same tick show the same frames, so a canvas remembered from one is good for
// the other, and a tick that moved on means the canvas in the cache holds a
// frame that has ended. A drawing with nothing moving in it is always zero.
func animationTick(n *Node, now float64, tick int) int {
	if n == nil {
		return tick
	}
	if n.Pic != nil {
		tick = tick*31 + n.Pic.frameIndex(now)
	}
	for _, kid := range n.Kids {
		tick = animationTick(kid, now, tick)
	}
	return tick
}
