package svg

import "fmt"

// structuralPseudo is the pseudo-classes about where an element stands among
// its siblings that need no argument.
var structuralPseudo = map[string]simpleKind{
	"root":          kindRoot,
	"scope":         kindScope,
	"empty":         kindEmpty,
	"first-child":   kindFirstChild,
	"last-child":    kindLastChild,
	"only-child":    kindOnlyChild,
	"first-of-type": kindFirstOfType,
	"last-of-type":  kindLastOfType,
	"only-of-type":  kindOnlyOfType,
}

// appendPseudoClass reads a `:name` that is not a pseudo-element: an argument
// is read for the functional ones, and a name CSS does not know is refused, so
// an unreadable selector is left out rather than matching something it did not
// mean. ns is the namespaces in scope, which an argument carrying a selector
// reads through.
func appendPseudoClass(c *compoundSelector, s, name string, next, at int, ns namespace) (int, error) {
	if k, ok := structuralPseudo[name]; ok {
		c.simples = append(c.simples, simpleSelector{kind: k})
		return next, nil
	}
	switch name {
	case "nth-child", "nth-last-child", "nth-of-type", "nth-last-of-type":
		return appendNth(c, s, name, next, at, ns)
	case "not", "is", "where":
		return appendLogical(c, s, name, next, at, ns)
	case "has":
		return appendHas(c, s, next, at, ns)
	case "lang":
		return appendLang(c, s, next, at)
	case "dir":
		return appendDir(c, s, next, at)
	}
	if isPseudoElement(name) {
		return appendUnmatchableElement(c, s, next, at)
	}
	if neverMatches[name] {
		return appendUnmatchable(c, s, next, at)
	}
	return at, fmt.Errorf("the pseudo-class :%s is not one CSS names", name)
}
