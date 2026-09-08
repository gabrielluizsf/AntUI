//go:build linux || freebsd || netbsd || openbsd || dragonfly

package antui

// defaultUIPoints is what a desktop toolkit draws its interface at when
// nobody has said otherwise. GNOME and KDE both start at 10 or 11.
const defaultUIPoints = 11

// systemFontPaths is where to look when fontconfig cannot be asked — a
// container with no fc-match in it, a machine with fontconfig stripped, an
// embedded system.
//
// It is a fallback and not the method. These are the files a distribution is
// likely to have resolved sans-serif to anyway, so the answer is usually the
// same one; when it is not, asking was right and this is a guess.
func systemFontPaths() []string {
	return []string{
		"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/google-noto/NotoSans-Regular.ttf",
		"/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/liberation-sans/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/usr/local/share/fonts/dejavu/DejaVuSans.ttf",
	}
}