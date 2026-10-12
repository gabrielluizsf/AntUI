//go:build android

package media

// Sound is what to hand to [antui.OpenAudio]: it fills a block with the
// file's sound and keeps the clock.
//
// It is nil when the file has no sound, and a caller should check rather
// than opening a silent device.
func (p *Player) Sound() func([]float32) {
	if p.audio == nil {
		return nil
	}
	return p.fill
}
