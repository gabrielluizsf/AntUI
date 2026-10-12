//go:build android

package media

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/canvas"
)

// Open reads a file and gets both its tracks ready.
//
// The path is a file on this device. An address from
// [antui/backend/android/picker] is not one, and opening it needs the content
// resolver
func Open(path string) (*Player, error) {
	m, err := ndk.OpenMedia(path)
	if err != nil {
		return nil, err
	}
	p := &Player{}

	if track, ok := m.Track(true); ok {
		d, err := m.DecodeVideo(track, track.Rotation)
		if err != nil {
			return nil, err
		}
		p.video = d
		p.info.HasVideo = true
		p.info.Rotation = track.Rotation
		p.info.Width, p.info.Height = track.Width, track.Height
		if track.Rotation == 90 || track.Rotation == 270 {
			p.info.Width, p.info.Height = track.Height, track.Width
		}
		p.info.Duration = track.Duration
	}
	if track, ok := m.Track(false); ok {
		d, err := m.DecodeAudio(track)
		if err != nil {
			if p.video != nil {
				p.video.Close()
			}
			return nil, err
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
	}
	if p.video == nil && p.audio == nil {
		return nil, errors.New("media: the file has nothing this device can play")
	}

	if p.video != nil {
		cv, err := canvas.NewCanvas(p.info.Width, p.info.Height)
		if err != nil {
			p.Close()
			return nil, err
		}
		p.cur = cv
		if cv, err = canvas.NewCanvas(p.info.Width, p.info.Height); err == nil {
			p.next = cv
		}
	}
	return p, nil
}
