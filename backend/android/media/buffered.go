//go:build android

package media

func (p *Player) buffered() int {
	if p.pcmTail >= p.pcmHead {
		return p.pcmTail - p.pcmHead
	}
	return len(p.pcm) - p.pcmHead + p.pcmTail
}
