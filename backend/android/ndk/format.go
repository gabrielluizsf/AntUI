//go:build android

package ndk

// Format is how the pixels in a window's buffer are laid out. The names are
// Android's, and Android names them in **byte order**, not in the order the
// bits sit in a word — so RGBA8888 is R,G,B,A in memory, which on a
// little-endian machine reads back as 0xAABBGGRR. A canvas.Color is
// 0xAARRGGBB, so the two differ by a red/blue swap and every blit pays for
// it. This is worth knowing before debugging a picture that comes out blue.
type Format int32
