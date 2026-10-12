//go:build android

package media

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
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
	if err := addVideo(p, m); err != nil {
		return nil, err
	}
	if err := addAudio(p, m); err != nil {
		return nil, err
	}
	if p.video == nil && p.audio == nil {
		return nil, errors.New("media: the file has nothing this device can play")
	}
	if err := prepareCanvases(p); err != nil {
		return nil, err
	}
	return p, nil
}
