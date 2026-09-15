package apk

import (
	"archive/zip"
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// What the Play Store will say about a build, said here instead.
//
// Every rule below is one an upload is actually measured against, and each
// of them is cheap to check and expensive to discover: the store answers
// hours later, in an email, about a build that has already been made, signed
// and versioned. A version code cannot be reused, so a rejected upload burns
// one.
//
// The rules are not ours and they move. The two that move on a schedule are
// the target API level, which rises every August, and the ABI requirement,
// which has only ever got stricter. Both live in [antui/backend/android/sdk] as
// constants so there is one place to change them.

// A Problem is one thing the store will act on.
type Problem struct {
	// Fatal is whether the store refuses the upload outright. A problem that
	// is not fatal is one it accepts and complains about, or one that will
	// cost the app users without anyone saying so.
	Fatal bool
	// Text is what is wrong, and Fix is what to do about it.
	Text string
	Fix  string
}

// String is the problem as one line, or two when it says what to do
// about it.
func (p Problem) String() string {
	mark := "warning"
	if p.Fatal {
		mark = "refused"
	}
	if p.Fix == "" {
		return mark + ": " + p.Text
	}
	return mark + ": " + p.Text + "\n         " + p.Fix
}

// Fatal reports whether any of these would stop an upload.
func Fatal(ps []Problem) bool {
	for _, p := range ps {
		if p.Fatal {
			return true
		}
	}
	return false
}

// PlayProblems is what the store will say about this config, answered before
// anything is built rather than after it is uploaded.
//
// It reads the config and nothing else, so it is worth calling early: the
// whole point is to fail in a second instead of after a two-minute build and
// a day of waiting.
func (c Config) PlayProblems() []Problem {
	// A copy, so that asking the question does not change the answer: check
	// fills in the defaults, and the defaults are most of what is being
	// judged.
	cfg := c
	if err := cfg.check(); err != nil {
		return []Problem{{Fatal: true, Text: err.Error()}}
	}
	var ps []Problem

	if !carries64(cfg.ABIs) {
		ps = append(ps, Problem{true,
			fmt.Sprintf("no 64-bit ABI (%s)", abiList(cfg.ABIs)),
			"the store has refused 32-bit-only uploads since August 2019; " +
				"build for " + string(sdk.Arm64)})
	}
	if !slices.Contains(cfg.ABIs, sdk.Arm64) {
		ps = append(ps, Problem{false,
			"no arm64-v8a",
			"almost every phone sold since 2016 is arm64; without it the app " +
				"is not offered to them at all"})
	}
	if cfg.TargetSDK < sdk.PlayTarget {
		ps = append(ps, Problem{true,
			fmt.Sprintf("targets API %d, and the store requires %d",
				cfg.TargetSDK, sdk.PlayTarget),
			"the floor rises every August and applies to updates as well as " +
				"new apps"})
	}
	if cfg.Debuggable {
		ps = append(ps, Problem{true,
			"the app is debuggable",
			"android:debuggable lets anything on the device read the app's " +
				"data and attach to it; drop -debug"})
	}
	if cfg.Keystore.Debug() {
		ps = append(ps, Problem{true,
			"signed with the shared debug key",
			"every Android developer has that key and its password; make one " +
				`of your own with "antuiapk keygen"`})
	}
	if strings.HasPrefix(cfg.Package, "com.example.") || cfg.Package == "com.example" {
		ps = append(ps, Problem{true,
			fmt.Sprintf("the application id is %q", cfg.Package),
			"the store reserves com.example; and the id can never be changed " +
				"after the first release"})
	}
	if cfg.Icon == "" {
		ps = append(ps, Problem{false,
			"no icon",
			"the app appears under the platform's blank one, and the listing " +
				"needs a 512x512 one separately"})
	}
	if cfg.Strip != nil && !*cfg.Strip {
		ps = append(ps, Problem{false,
			"the library ships with its symbols in it",
			"it is about a third larger than it needs to be, and the store " +
				"asks for the symbols as a separate file anyway"})
	}
	if cfg.VersionCode <= 0 {
		ps = append(ps, Problem{true,
			"no version code",
			"it must be a positive number and must rise with every upload; a " +
				"code that has been used once can never be used again"})
	}
	return ps
}

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

func abiList(abis []sdk.ABI) string {
	out := make([]string, len(abis))
	for i, a := range abis {
		out[i] = string(a)
	}
	return strings.Join(out, ", ")
}

func has64(abis map[string]bool) bool {
	for name := range abis {
		if abi, err := sdk.ParseABI(name); err == nil && abi.Bits() == 64 {
			return true
		}
	}
	return false
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func plural(n int) string {
	if n == 1 {
		return "y is"
	}
	return "ies are"
}
