package canvas

// A Color is 0xAARRGGBB. A display surface almost never is, and this is the
// difference for the one format that has no packer already — [To565] is the
// other, and is used rather than written again here, because two functions
// that pack the same format slightly differently is a bug waiting to be
// found in a screenshot.

// packRGBA writes a colour the way a surface named RGBA_8888 or RGBX_8888
// wants it.
//
// The name is in **byte order**: R, G, B, A in memory. On a little-endian
// machine — which is every Android device — that word reads 0xAABBGGRR, so
// red and blue trade places and alpha and green stay where they are. On a
// surface that ignores alpha the top byte is written anyway, because not
// writing it costs a mask and saves nothing.
func packRGBA(c Color) uint32 {
	v := uint32(c)
	return v&0xFF00FF00 | (v&0x00FF0000)>>16 | (v&0x000000FF)<<16
}