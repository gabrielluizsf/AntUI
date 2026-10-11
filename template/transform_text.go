package template

import (
	"strings"

	"github.com/gabrielluizsf/antui/template/css"
)

// transformText applies text-transform to a run.
func transformText(s string, t uint8) string {
	switch t {
	case css.TextTransformUppercase:
		return strings.ToUpper(s)
	case css.TextTransformLowercase:
		return strings.ToLower(s)
	case css.TextTransformCapitalize:
		return capitalize(s)
	}
	return s
}
