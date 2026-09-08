package antui

// defaultUIPoints is what macOS draws its own interface at.
const defaultUIPoints = 13

// systemFontPaths is where to look when fontconfig is not installed, which
// on macOS is the usual case: it is San Francisco, which macOS has used
// since 10.11, and Helvetica behind it.
//
// **CoreText is the right answer here and is not done.**
// CTFontCreateUIFontForLanguage names the interface font and
// kCTFontURLAttribute gives its file, which is the system telling you rather
// than a list agreeing with it — and this library already allows cgo on
// macOS, so nothing is in the way except that it cannot be compiled or run
// on the machine this was written on. Saying so is better than pretending
// the list is the method.
//
// SFNS is a variable font in a .ttf, and what is read here is its default
// instance — the regular weight — which is what an interface is written in
// anyway. Anything that needs a weight other than regular needs variable
// font support, which is not here.
func systemFontPaths() []string {
	return []string{
		"/System/Library/Fonts/SFNS.ttf",
		"/System/Library/Fonts/SFNSText.ttf",
		"/System/Library/Fonts/SFNSDisplay.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/System/Library/Fonts/HelveticaNeue.ttc",
		"/Library/Fonts/Arial.ttf",
	}
}