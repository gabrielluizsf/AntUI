//go:build linux && !android

package antui

func (t *tags) u8(v byte) { t.b = append(t.b, tagU8, v) }
