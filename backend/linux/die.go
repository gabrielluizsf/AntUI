//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// die gives up on the connection, remembering why.
func (x *Driver) die(reason string) {
	if x.alive {
		x.alive = false
		x.dead = reason
	}
}
