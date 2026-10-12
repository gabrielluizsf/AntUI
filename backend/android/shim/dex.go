package shim

import _ "embed"

// Dex is the compiled shim, ready to be put in an APK as classes.dex.
//
//go:embed classes.dex
var Dex []byte
