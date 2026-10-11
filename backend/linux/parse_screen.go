//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "encoding/binary"

// parseScreen reads the first screen's geometry and root visual, and picks
// the pixel format that goes with the root depth.
func (x *Driver) parseScreen(body []byte, screenAt, formatsAt, formats int) {
	x.root = binary.LittleEndian.Uint32(body[screenAt:])
	x.screenWidth = binary.LittleEndian.Uint16(body[screenAt+20:])
	x.screenHeight = binary.LittleEndian.Uint16(body[screenAt+22:])
	x.screenWidthMM = binary.LittleEndian.Uint16(body[screenAt+24:])
	x.screenHeightMM = binary.LittleEndian.Uint16(body[screenAt+26:])
	rootVisual := binary.LittleEndian.Uint32(body[screenAt+32:])
	depth := body[screenAt+38]
	x.visual = rootVisual
	x.rootDepth = depth

	for i := range formats {
		f := body[formatsAt+i*8:]
		if f[0] == depth {
			x.bitsPerPixel = f[1]
			break
		}
	}

	x.detectPixels(body, screenAt, rootVisual)
}
