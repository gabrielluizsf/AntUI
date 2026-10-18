package svg

import "strings"

// parseSheet reads the CSS of a `<style>` element into the rules it holds. A
// rule is a selector list and a block of declarations. An at-rule is read past
// whole — a still drawing has no screen for `@media` to query and no sheet for
// `@import` to bring in — and is warned about rather than half-applied. ns is
// the namespaces in scope where the `<style>` stands, which a selector with a
// prefix is read through.
func parseSheet(css string, ns namespace, warn func(string, ...any)) []cssRule {
	css = stripCSSComments(css)
	var rules []cssRule
	for i := 0; i < len(css); {
		open := strings.IndexByte(css[i:], '{')
		if open < 0 {
			break
		}
		open += i
		prelude := strings.TrimSpace(css[i:open])
		end := blockEnd(css, open)
		if end < 0 {
			warn("a rule in a style sheet is never closed, so %q is left unread", prelude)
			break
		}
		body := css[open+1 : end]
		i = end + 1
		if prelude == "" {
			continue
		}
		if prelude[0] == '@' {
			warn("the at-rule @%s is not one a drawing applies, so it is left out", atRuleName(prelude))
			continue
		}
		selectors, err := parseSelectorList(prelude, ns)
		if err != nil {
			warn("the selector %q is not one a drawing can read, so it is left out: %v", prelude, err)
			continue
		}
		if len(selectors) > 0 {
			rules = append(rules, cssRule{selectors: selectors, decls: parseDecls(body)})
		}
	}
	return rules
}

// atRuleName is the name an at-rule is known by, `media` in `@media …`.
func atRuleName(prelude string) string {
	prelude = prelude[1:]
	i := strings.IndexAny(prelude, " \t\r\n(")
	if i < 0 {
		return prelude
	}
	return prelude[:i]
}
