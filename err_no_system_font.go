package antui

import (
	"errors"
)

// ErrNoSystemFont is what [SystemFace] gives back when the system has no
// answer and nothing on the fallback list is there.
//
// It is worth telling apart from a broken font file, because the response is
// the same either way — carry on with the built-in — and a program should
// not stop for it.
var ErrNoSystemFont = errors.New("antui: no system font found")
