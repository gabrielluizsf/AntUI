package apk

import (
	"fmt"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/gabrielluizsf/antui/canvas"
)

// writePNG saves a canvas.
func writePNG(path string, cv *canvas.Canvas) error {
	if cv == nil {
		return fmt.Errorf("apk: nothing to write to %s", path)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, canvasImage{cv})
}
