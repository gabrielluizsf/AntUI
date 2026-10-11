package canvas

// boxBlurLine is one row or one column of a pass. The four channels are
// averaged side by side rather than one after another — each keeps its own
// running sum and its own order of additions, so the numbers are the same ones
// four separate walks would have produced, and the line is read once instead of
// four times.
//
// The window reaches past both ends of the line and has to be clamped there,
// which is two samples of bookkeeping for every pixel of a shadow. The middle
// of the line, where the window stays inside, skips it.
func boxBlurLine(src, dst []float64, at, stride, count, r int, span float64) {
	var acc blurSum
	for k := -r; k <= r; k++ {
		i := at + clampInt(k, 0, count-1)*stride
		acc.r += src[i]
		acc.g += src[i+1]
		acc.b += src[i+2]
		acc.a += src[i+3]
	}

	// ends is the run of samples whose window reaches past the start of the
	// line, mid is where it stops reaching past the end, and everything between
	// them is a sample both of whose ends are inside the line: n-r is a real
	// sample from n=r on, and n+r+1 is one until n = count-r-2.
	ends := min(r, count)
	mid := min(max(count-r-1, ends), count)

	for n := range ends {
		o := at + n*stride
		acc.push(dst, o, span)
		lo := at + clampInt(n-r, 0, count-1)*stride
		hi := at + clampInt(n+r+1, 0, count-1)*stride
		acc.slide(src, lo, hi)
	}
	for n := ends; n < mid; n++ {
		o := at + n*stride
		acc.push(dst, o, span)
		acc.slide(src, o-r*stride, o+(r+1)*stride)
	}
	for n := mid; n < count; n++ {
		o := at + n*stride
		acc.push(dst, o, span)
		lo := at + clampInt(n-r, 0, count-1)*stride
		hi := at + clampInt(n+r+1, 0, count-1)*stride
		acc.slide(src, lo, hi)
	}
}
