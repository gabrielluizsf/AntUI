package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
	"io"
)

// WriteWindowsResource writes a .syso: an object file the Go linker puts into
// the program, holding the icon the file browser shows.
//
// Put what it writes in the main package as `icon_windows.syso` — the name
// ends in `_windows` so that the toolchain leaves it out of every other
// platform's build, which is what makes a checked-in resource harmless
// everywhere else.
func WriteWindowsResource(w io.Writer, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}
	section := buildResourceSection(pictures)
	return writeCOFF(w, section)
}
