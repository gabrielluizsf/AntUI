package apk

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Symbols are how a crash on someone else's phone stays readable.
//
// A native crash arrives from the store as a list of addresses. Turning
// those back into function names needs a copy of the library with its symbol
// table still in it, and the copy has to be the *same build* — a rebuild from
// the same source lands the code somewhere else and symbolises to nonsense.
// So the moment to keep it is the moment it is made, and there is no second
// chance: once the stripped library has shipped, a crash in it is a list of
// addresses forever.
//
// This is why stripping happens here rather than at the link. Go's -s and -w
// tell the linker never to write the symbols at all, which costs the same as
// stripping afterwards — the shipped library comes out the same size to the
// byte — and throws away the only copy. Building whole and stripping after
// leaves something to keep.
//
// Sizes, for one real arm64 app:
//
//	4.4 MB   as built, with DWARF
//	3.3 MB   the symbol file: DWARF dropped, symbol table kept
//	3.0 MB   what ships, stripped
//	3.0 MB   what -s -w used to ship, with nothing kept
//
// The archive is what Play asks for on the release page, and [Bundle] also
// puts it inside the bundle so that one upload carries both.
const symbolSuffix = ".sym"

// symbolFile is where BuildLibs leaves the symbols for a library: beside it,
// with .sym on the end. It is a sibling rather than a separate directory so
// that a caller who keeps the build directory has them both together.
func symbolFile(lib string) string { return lib + symbolSuffix }

// keepSymbols splits a freshly built library in two: a symbol file that can
// read a crash, and the library that ships.
//
// The symbol file keeps the symbol table and drops the DWARF, which is what
// the store's symbolication uses and a third of the size of keeping
// everything. Function names come back; line numbers do not.
func keepSymbols(tc *sdk.Toolchain, lib string) error {
	objcopy := tc.NDK.Tool("objcopy")
	strip := tc.NDK.Tool("strip")

	if err := copyFile(lib, symbolFile(lib)); err != nil {
		return fmt.Errorf("apk: keeping symbols: %w", err)
	}
	if _, err := run(objcopy, "--strip-debug", symbolFile(lib)); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	if _, err := run(strip, "--strip-all", lib); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}

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

// bundleSymbols is the same set of files under the path a bundle carries
// them at, so that uploading the bundle uploads the symbols with it.
const bundleSymbols = "BUNDLE-METADATA/com.android.tools.build.debugsymbols"

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

// addSymbol puts one symbol file in a zip under dir.
func addSymbol(z *zip.Writer, dir, lib, src string) error {
	name := path.Join(dir, "lib"+lib+".so"+symbolSuffix)
	w, err := z.Create(name)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer in.Close()
	if _, err := io.Copy(w, in); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return nil
}

// sorted puts the ABIs of a map in a fixed order, so that two builds of the
// same source produce the same archive.
func sorted(m map[sdk.ABI]string) []sdk.ABI {
	out := make([]sdk.ABI, 0, len(m))
	for abi := range m {
		out = append(out, abi)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// symbolsPath is where the archive goes when the config does not say: beside
// the package, named after it.
//
//	hello.apk  →  hello-symbols.zip
//	hello.aab  →  hello-symbols.zip
func symbolsPath(cfg *Config) string {
	if cfg.Symbols != "" {
		return cfg.Symbols
	}
	out := cfg.Out
	return strings.TrimSuffix(out, filepath.Ext(out)) + "-symbols.zip"
}

// copyFile is a plain copy, keeping the mode.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	st, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, st.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
