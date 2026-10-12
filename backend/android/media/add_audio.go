//go:build android

package media

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// addAudio prepares the audio track, if the file has one.
func addAudio(p *Player, m *ndk.Media) error {
	track, ok := m.Track(false)
	if !ok {
		return nil
	}
	d, err := m.DecodeAudio(track)
	if err != nil {
		if p.video != nil {
			p.video.Close()
		}
		return err
	}
	p.audio = d
	p.info.HasAudio = true
	p.info.SampleRate = track.SampleRate
	p.info.Channels = track.Channels
	if p.info.Duration == 0 {
		p.info.Duration = track.Duration
	}
	// About a quarter of a second of slack between the decoder and the
	// device, which is enough to ride out a slow frame and short enough
	// that a seek is not heard as a delay.
	p.pcm = make([]float32, max(track.SampleRate*track.Channels/4, 4096))
	return nil
}
