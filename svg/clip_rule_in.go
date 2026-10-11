package svg

import "github.com/gabrielluizsf/antui/canvas"

// clipRuleIn is the clip-rule one element stands under: the one written on it
// — as an attribute or as a declaration in its `style` — or the one it
// inherited from above, which is what makes clip-rule an inherited property
// like the rest of the style. It is read without a warn because there is
// nothing here that cannot be read: a rule that is not `evenodd` is simply
// the nonzero one, which is also what it starts as.
func clipRuleIn(e *element, inherited canvas.FillRule) canvas.FillRule {
	if e.hasAttr("clip-rule") {
		return clipRuleOf(e.attr("clip-rule"))
	}
	if s := e.attr("style"); s != "" {
		if raw, ok := parseDeclarations(s)["clip-rule"]; ok {
			return clipRuleOf(raw)
		}
	}
	return inherited
}
