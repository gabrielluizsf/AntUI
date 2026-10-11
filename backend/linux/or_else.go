//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// orElse is v, or the fallback when v says nothing.
func orElse(v, fallback int) int {
	if v > 0 {
		return v
	}
	return fallback
}
