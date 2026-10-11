//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// DisplayScale is how much the display is scaled by, and whether it said
// one. The number is read from the server when the desktop writes it and
// kept in between: this is asked on every frame of a window that is up, and
// asking the server for it is a round trip a frame would spend on an answer
// that changes when the user changes it. [Driver.watchScale] is what makes
// that change reach here.
func (x *Driver) DisplayScale() (float64, bool) {
	if !x.alive {
		return 0, false
	}
	if x.scaleDirty {
		// Cleared before the read rather than after: the read is a round
		// trip, and a desktop that rewrites the database while it is in
		// flight must leave the number dirty for the next frame rather
		// than have this one put it back.
		x.scaleDirty = false
		x.scale, x.scaleKnown = x.askScale()
	}
	return x.scale, x.scaleKnown
}
