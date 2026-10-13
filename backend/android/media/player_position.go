//go:build android

package media

import "time"

// position is Position without the lock, which the callers below already
// hold.
func (p *Player) position() time.Duration {
	if p.info.HasAudio && p.info.SampleRate > 0 && p.played > 0 {
		frames := p.played / int64(max(p.info.Channels, 1))
		return time.Duration(frames) * time.Second / time.Duration(p.info.SampleRate)
	}
	if !p.playing {
		return p.base
	}
	return p.base + time.Since(p.started)
}
