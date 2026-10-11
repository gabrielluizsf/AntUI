//go:build android

package antui

// systemDark is the device's night mode, which Android carries in every
// configuration it hands the app and so always has an answer to give.
func (b *androidWindow) systemDark() (bool, bool) { return b.dark, true }
