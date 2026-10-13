package svg

// refreshBuild walks every built picture and refreshes shape data derived from
// Style. Definitions are included because a marker, mask, clip or pattern is
// painted through its own node tree rather than through Root.Kids.
func (img *Image) refreshBuild() bool {
	seen := map[*Node]bool{}
	changed := false
	var walk func(*Node)
	walk = func(n *Node) {
		if n == nil || seen[n] {
			return
		}
		seen[n] = true
		changed = refreshShape(n) || changed
		for _, kid := range n.Kids {
			walk(kid)
		}
	}
	walk(img.Root)
	for _, def := range img.markers {
		walk(def.content)
	}
	for _, def := range img.masks {
		walk(def.content)
	}
	for _, def := range img.patterns {
		walk(def.content)
	}
	for _, def := range img.clips {
		walk(def.content)
	}
	if changed {
		img.clearPatternPictures()
	}
	return changed
}
