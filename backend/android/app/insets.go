//go:build android

package app

// Insets is how much of the window the system is covering, in pixels.
//
// A fullscreen app gets a surface the whole size of the screen and the
// system then draws its bars **over** it — so the app owns every pixel and
// some of them are not visible. This is which ones.
type Insets struct {
	Left, Top, Right, Bottom int
	// Keyboard is how tall the on-screen keyboard is, and 0 when it is
	// down. It is reported separately because it comes and goes while the
	// bars do not, and because a layout usually wants to move rather than
	// shrink for it.
	//
	// It is only available from API 30. Below that the keyboard is included
	// in Bottom instead, and there is no way to tell the two apart.
	Keyboard int
}
