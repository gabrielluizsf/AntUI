//go:build linux && !android

package antui

func (t *tags) boolean(v bool) {
	if v {
		t.b = append(t.b, tagTrue)
		return
	}
	t.b = append(t.b, tagFalse)
}
