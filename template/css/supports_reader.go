package css

import (
	"strings"
)

// supportsReader walks an @supports prelude. A condition is an or-list of
// and-lists of negations, each of which is either a parenthesised group — a
// declaration when a colon splits it, a nested condition when it does not — or
// one of the condition's own functions.
type supportsReader struct {
	src string
	i   int
	bad bool
}

func (r *supportsReader) eof() bool { return r.i >= len(r.src) }

func (r *supportsReader) skipSpace() {
	for !r.eof() && isSpace(r.src[r.i]) {
		r.i++
	}
}

// atEnd reports the reader having consumed the whole prelude.
func (r *supportsReader) atEnd() bool {
	r.skipSpace()
	return r.eof()
}

// word takes w when it is the next thing, and leaves the reader past the
// spaces before it when it is not, so a word that does not match costs
// nothing. Keywords are read without case, the way the media grammar reads
// them.
func (r *supportsReader) word(w string) bool {
	r.skipSpace()
	rest := r.src[r.i:]
	if len(rest) < len(w) || !strings.EqualFold(rest[:len(w)], w) {
		return false
	}
	if r.i+len(w) < len(r.src) && isIdentByte(r.src[r.i+len(w)]) {
		return false // a longer word that merely starts with w
	}
	r.i += len(w)
	return true
}
