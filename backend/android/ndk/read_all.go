//go:build android

package ndk

import (
	"io"
)

// ReadAll is the whole asset as Go memory, mapped where it can be and
// inflated where it cannot.
func (a *Asset) ReadAll() ([]byte, error) {
	if mapped, ok := a.Bytes(); ok {
		out := make([]byte, len(mapped))
		copy(out, mapped)
		return out, nil
	}
	if a.size <= 0 {
		return nil, nil
	}
	out := make([]byte, a.size)
	read := 0
	for read < len(out) {
		n, err := a.Read(out[read:])
		if n > 0 {
			read += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return out[:read], nil
}
