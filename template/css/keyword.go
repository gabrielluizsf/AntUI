package css

import "strings"

// keyword folds a CSS keyword into the form the parsers compare against: the
// surrounding space gone and the case gone, since CSS keywords and property
// names do not care about either. Every `toLower(trim(...))` in the package is
// this.
func keyword(raw string) string { return strings.ToLower(strings.TrimSpace(raw)) }
