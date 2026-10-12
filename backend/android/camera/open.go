//go:build android

package camera

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/backend/android/permission"
)

// Open starts a camera.
//
// It needs the CAMERA permission, and says so rather than producing a black
// picture without it. The camera is a shared resource: another app may take
// it away at any moment, which is what [Preview.Lost] reports.
func Open(info Info, opt Options) (*Preview, error) {
	held, err := permission.Held(permission.Camera)
	if err != nil {
		return nil, err
	}
	if !held {
		return nil, ErrNoPermission
	}
	if opt.Width <= 0 || opt.Height <= 0 {
		opt.Width, opt.Height = 1280, 720
	}

	m, err := ndk.Cameras()
	if err != nil {
		return nil, err
	}
	d, err := m.Open(info.ID)
	if err != nil {
		m.Close()
		return nil, err
	}
	if err := d.StartPreview(opt.Width, opt.Height); err != nil {
		d.Close()
		m.Close()
		return nil, err
	}

	p := &Preview{info: info, manager: m, device: d, askedW: opt.Width, askedH: opt.Height}
	if opt.Upright {
		p.upright = uprightRotation(info)
	}
	return p, nil
}
