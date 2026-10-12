package apk

import (
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

func abiList(abis []sdk.ABI) string {
	out := make([]string, len(abis))
	for i, a := range abis {
		out[i] = string(a)
	}
	return strings.Join(out, ", ")
}
