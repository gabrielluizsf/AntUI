package canvas

// discCoverage is the antialiased coverage of a disc, 0 outside to 255 in.
//
// The square root is only worth taking in the one-pixel band along the edge.
// Inside it the answer is 255 and outside it is 0, and both are decided by
// comparing squares: a circle is mostly its own interior, and thirty-two
// shifts per pixel to learn what a multiply would have said is most of what
// drawing one costs.
func discCoverage(dx, dy, radius int) int {
	dist2 := uint64(dx*dx + dy*dy)
	edge := int64(radius)*256 + 128
	// floor(256*sqrt(dist2)) is at most t exactly when dist2*65536 < (t+1)^2,
	// which for a whole dist2 is dist2 < ceil((t+1)^2/65536).
	if dist2 >= ceilSquare(uint64(edge)) {
		return 0
	}
	if dist2 < ceilSquare(uint64(edge-255)) {
		return 255
	}
	dist := int64(isqrt64(dist2 << 16))
	cover := edge - dist
	if cover <= 0 {
		return 0
	}
	if cover >= 256 {
		return 255
	}
	return int(cover)
}
