package apk

import (
	"fmt"
	"strings"
)

// packageProblems is what the store would say about an already-built
// package: whether there is a 64-bit library, whether the libraries are
// stored rather than compressed, whether a bundle carries its symbols, and
// whether it is inside the size the store allows a device to download.
func packageProblems(abis map[string]bool, compressed []string, symbols, bundle bool, size int64) []Problem {
	var ps []Problem
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
	return ps
}
