//go:build windows

package windows

import "unsafe"

// blit hands the whole canvas to GDI as a top-down 32-bit DIB. The negative
// height is what "top-down" means to Windows; without it the frame arrives
// upside down.
func (d *Driver) blit(hdc uintptr) {
	cv := d.win.Canvas()
	if cv.Pixels == nil || cv.Width <= 0 || cv.Height <= 0 {
		return
	}
	d.info = bitmapInfoHeader{
		Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		Width:       int32(cv.Width),
		Height:      -int32(cv.Height),
		Planes:      1,
		BitCount:    32,
		Compression: biRGB,
	}
	procStretchDIBits.Call(
		hdc,
		0, 0, uintptr(cv.Width), uintptr(cv.Height),
		0, 0, uintptr(cv.Width), uintptr(cv.Height),
		uintptr(unsafe.Pointer(&cv.Pixels[0])),
		uintptr(unsafe.Pointer(&d.info)),
		dibRGBColor, srcCopy,
	)
}
