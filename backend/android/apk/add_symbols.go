package apk

import (
	"archive/zip"
	"os"
	"path"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// addSymbols writes the symbol files into a bundle being assembled.
func addSymbols(z *zip.Writer, libs map[sdk.ABI]string, lib string) error {
	for _, abi := range sorted(libs) {
		sym := symbolFile(libs[abi])
		if _, err := os.Stat(sym); err != nil {
			continue
		}
		if err := addSymbol(z, path.Join(bundleSymbols, string(abi)), lib, sym); err != nil {
			return err
		}
	}
	return nil
}
