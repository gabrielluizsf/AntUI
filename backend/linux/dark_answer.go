//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// darkAnswer is one source's word: ok says whether it had any, dark what it
// said. A source that had none leaves ok false rather than claiming light,
// and the question moves on to the next one.
type darkAnswer struct {
	dark bool
	ok   bool
}
