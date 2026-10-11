//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// askScale is the display's scale, read straight from the server.
//
// X11 has no answer to this, which is the honest summary — there is a
// physical size in the screen record and it is very often a lie, since a
// server with no EDID makes one up that lands on exactly 96 dpi. What there
// is instead is a convention: every toolkit reads Xft.dpi out of the X
// resource database, and that is what the user's desktop settings write. So
// that is read first and the physical size is the fallback, believed only
// when it lands somewhere a real display could be.
func (x *Driver) askScale() (float64, bool) {
	if dpi := x.xftDPI(); dpi > 0 {
		return dpi / 96, true
	}
	if x.screenWidth > 0 && x.screenWidthMM > 0 {
		dpi := float64(x.screenWidth) * 25.4 / float64(x.screenWidthMM)
		// Under 96 is a television or a made-up number, and over 400 is a
		// screen size the server got wrong by a factor. Neither is a scale
		// anyone chose, and guessing wrong here makes every window the wrong
		// size — so an unbelievable answer is no answer.
		if dpi >= 96 && dpi <= 400 {
			return dpi / 96, true
		}
	}
	return 0, false
}
