package sdk

import (
	"strconv"
	"strings"
)

// parsePlatform reads "android-30", "android-36.1" and "android-34-ext11".
func parsePlatform(name string) (Platform, bool) {
	rest, ok := strings.CutPrefix(name, "android-")
	if !ok {
		return Platform{}, false
	}
	p := Platform{Name: name}
	if base, ext, found := strings.Cut(rest, "-ext"); found {
		n, err := strconv.Atoi(ext)
		if err != nil {
			return Platform{}, false
		}
		p.Ext, rest = n, base
	}
	major, minor, _ := strings.Cut(rest, ".")
	n, err := strconv.Atoi(major)
	if err != nil {
		// A preview platform is named by a letter — "android-VanillaIceCream"
		// — and has no API level to compare, so it is skipped.
		return Platform{}, false
	}
	p.API = n
	if minor != "" {
		p.Minor, _ = strconv.Atoi(minor)
	}
	return p, true
}
