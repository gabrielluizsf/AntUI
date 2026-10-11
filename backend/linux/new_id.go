//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

func (x *Driver) newID() uint32 {
	id := x.idBase | (x.idNext & x.idMask)
	x.idNext++
	return id
}
