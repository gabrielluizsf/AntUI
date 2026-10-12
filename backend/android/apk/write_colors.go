package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
)

// writeColors writes the background colour the adaptive icon refers to.
func writeColors(dir string, background canvas.Color) error {
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
