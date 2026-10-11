package canvas

// blurSum is a running sum of the four channels of a line as the blur window
// slides along it. Adding src[hi]-src[lo] would be one addition fewer than the
// subtraction then addition the walk was written with, and a different answer
// in the last bit, so slide keeps the original order.
type blurSum struct{ r, g, b, a float64 }

func (s *blurSum) push(dst []float64, o int, span float64) {
	dst[o] = s.r / span
	dst[o+1] = s.g / span
	dst[o+2] = s.b / span
	dst[o+3] = s.a / span
}

func (s *blurSum) slide(src []float64, lo, hi int) {
	s.r = s.r - src[lo] + src[hi]
	s.g = s.g - src[lo+1] + src[hi+1]
	s.b = s.b - src[lo+2] + src[hi+2]
	s.a = s.a - src[lo+3] + src[hi+3]
}
