//go:build linux && !android

package antui

import "errors"

var errShort = errors.New("antui: the server's answer stopped in the middle")
