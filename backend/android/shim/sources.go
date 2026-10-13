package shim

import "embed"

// Sources is the Java it was compiled from — every file, so that [Build] can
// reproduce the dex without knowing where the repository is.
//
//go:embed *.java
var Sources embed.FS
