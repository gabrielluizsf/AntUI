//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

type xAuth struct {
	name string
	data []byte
}
