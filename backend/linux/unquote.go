//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"strings"
)

// unquote takes the quotes the settings tool puts around the value it
// prints, and the whitespace it puts around that.
func unquote(value string) string {
	return strings.Trim(strings.TrimSpace(value), `'"`)
}
