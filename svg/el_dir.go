package svg

import "strings"

// elDir answers the direction of writing the element at node stands in: the
// nearest `dir` above it that gives one, where a value of `auto` is settled by
// the first strong directional character the element writes. With none saying
// otherwise the drawing runs left to right, as a document with nothing to say
// about it does.
func elDir(n *selNode) string {
	for m := n; m != nil; m = m.parent {
		switch v := strings.ToLower(m.el.Attr["dir"]); v {
		case "ltr", "rtl":
			return v
		case "auto":
			return firstStrongDir(m.el)
		}
	}
	return "ltr"
}
