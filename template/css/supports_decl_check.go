package css

import (
	"strings"
)

// supportsDecl reports whether the engine would take a declaration: a
// property it has, holding a value it reads. Both halves come from applyDecl
// itself — the fold the cascade runs — so @supports can only claim what the
// engine will really do with the declaration.
func supportsDecl(prop, raw string) bool {
	prop = strings.ToLower(strings.TrimSpace(prop))
	raw = strings.TrimSpace(raw)
	if prop == "" || raw == "" {
		return false
	}
	if strings.HasPrefix(prop, "--") {
		// A custom property is the author's to declare: the engine holds it
		// for var() without being asked what it holds.
		return true
	}
	if strings.HasPrefix(prop, "-") {
		// A vendor prefix the engine has no case for, which applyDecl says
		// out loud whenever it meets one in a stylesheet.
		return false
	}
	if !supportedProperty(prop) {
		return false
	}
	_, set := foldSupport(prop, raw)
	return set[prop]
}

// supportedProperty reports whether the engine has the property at all. It
// asks with a value no grammar reads, where the only answer that can come
// back is about the property's own name: a value the engine cannot read is
// dropped in silence, and only a name it has no case for is reported.
func supportedProperty(prop string) bool {
	warns, _ := foldSupport(prop, unreadableValue)
	return len(warns) == 0
}

const unreadableValue = "?"
