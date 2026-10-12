package apk

import (
	"os"
	"path/filepath"
)

// writeAdaptive writes the adaptive icon, which is a description rather than
// a picture: a background, a foreground, and — from Android 13 — a
// monochrome layer the system tints to match the user's wallpaper.
// anydpi-v26 means "every density, from API 26", so it wins over the PNGs
// wherever adaptive icons exist and is ignored where they do not.
func writeAdaptive(dir string) error {
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
	return nil
}
