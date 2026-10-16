package svg

// warn records one thing the drawing asked for that could not be done.
func (img *Image) warn(e *element, format string, args ...any) {
	img.warnNamed(e.Name, format, args...)
}
