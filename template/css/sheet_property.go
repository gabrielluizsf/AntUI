package css

import "strings"

// Property reads the winning raw value for one property of the element at
// rest, the way the cascade would compute it: the last matching declaration
// in stylesheet order, with any !important declaration overriding a normal one.
// It is how a template answers "what does my button's background say" without
// computing the whole style. Media queries are read against a window as tall
// as baseWidth is wide; [Sheet.StyleViewport] reads them against a real one.
func (sh *Sheet) Property(q PropertyQuery) (string, bool) {
	prop := keyword(q.Prop)
	var normalRaw string
	var impRaw string
	var hasNormal, hasImp bool
	for _, r := range sh.rules {
		if !r.Media.matches(Viewport{Width: q.Viewport.Width, Height: q.Viewport.Height}) {
			continue
		}
		match := false
		for _, sel := range r.Selectors {
			if sel.match(q.Tag, q.Classes, StateNone) {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		for _, d := range r.Decls {
			if strings.EqualFold(strings.TrimSpace(d.Prop), prop) {
				raw := strings.TrimSpace(d.Raw)
				if d.Important {
					impRaw, hasImp = raw, true
				} else {
					normalRaw, hasNormal = raw, true
				}
			}
		}
	}
	if hasImp {
		return impRaw, true
	}
	return normalRaw, hasNormal
}
