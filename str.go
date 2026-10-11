//go:build linux && !android

package antui

// str writes a string, which carries its terminator: the protocol has no
// length in front of one.
func (t *tags) str(s string) {
	t.b = append(t.b, tagString)
	t.b = append(t.b, s...)
	t.b = append(t.b, 0)
}
