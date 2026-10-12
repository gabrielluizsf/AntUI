//go:build android

package media

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// addVideo prepares the video track, if the file has one.
func addVideo(p *Player, m *ndk.Media) error {
	track, ok := m.Track(true)
	if !ok {
		return nil
	}
	d, err := m.DecodeVideo(track, track.Rotation)
	if err != nil {
		return err
	}
	p.video = d
	p.info.HasVideo = true
	p.info.Rotation = track.Rotation
	p.info.Width, p.info.Height = track.Width, track.Height
	if track.Rotation == 90 || track.Rotation == 270 {
		p.info.Width, p.info.Height = track.Height, track.Width
	}
	p.info.Duration = track.Duration
	return nil
}
