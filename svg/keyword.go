package svg

import "strings"

// keyword folds a keyword read out of a drawing into the form the parsers
// compare against: the surrounding space gone and the case gone, since XML
// attribute values do not care about either.
func keyword(raw string) string { return strings.ToLower(strings.TrimSpace(raw)) }
