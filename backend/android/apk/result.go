package apk

import (
	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Result is what a build produced.
type Result struct {
	// APK is the file written.
	APK string
	// Libs is the compiled library for each ABI, in the temporary directory
	// the build used. They are gone by the time Build returns unless Keep
	// was set.
	Libs map[sdk.ABI]string
	// Size is the finished package, in bytes.
	Size int64
	// Symbols is the archive of symbol files written beside the package, or
	// empty when the build was told not to strip and the symbols are
	// therefore still in the shipped library. See [WriteSymbols].
	Symbols string
}
