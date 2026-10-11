//go:build linux && !android

package antui

func (p *pulse) next() uint32 {
	p.serial++
	return p.serial
}
