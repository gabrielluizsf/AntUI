//go:build android

package ndk

// The image format a preview asks for. YUV420 is the one every camera
// supports and the only one worth asking for.
const FormatYUV420 = 0x23
