package sdk

import (
	"cmp"
	"strconv"
	"strings"
)

// Version is a dotted version out of the Android SDK — "37.0.0" for the build
// tools, "30.0.16138531" for the NDK — with the suffix a preview carries.
//
// The parts are kept as they were written rather than folded into three named
// fields, because the SDK is not consistent about how many there are.
type Version struct {
	Parts []int
	// Pre is what followed a hyphen: "rc2", "beta1". Empty on a release.
	Pre string
}

// ParseVersion reads a version. It never fails: anything it cannot make sense
// of contributes nothing, so a directory named by something other than a
// version sorts below every real one instead of stopping the search.
func ParseVersion(s string) Version {
	var v Version
	if i := strings.IndexByte(s, '-'); i >= 0 {
		v.Pre, s = s[i+1:], s[:i]
	}
	for _, part := range strings.Split(s, ".") {
		n, err := strconv.Atoi(part)
		if err != nil {
			break
		}
		v.Parts = append(v.Parts, n)
	}
	return v
}

// String writes the version back out the way it came in.
func (v Version) String() string {
	if len(v.Parts) == 0 {
		return ""
	}
	var b strings.Builder
	for i, n := range v.Parts {
		if i > 0 {
			b.WriteByte('.')
		}
		b.WriteString(strconv.Itoa(n))
	}
	if v.Pre != "" {
		b.WriteByte('-')
		b.WriteString(v.Pre)
	}
	return b.String()
}

// Major is the first part, or 0 when there is none.
func (v Version) Major() int {
	if len(v.Parts) == 0 {
		return 0
	}
	return v.Parts[0]
}

// Release reports whether this is a release rather than a preview. Previews
// are never chosen on their own; see [pick].
func (v Version) Release() bool { return v.Pre == "" }

// Compare orders two versions: by their numbers, then by whether they are a
// release, so 37.0.0 sorts above 37.0.0-rc2. A version with fewer parts is
// the smaller one when the parts they share are equal, so 37.0 is below
// 37.0.1.
func (v Version) Compare(w Version) int {
	for i := range max(len(v.Parts), len(w.Parts)) {
		var a, b int
		if i < len(v.Parts) {
			a = v.Parts[i]
		}
		if i < len(w.Parts) {
			b = w.Parts[i]
		}
		if c := cmp.Compare(a, b); c != 0 {
			return c
		}
	}
	switch {
	case v.Pre == w.Pre:
		return 0
	case v.Pre == "":
		return 1
	case w.Pre == "":
		return -1
	}
	return cmp.Compare(v.Pre, w.Pre)
}

// AtLeast reports whether v is w or newer.
func (v Version) AtLeast(w Version) bool { return v.Compare(w) >= 0 }

// pick chooses the newest of what was found, preferring a release over any
// preview however new the preview is — an rc of the build tools has shipped
// broken often enough that it is not worth taking by accident.
func pick[T any](found []T, version func(T) Version) (T, bool) {
	var best T
	var bestV Version
	ok := false
	for _, f := range found {
		v := version(f)
		switch {
		case !ok:
		case bestV.Release() && !v.Release():
			continue
		case !bestV.Release() && v.Release():
		case v.Compare(bestV) <= 0:
			continue
		}
		best, bestV, ok = f, v, true
	}
	return best, ok
}
