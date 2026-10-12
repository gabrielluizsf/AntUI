package apk

import (
	"archive/zip"
	"path"
	"strings"
)

// scanPackage reads the archive once and reports what a store check cares
// about: which ABIs have a library, which libraries are stored compressed,
// whether the symbols are there, and the total compressed size.
func scanPackage(z *zip.Reader) (abis map[string]bool, compressed []string, symbols bool, size int64) {
	abis = map[string]bool{}
	for _, e := range z.File {
		size += int64(e.CompressedSize64)
		switch {
		case strings.HasPrefix(e.Name, bundleSymbols+"/"):
			symbols = true
		case strings.HasSuffix(e.Name, ".so"):
			// lib/<abi>/x.so in a package, base/lib/<abi>/x.so in a bundle.
			if dir := path.Dir(e.Name); dir != "." {
				abis[path.Base(dir)] = true
			}
			// A library the app maps out of the package rather than
			// unpacking has to be stored: a deflated entry cannot be mapped,
			// and the app fails to load on the device with nothing here
			// saying why.
			if e.Method != zip.Store {
				compressed = append(compressed, e.Name)
			}
		}
	}
	return abis, compressed, symbols, size
}
