package apk

// symbolFile is where BuildLibs leaves the symbols for a library: beside it,
// with .sym on the end. It is a sibling rather than a separate directory so
// that a caller who keeps the build directory has them both together.
func symbolFile(lib string) string { return lib + symbolSuffix }
