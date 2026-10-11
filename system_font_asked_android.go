//go:build android

package antui

import (
	"encoding/xml"
	"os"
)

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
	for _, path := range fontsXML {
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var set familySet
		if err := xml.Unmarshal(body, &set); err != nil {
			continue
		}
		if out := askedFontFiles(set); len(out) > 0 {
			return out
		}
	}
	return nil
}
