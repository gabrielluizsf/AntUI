//go:build android

package display

func (e errString) Error() string { return string(e) }
