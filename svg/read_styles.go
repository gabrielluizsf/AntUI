package svg

// readStyles reads the drawing's own `<style>` sheets and leaves on every
// element what the rules matching it say. It runs before anything is built,
// because a sheet reaches every element of the drawing — the shapes a clip, a
// mask, a marker or a pattern holds as much as the ones painted where they
// stand — and the style read afterwards is the whole cascade: the sheet by
// weight, then the element's own `style`, and the presentation attributes
// under both.
func (img *Image) readStyles(root *element) {
	rules := img.styleRules(root, namespace{})
	if len(rules) == 0 {
		return
	}
	img.matchRules(selTree(root, nil, nil), rules)
}

// styleRules gathers the rules of every `<style>` in the drawing, in the order
// they were written, which is the order CSS breaks a tie by. ns is the
// namespaces in scope at each element, so a selector with a prefix is read
// through the ones the drawing itself binds where the sheet stands.
func (img *Image) styleRules(e *element, ns namespace) []cssRule {
	inner, err := ns.child(e)
	if err != nil {
		inner = ns
	}
	var rules []cssRule
	if e.Name == "style" {
		rules = append(rules, parseSheet(e.Text, inner, func(f string, a ...any) { img.warn(e, f, a...) })...)
	}
	for _, k := range e.Kids {
		rules = append(rules, img.styleRules(k, inner)...)
	}
	return rules
}

// matchRules walks the drawing and leaves the matched declarations on each
// element.
func (img *Image) matchRules(n *selNode, rules []cssRule) {
	if d := matchDecls(n, rules); d != nil {
		n.el.sheet = d
	}
	for _, k := range n.kids {
		img.matchRules(k, rules)
	}
}
