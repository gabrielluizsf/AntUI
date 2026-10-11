//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/canvas"

// SetIcon gives the window an icon for the task bar, by writing _NET_WM_ICON,
// which is what every window manager from this century reads.
//
// The property is one long list of 32-bit numbers: for each size, its width,
// its height, and then its pixels as ARGB, a row at a time. Several sizes go
// one after another in the same property and the window manager picks the one
// it wants — which is why this takes a list rather than a picture.
//
// The values are CARDINAL, and a CARDINAL is 32 bits on the wire whatever the
// machine is, so a pixel is one number and not two.
func (x *Driver) SetIcon(images []*canvas.Canvas) bool {
	if !x.alive || len(images) == 0 {
		return false
	}
	if x.atomNetWMIcon == 0 {
		x.atomNetWMIcon = x.internAtom("_NET_WM_ICON")
		if x.atomNetWMIcon == 0 {
			return false
		}
	}

	count := 0
	for _, image := range images {
		count += 2 + image.Width*image.Height
	}
	// A property is sent in one request, so it has to fit in one. The biggest
	// sizes go in first — they are what a task switcher wants — and anything
	// past what will fit is left out rather than making a request the server
	// refuses whole.
	most := int(x.maxRequestBytes/4) - 8
	if most <= 0 {
		most = 16384
	}

	values := iconValues(images, min(count, most), most)
	if len(values) == 0 {
		return false
	}

	x.changeProperty(x.atomNetWMIcon, atomCardinal, 32,
		le32(values...), uint32(len(values)))
	return x.alive
}
