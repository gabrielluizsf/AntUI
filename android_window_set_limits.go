//go:build android

package antui

// setLimits and setSize: the size of the window is the system's decision and
// is not negotiable. Reporting false is the honest answer, and window_size.go
// is written to accept it.
func (b *androidWindow) setLimits(*Window, Limits) {}
