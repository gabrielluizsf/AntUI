package css

import "strings"

// applyFaceWeightStyle parses the weight and slant declarations onto the
// face, warning and keeping the defaults on anything unreadable.
func applyFaceWeightStyle(sh *Sheet, ff *FontFace, weight, style, family string) {
	if s := strings.TrimSpace(weight); s != "" {
		w, ok := parseFontWeight(s)
		if !ok {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring font-weight %q in @font-face for %q", weight, family).Error())
		} else {
			ff.Weight = w
		}
	}
	if s := strings.TrimSpace(style); s != "" {
		v, ok := parseFontStyle(s)
		if !ok {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring font-style %q in @font-face for %q", style, family).Error())
		} else {
			ff.Style = v
		}
	}
}
