package canvas

// ellipseCoverage is [discCoverage] for an axis-aligned ellipse: it measures
// the antialiased coverage in the ellipse's own units, so a flat rx=ry ellipse
// blurs exactly like a disc of the same radius.
func ellipseCoverage(dx, dy, rx, ry int) int {
	if rx <= 0 || ry <= 0 {
		return 0
	}
	// Normalise the point into the ellipse's unit space: (dx/rx)^2+(dy/ry)^2
	// scaled to integers is ((dx*ry)^2+(dy*rx)^2) / (rx*ry)^2.
	ux := int64(dx) * int64(ry)
	uy := int64(dy) * int64(rx)
	dist2 := uint64(ux*ux + uy*uy)
	den := int64(rx) * int64(ry)
	edge := den*256 + 128

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
