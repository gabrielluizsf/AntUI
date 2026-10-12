package apk

import (
	"archive/zip"
	"fmt"
	"path"
	"strings"
)

// PlayProblemsIn is the same question asked of a package that has already
// been built — an .apk or an .aab — because a config is what was meant and
// the file is what will be uploaded.
//
// It reads the archive and does not run anything, so it works on a package
// built somewhere else.
func PlayProblemsIn(pkg string) ([]Problem, error) {
	z, err := zip.OpenReader(pkg)
	if err != nil {
		return nil, fmt.Errorf("apk: reading %s: %w", pkg, err)
	}
	defer z.Close()

	var ps []Problem
	abis := map[string]bool{}
	var compressed []string
	symbols := false
	var size int64

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

	bundle := strings.HasSuffix(pkg, ".aab")
	if len(abis) == 0 {
		ps = append(ps, Problem{true, "no native library in the package",
			"a Go app is native code; the package has nothing to run"})
	} else if !has64(abis) {
		ps = append(ps, Problem{true,
			"no 64-bit library (" + strings.Join(keys(abis), ", ") + ")",
			"the store has refused 32-bit-only uploads since August 2019"})
	}
	if len(compressed) > 0 && !bundle {
		ps = append(ps, Problem{true,
			fmt.Sprintf("%d native librar%s stored compressed", len(compressed),
				plural(len(compressed))),
			"the manifest says the loader maps them out of the package, and a " +
				"compressed entry cannot be mapped"})
	}
	if bundle && !symbols {
		ps = append(ps, Problem{false, "the bundle carries no symbols",
			"a native crash from a user's phone arrives as a list of " +
				"addresses and there is no way to read it later"})
	}
	if limit := int64(150 << 20); bundle && size > limit {
		ps = append(ps, Problem{true,
			fmt.Sprintf("the bundle is %.0f MB", float64(size)/(1<<20)),
			"the store's limit for what a device downloads is 150 MB; the " +
				"rest has to be asset packs"})
	}
	return ps, nil
}
