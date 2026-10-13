//go:build android

package ndk

import "fmt"

// String names the pixel format the way Android does.
func (f Format) String() string {
	switch f {
	case RGBA8888:
		return "RGBA_8888"
	case RGBX8888:
		return "RGBX_8888"
	case RGB888:
		return "RGB_888"
	case RGB565:
		return "RGB_565"
	}
	return fmt.Sprintf("format(%d)", int32(f))
}
