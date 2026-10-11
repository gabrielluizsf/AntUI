package font

// isPNGImage reports whether a CBDT image format holds a PNG, which is the
// only colour picture this engine reads.
func isPNGImage(format int) bool { return format == 17 || format == 18 || format == 19 }
