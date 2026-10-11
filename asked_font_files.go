//go:build android

package antui

import (
	"path/filepath"
	"strings"
)

// askedFontFiles picks the regular upright files from a parsed font config.
func askedFontFiles(set familySet) []string {
	var out []string
	for _, family := range set.Families {
		if family.Name != "" && family.Name != "sans-serif" {
			continue
		}
		for _, font := range family.Fonts {
			// The regular weight, upright. A file with no weight at all is
			// taken too: a variable font is listed once and covers every
			// weight, which is how the newest devices write it.
			if font.Style != "" && font.Style != "normal" {
				continue
			}
			if font.Weight != "" && font.Weight != "400" {
				continue
			}
			name := strings.TrimSpace(font.File)
			if name == "" {
				continue
			}
			if !filepath.IsAbs(name) {
				name = filepath.Join(fontsDir, name)
			}
			out = append(out, name)
		}
	}
	return out
}
