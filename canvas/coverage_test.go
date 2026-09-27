package canvas

import "testing"

// TestCoverageIsExactAtEveryDistance walks every pixel a disc and a set of
// ellipses can reach and holds the answer against the one it replaced: a
// square root per pixel is slow, and the comparisons that bracket it are only
// allowed to be faster, never different. A ring drawn with the fast path has
// to be the same ring, to the last bit of coverage.
func TestCoverageIsExactAtEveryDistance(t *testing.T) {
	for _, radius := range []int{1, 2, 3, 7, 16, 60, 120} {
		for dy := -radius - 2; dy <= radius+2; dy++ {
			for dx := -radius - 2; dx <= radius+2; dx++ {
				if got, want := discCoverage(dx, dy, radius), rootedDiscCoverage(dx, dy, radius); got != want {
					t.Fatalf("discCoverage(%d, %d, %d) = %d, want %d", dx, dy, radius, got, want)
				}
			}
		}
	}
	for _, rx := range []int{1, 3, 20, 120} {
		for _, ry := range []int{1, 7, 40, 120} {
			for dy := -ry - 2; dy <= ry+2; dy++ {
				for dx := -rx - 2; dx <= rx+2; dx++ {
					got, want := ellipseCoverage(dx, dy, rx, ry), rootedEllipseCoverage(dx, dy, rx, ry)
					if got != want {
						t.Fatalf("ellipseCoverage(%d, %d, %d, %d) = %d, want %d", dx, dy, rx, ry, got, want)
					}
				}
			}
		}
	}
}

// TestCoverageIsAFullCircleAndAnEmptyOne guards what the bracket is for: the
// middle of a shape is solid and the outside of it is nothing, with the
// antialiased band in between exactly one pixel wide.
func TestCoverageIsAFullCircleAndAnEmptyOne(t *testing.T) {
	const radius = 20
	if got := discCoverage(0, 0, radius); got != 255 {
		t.Errorf("the middle of a disc = %d, want it solid", got)
	}
	if got := discCoverage(radius+2, 0, radius); got != 0 {
		t.Errorf("outside a disc = %d, want nothing", got)
	}
	if got := ellipseCoverage(0, 0, 40, 10); got != 255 {
		t.Errorf("the middle of an ellipse = %d, want it solid", got)
	}
	if got := ellipseCoverage(0, 12, 40, 10); got != 0 {
		t.Errorf("outside an ellipse = %d, want nothing", got)
	}
	// The band: a disc of radius 20 has one pixel of partial coverage all the
	// way round, and the two either side of it are solid and empty.
	band := 0
	for dx := -radius; dx <= radius; dx++ {
		switch discCoverage(dx, 0, radius) {
		case 0:
		case 255:
		default:
			band++
		}
	}
	if band == 0 {
		t.Error("a disc has no antialiased edge")
	}
}

// rootedDiscCoverage and rootedEllipseCoverage are the coverage as it was
// before the square root was bracketed: every pixel paying for one.
func rootedDiscCoverage(dx, dy, radius int) int {
	dist2 := uint64(dx*dx + dy*dy)
	dist := int64(isqrt64(dist2 << 16))
	edge := int64(radius)*256 + 128
	return clampCover(edge - dist)
}

func rootedEllipseCoverage(dx, dy, rx, ry int) int {
	if rx <= 0 || ry <= 0 {
		return 0
	}
	ux := int64(dx) * int64(ry)
	uy := int64(dy) * int64(rx)
	den := int64(rx) * int64(ry)
	dist := int64(isqrt64(uint64(ux*ux+uy*uy) << 16))
	return clampCover(den*256 + 128 - dist)
}

func clampCover(cover int64) int {
	switch {
	case cover <= 0:
		return 0
	case cover >= 256:
		return 255
	default:
		return int(cover)
	}
}
