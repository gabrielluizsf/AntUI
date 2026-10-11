package antui

import (
	"github.com/gabrielluizsf/antui/canvas"
	"io"
)

// WriteWindowsResourceFor is WriteWindowsResource for an architecture other
// than the usual one.
func WriteWindowsResourceFor(w io.Writer, arch string, images ...*canvas.Canvas) error {
	pictures, err := iconPictures(images)
	if err != nil {
		return err
	}
	return writeCOFFFor(w, arch, buildResourceSection(pictures))
}
