//go:build android

package app

// nextRequestCode is a code nothing else is using. It stays under 32768
// because startActivityForResult only carries the low bits of it on some
// platform versions, and zero is skipped because parts of the platform use
// that to mean "no request".
func nextRequestCode() int32 {
	for {
		if c := lastCode.Add(1) & 0x7FFF; c != 0 {
			return c
		}
	}
}
