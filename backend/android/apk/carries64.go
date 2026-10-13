package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

func carries64(abis []sdk.ABI) bool {
	for _, a := range abis {
		if a.Bits() == 64 {
			return true
		}
	}
	return false
}
