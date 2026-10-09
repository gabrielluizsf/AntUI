package css

import "strings"

// validGridName is a line name: an identifier that is none of the keywords the
// grid grammar spells elsewhere.
func validGridName(raw string) bool {
	if !validGridIdent(raw) {
		return false
	}
	switch strings.ToLower(raw) {
	case "auto", "span", "min-content", "max-content", "fit-content", "minmax",
		"repeat", "subgrid", "initial", "inherit", "unset", "none":
		return false
	}
	return true
}
