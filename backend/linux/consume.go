//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

func (x *Driver) consume(count int) { x.rpos += count }
