package svg

import "strings"

// keepsSpace is whether the writing inside an element is drawn exactly as the
// file wrote it. `xml:space="preserve"` turns it on and `xml:space="default"`
// turns it off, and with neither written the answer the element above gave is
// the answer here — which is what makes the attribute worth being inherited.
func keepsSpace(e *element, inherited bool) bool {
	switch strings.ToLower(strings.TrimSpace(e.attr("xml:space"))) {
	case "preserve":
		return true
	case "default":
		return false
	}
	return inherited
}
