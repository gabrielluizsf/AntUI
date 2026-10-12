//go:build android

package ndk

import (
	"fmt"
)

// DecodeAudio starts decoding a sound track into PCM.
func (m *Media) DecodeAudio(track Track) (*Decoder, error) {
	if !track.Audio() {
		return nil, fmt.Errorf("ndk: track %d is %s, not a sound", track.Index, track.Mime)
	}
	d := &Decoder{track: track}
	if err := d.open(m.path, track); err != nil {
		return nil, err
	}
	if err := d.configure(nil); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}
