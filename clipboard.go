//go:build android

package antui

// The clipboard is Java's, and arrives with the JNI layer.
func (b *androidWindow) clipboard() (string, []string) { return "", nil }
