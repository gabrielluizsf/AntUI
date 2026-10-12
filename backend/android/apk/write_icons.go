package apk

import (
	"fmt"
	_ "image/jpeg"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
)

// WriteIcons turns one picture into every file the launcher wants, under
// dir: a bitmap at each density, a round one, and the adaptive pair that
// every device since Android 8 draws.
//
// It is exported because it is useful on its own — "antuiapk icon" is this
// function and nothing else — and because an app that has its own build
// wants the same set of files without the rest of the package.
//
// It uses antui's own scaler, which is an alpha-weighted box filter: every
// pixel of the result is the average of the block of the original that lands
// on it, weighted so that the colour of a transparent corner does not bleed
// into the edge of the shape. An icon is small and square and comes from
// something large and square, which is the one case that filter is exactly
// right for.
func WriteIcons(dir, source string, background canvas.Color) error {
	src, err := loadCanvas(source)
	if err != nil {
		return err
	}
	if src.Width < 192 || src.Height < 192 {
		// Not an error — a small source still makes every size — but a
		// launcher icon at 192 is what a modern phone shows, and scaling up
		// to it gains nothing.
		fmt.Fprintf(os.Stderr,
			"antuiapk: the icon is %dx%d; 512x512 is what the store asks for\n",
			src.Width, src.Height)
	}

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

	// The adaptive icon itself, which is a description rather than a
	// picture: a background, a foreground, and — from Android 13 — a
	// monochrome layer the system tints to match the user's wallpaper.
	// Anydpi-v26 means "every density, from API 26", so it wins over the
	// PNGs above wherever adaptive icons exist and is ignored where they do
	// not.
	folder := filepath.Join(dir, "mipmap-anydpi-v26")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return err
	}
	adaptive := `<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background" />
    <foreground android:drawable="@mipmap/ic_launcher_foreground" />
    <monochrome android:drawable="@mipmap/ic_launcher_foreground" />
</adaptive-icon>
`
	for _, name := range []string{"ic_launcher.xml", "ic_launcher_round.xml"} {
		if err := os.WriteFile(filepath.Join(folder, name), []byte(adaptive), 0o644); err != nil {
			return err
		}
	}

	values := filepath.Join(dir, "values")
	if err := os.MkdirAll(values, 0o755); err != nil {
		return err
	}
	colors := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <color name="ic_launcher_background">#%02X%02X%02X%02X</color>
</resources>
`, background.A(), background.R(), background.G(), background.B())
	return os.WriteFile(filepath.Join(values, "colors.xml"), []byte(colors), 0o644)
}
