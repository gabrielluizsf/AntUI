package antui

import (
	"github.com/gabrielluizsf/antui/font"
	"os"
)

// SystemFont is the file this platform writes its interface in, and the
// bytes of it.
//
// The path is there so a program can say which font it chose — a settings
// screen showing "Noto Sans" is a great deal more use than one showing
// nothing — and so that a bug report can name it.
func SystemFont() (path string, data []byte, err error) {
	sysOnce.Do(func() {
		// Asked for first, guessed at second.
		for _, candidate := range append(systemFontAsked(), systemFontPaths()...) {
			body, err := os.ReadFile(candidate)
			if err != nil {
				continue
			}
			// Read rather than merely found: a path that exists and is not a
			// TrueType file — a CFF .otf, a link to nothing, a stub — is no
			// more use than a missing one, and the next candidate is better
			// than failing here.
			if _, err := font.ParseTTF(body); err != nil {
				continue
			}
			sysPath, sysData = candidate, body
			return
		}
		sysErr = ErrNoSystemFont
	})
	return sysPath, sysData, sysErr
}
