package svg

import "strings"

// readPatternHref follows an href to another pattern, the same way gradients
// follow href references. The referenced pattern must be named with an id and
// be within this drawing; its content and transform are taken, and its
// boxUnits and contentBox settings are inherited.
func (img *Image) readPatternHref(p *pattern, e *element, warn func(string, ...any)) {
	for _, name := range []string{"href", "xlink:href"} {
		ref, ok := hashRef(e.attr(name))
		if !ok || ref == "" {
			continue
		}
		id := strings.TrimPrefix(ref, "#")
		if refPattern, ok := img.patterns[id]; ok {
			p.content = refPattern.content
			p.transform = refPattern.transform
			p.boxUnits = refPattern.boxUnits
			p.contentBox = refPattern.contentBox
			warn("the pattern takes its content from #%s", id)
		} else {
			warn("the pattern references #%s, which this drawing does not have, so its content is empty", id)
		}
		break
	}
}
