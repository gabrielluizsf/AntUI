package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
)

// writeDensities writes the square, round and adaptive-foreground PNG at
// every density a launcher asks for.
func writeDensities(dir string, src *canvas.Canvas) error {
	for _, d := range densities {
		folder := filepath.Join(dir, "mipmap-"+d.name)
		if err := os.MkdirAll(folder, 0o755); err != nil {
			return err
		}
		square := canvas.IconScaled(src, d.size)
		if square == nil {
			return fmt.Errorf("apk: cannot scale the icon to %d", d.size)
		}
		if err := writePNG(filepath.Join(folder, "ic_launcher.png"), square); err != nil {
			return err
		}
		// The round one is the same picture with the corners taken off. A
		// launcher that asks for it is one that will not mask what it is
		// given, so the masking has to be done here.
		if err := writePNG(filepath.Join(folder, "ic_launcher_round.png"),
			rounded(square)); err != nil {
			return err
		}
		if err := writePNG(filepath.Join(folder, "ic_launcher_foreground.png"),
			adaptiveLayer(src, d.adaptive)); err != nil {
			return err
		}
	}
	return nil
}
