//go:build android

package antui

func (b *androidWindow) setClipboard(string) bool { return false }
