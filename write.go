//go:build linux && !android

package antui

import (
	"encoding/binary"
	"math"
)

// write pulls one block from the game and sends it.
func (p *pulse) write(samples int) error {
	block := p.block[:samples]
	p.audio.pull(block)
	for i, sample := range block {
		binary.LittleEndian.PutUint32(p.out[i*4:], math.Float32bits(sample))
	}
	return p.send(p.channel, p.out[:samples*4])
}
