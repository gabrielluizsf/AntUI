package svg

import "strings"

// readPreserveSpace reads `xml:space`: only the two words mean anything, and
// anything else leaves the inherited answer alone rather than switching the
// writing off its own whitespace by accident.
func readPreserveSpace(raw string, inherited bool) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "preserve":
		return true
	case "default":
		return false
	}
	return inherited
}
