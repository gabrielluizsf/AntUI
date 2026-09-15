//go:build android

package camera

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/gabrielluizsf/antui/backend/android/display"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
	"github.com/gabrielluizsf/antui/backend/android/permission"
	"github.com/gabrielluizsf/antui/canvas"
)

// Facing is which way a camera points.
type Facing int

// The three the platform distinguishes.
const (
	Front    Facing = 0
	Back     Facing = 1
	External Facing = 2
)

// String names the direction the camera points.
func (f Facing) String() string {
	switch f {
	case Front:
		return "front"
	case Back:
		return "back"
	}
	return "external"
}

// Info describes one camera.
type Info struct {
	// ID is what [Open] takes. It is an opaque string, not an index.
	ID     string
	Facing Facing
	// Orientation is how far the sensor's image is turned from the device's
	// natural orientation, in degrees clockwise. It is 90 on almost every
	// phone, because the sensor is mounted on its side — which is why a
	// photograph taken by an app that ignores this comes out sideways.
	Orientation int
}

// ErrNoPermission is opening a camera without having been granted one.
var ErrNoPermission = errors.New("camera: this app has not been granted CAMERA; " +
	"ask for it with antui/backend/android/permission first")

// List is every camera the app may open.
func List() ([]Info, error) {
	m, err := ndk.Cameras()
	if err != nil {
		return nil, err
	}
	defer m.Close()
	ids, err := m.IDs()
	if err != nil {
		return nil, err
	}
	out := make([]Info, 0, len(ids))
	for _, id := range ids {
		out = append(out, Info{
			ID:          id,
			Facing:      Facing(m.Facing(id)),
			Orientation: m.Orientation(id),
		})
	}
	return out, nil
}

// Default is the first camera pointing a given way.
func Default(f Facing) (Info, bool, error) {
	all, err := List()
	if err != nil {
		return Info{}, false, err
	}
	for _, c := range all {
		if c.Facing == f {
			return c, true, nil
		}
	}
	return Info{}, false, nil
}

// Preview is a running camera, producing frames.
type Preview struct {
	info    Info
	manager *ndk.CameraManager
	device  *ndk.CameraDevice

	mu     sync.Mutex
	canvas *canvas.Canvas
	closed bool
	// asked is the size the image reader was made at, which is the size
	// every frame will be.
	askedW, askedH int

	// upright is how far a frame has to be turned to look right, worked out
	// once from the sensor's mounting and the screen's rotation.
	upright int
}

// Options is how a preview is opened.
type Options struct {
	// Width and Height are what to ask the camera for. They are a request:
	// the camera picks the nearest size it supports, and the frames that
	// come back say what it chose. Zero asks for 1280x720, which every
	// camera has.
	Width, Height int
	// Upright turns each frame so that it matches what the user is looking
	// at, which is what almost every app wants and what nothing does by
	// itself. Turning it off gives the sensor's own image, sideways.
	Upright bool
}

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

// uprightRotation works out how far a frame has to be turned.
//
// Two rotations are involved and they go opposite ways. The sensor is
// mounted turned by Orientation; the screen is turned by whatever the user
// has done with the device. A front camera is also mirrored, which changes
// the sign. This is the formula the platform's own documentation gives, and
// it is worth copying rather than deriving.
func uprightRotation(info Info) int {
	screen := 0
	if r, err := display.Rotation(); err == nil {
		screen = r * 90
	}
	if info.Facing == Front {
		return (info.Orientation + screen) % 360
	}
	return (info.Orientation - screen + 360) % 360
}

// Frame is the newest picture from the camera, or false when none has
// arrived since the last call.
//
// The canvas belongs to the preview and is reused: a 1080p frame is eight
// megabytes, and making a new one thirty times a second would do nothing but
// keep the collector busy. Copy anything that has to outlive the next call.
func (p *Preview) Frame() (*canvas.Canvas, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, false, errors.New("camera: the preview is closed")
	}
	if p.canvas == nil {
		if err := p.makeCanvas(p.askedW, p.askedH); err != nil {
			return nil, false, err
		}
	}
	cv := p.canvas
	// canvas.Color is a uint32 underneath, and the converter writes words.
	words := unsafe.Slice((*uint32)(unsafe.Pointer(&cv.Pixels[0])), len(cv.Pixels))
	w, h, ok, err := p.device.Frame(words, cv.Stride, p.upright)
	if err != nil || !ok {
		return nil, false, err
	}
	if w != cv.Width || h != cv.Height {
		// The camera settled on something other than what was asked for.
		// That frame is lost — it was written into a canvas the wrong shape
		// — and the next one lands in a canvas that fits.
		p.askedW, p.askedH = w, h
		p.canvas = nil
		return nil, false, nil
	}
	return cv, true, nil
}

// makeCanvas makes the canvas frames are converted into.
//
// The size is known in advance and does not have to be discovered: the image
// reader was created at exactly this size, so that is what every frame is.
// A quarter turn swaps it.
func (p *Preview) makeCanvas(w, h int) error {
	if p.upright == 90 || p.upright == 270 {
		w, h = h, w
	}
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		return err
	}
	p.canvas = cv
	return nil
}

// Info is which camera this is.
func (p *Preview) Info() Info { return p.info }

// Rotation is how far each frame is being turned, in degrees clockwise. It
// is 0 unless [Options.Upright] was set.
func (p *Preview) Rotation() int { return p.upright }

// Lost reports whether the camera was taken away — by another app, or by
// being unplugged. A preview that has been lost produces no more frames and
// has to be closed and opened again.
func (p *Preview) Lost() bool { return p.device.Lost() }

// Close stops the camera and gives it back. A camera left open is a camera
// no other app can use, and the platform does not take it back on its own.
func (p *Preview) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	p.device.Close()
	p.manager.Close()
	p.canvas = nil
	return nil
}

// String describes the running camera: which one, and how its frames are
// turned to be the right way up.
func (p *Preview) String() string {
	return fmt.Sprintf("%s camera %s, turned %d°", p.info.Facing, p.info.ID, p.upright)
}
