package svg

// warnNamed records one thing the drawing asked for that could not be done,
// naming the element it came from. It is the one door warnings go through,
// because a fetch running off the frame can say a picture did not arrive while
// the drawing is being painted, and a list written to from two goroutines has
// to be read from one of them at a time.
func (img *Image) warnNamed(element, format string, args ...any) {
	img.warnMu.Lock()
	img.warnings.warn(element, format, args...)
	img.warnMu.Unlock()
}
