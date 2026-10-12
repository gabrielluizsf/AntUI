package apk

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// WriteSymbols collects the symbol files for a build into the archive Play
// wants: one entry per ABI, named for the library.
//
// It returns false when there is nothing to write, which is what happens
// when the build was told not to strip — the symbols are then in the shipped
// library itself and a separate copy would be the same file twice.
func WriteSymbols(libs map[sdk.ABI]string, lib, out string) (bool, error) {
	syms := make(map[sdk.ABI]string, len(libs))
	for abi, p := range libs {
		if _, err := os.Stat(symbolFile(p)); err == nil {
			syms[abi] = symbolFile(p)
		}
	}
	if len(syms) == 0 {
		return false, nil
	}

	f, err := os.Create(out)
	if err != nil {
		return false, fmt.Errorf("apk: %w", err)
	}
	defer f.Close()
	z := zip.NewWriter(f)
	for _, abi := range sorted(syms) {
		if err := addSymbol(z, string(abi), lib, syms[abi]); err != nil {
			return false, err
		}
	}
	if err := z.Close(); err != nil {
		return false, fmt.Errorf("apk: %w", err)
	}
	return true, nil
}
