package svg

import "strings"

// readPatternViewBox reads `viewBox`, which establishes the pattern's
// coordinate system, overriding x, y, width, height and
// patternUnits/patternContentUnits.
func (img *Image) readPatternViewBox(p *pattern, e *element, warn func(string, ...any)) {
	if !e.hasAttr("viewBox") {
		return
	}
	vbStr := strings.TrimSpace(e.attr("viewBox"))
	vb, ok := parseViewBox(vbStr)
	if !ok {
		warn("the pattern has an invalid viewBox %q, so its picture is drawn where it was written", vbStr)
		return
	}
	p.x, p.y, p.w, p.h = vb[0], vb[1], vb[2], vb[3]
	p.boxUnits = true
	// patternContentUnits is also ignored per spec when viewBox is given, but
	// we keep the current behavior for backward compatibility.
	if e.hasAttr("patternContentUnits") {
		p.contentBox = strings.EqualFold(strings.TrimSpace(e.attr("patternContentUnits")), "objectboundingbox")
	}
	warn("the pattern has a viewBox, which sets the tile coordinates; " +
		"x, y, width and height attributes are ignored")
}
