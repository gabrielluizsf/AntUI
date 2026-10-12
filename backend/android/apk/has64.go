package apk

import (
	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

func has64(abis map[string]bool) bool {
	for name := range abis {
		if abi, err := sdk.ParseABI(name); err == nil && abi.Bits() == 64 {
			return true
		}
	}
	return false
}
