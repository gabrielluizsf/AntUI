//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

func (x *Driver) available() int { return len(x.rbuf) - x.rpos }
