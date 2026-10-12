//go:build android

package app

var eventNames = [...]string{
	"None",
	"Start", "Resume", "Pause", "Stop", "Destroy",
	"WindowUp", "WindowResized", "Redraw", "WindowDown",
	"FocusGained", "FocusLost",
	"ContentRect", "ConfigChanged", "LowMemory",
	"InputUp", "InputDown",
}
