package sdk

// Toolchain is everything found. [Find] always looks for all of it — a part
// that is simply not installed is zero here, and [Require] decides only
// which of those absences is an error.
type Toolchain struct {
	Root       string // the SDK root
	BuildTools BuildTools
	Platform   Platform
	NDK        NDK
	JDK        JDK
	// Bundletool is the jar that takes an App Bundle apart again, or empty
	// when it is not on this machine. It is **not part of the SDK**: it is a
	// single jar from Google's own releases, and nothing needs it to *build*
	// a bundle — only to check one.
	Bundletool string
}
