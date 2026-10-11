//go:build linux && !android

package antui

func (t *tags) noStr() { t.b = append(t.b, tagStringNil) }
