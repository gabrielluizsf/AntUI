//go:build android

package antui

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// defaultUIPoints is what Android draws its own interface at: 14 scaled
// pixels, which is points here because the display's density is applied
// separately.
const defaultUIPoints = 14

// fontsXML is the platform's own font configuration, in the order the
// platform has kept it.
//
// /system/etc/fonts.xml is where it has been since API 21. The one before it
// is listed because a device that old still runs this library, and /product
// is where a manufacturer puts its own — a phone whose maker replaced the
// interface font says so there and nowhere else.
var fontsXML = []string{
	"/system/etc/fonts.xml",
	"/product/etc/fonts.xml",
	"/system/etc/system_fonts.xml",
}

// fontsDir is what the names in that file are relative to.
const fontsDir = "/system/fonts"

// familySet is as much of fonts.xml as this needs.
type familySet struct {
	Families []struct {
		Name  string `xml:"name,attr"`
		Fonts []struct {
			Weight string `xml:"weight,attr"`
			Style  string `xml:"style,attr"`
			File   string `xml:",chardata"`
		} `xml:"font"`
	} `xml:"family"`
}

// systemFontAsked reads Android's own font configuration rather than
// assuming a filename.
//
// **The file says not to.** Its own comment warns that third-party apps
// parsing it will "almost certainly break with the next major Android
// release". That is a fair warning and it is why [systemFontPaths] is there
// behind this: reading the platform's answer is right when it works, and
// when the format moves this falls through to the files that have been in
// /system/fonts for a decade, and neither case is a crash or a wrong font.
//
// The interface font is the "sans-serif" family, and which file that is has
// changed more than once: Droid Sans, then Roboto as one file per weight,
// then Roboto as a single variable font. A manufacturer may have replaced it
// entirely — plenty do — and the only place that is written down is this
// file.
//
// The first family with no name is the default one, which is what the
// platform falls back to and what older versions of this file use instead of
// naming sans-serif.
func systemFontAsked() []string {
	var out []string
	for _, path := range fontsXML {
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var set familySet
		if err := xml.Unmarshal(body, &set); err != nil {
			continue
		}
		for _, family := range set.Families {
			if family.Name != "" && family.Name != "sans-serif" {
				continue
			}
			for _, font := range family.Fonts {
				// The regular weight, upright. A file with no weight at all
				// is taken too: a variable font is listed once and covers
				// every weight, which is how the newest devices write it.
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
		if len(out) > 0 {
			return out
		}
	}
	return out
}

// systemFontPaths is the fallback, for a device whose configuration file
// cannot be read or does not parse.
func systemFontPaths() []string {
	return []string{
		"/system/fonts/Roboto-Regular.ttf",
		"/system/fonts/RobotoStatic-Regular.ttf",
		"/system/fonts/NotoSans-Regular.ttf",
		"/system/fonts/DroidSans.ttf",
	}
}
