package apk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// The one part of a build that does the same work twice.
//
// Almost nothing here is worth caching. The Go compiler has its own cache
// and is already incremental; linking a shared library has to happen every
// time; signing and aligning are a fraction of a second. What is left is the
// icon: seventeen reductions of one 512-pixel picture, and then aapt2
// compiling them, producing byte-for-byte the same archive on every build of
// an app whose icon has not changed. Measured on the jnitest example that is
// about 0.4 s of a 2.1 s rebuild, which is the difference between an
// edit-and-run loop that feels immediate and one that does not.
//
// The key is the content of the icon, so there is nothing to invalidate: a
// changed icon is a different key, and an old entry is simply never asked
// for again.

// cacheVersion changes when what is stored changes shape. An entry made by
// an older version is under a different key and is ignored rather than
// misread.
const cacheVersion = "1"

// compiledIcons returns an aapt2-compiled archive of the launcher icons,
// from the cache if it is there and freshly made if it is not.
//
// A cache that cannot be reached is not an error: it does the work and says
// nothing. Being unable to write to a cache directory is a reason to be
// slower and not a reason to fail a build.
func compiledIcons(tc *sdk.Toolchain, work, icon string, bg canvas.Color) (string, error) {
	key, err := iconKey(tc, icon, bg)
	cached := ""
	if err == nil {
		if dir, err := cacheDir(); err == nil {
			cached = filepath.Join(dir, key+".zip")
			if _, err := os.Stat(cached); err == nil {
				return cached, nil
			}
		}
	}

	res := filepath.Join(work, "res")
	if err := WriteIcons(res, icon, bg); err != nil {
		return "", err
	}
	out := filepath.Join(work, "resources.zip")
	if _, err := run(tc.BuildTools.Aapt2(), "compile", "--dir", res, "-o", out); err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	if cached != "" {
		// Written beside and renamed, so that two builds at once cannot
		// leave a half-written archive for a third to read.
		save(out, cached)
	}
	return out, nil
}

// iconKey is what the archive depends on: the picture itself, the colour
// behind it, and the tool that compiled it.
func iconKey(tc *sdk.Toolchain, icon string, bg canvas.Color) (string, error) {
	body, err := os.ReadFile(icon)
	if err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	h := sha256.New()
	fmt.Fprintf(h, "%s\n%s\n%s\n", cacheVersion, bg, tc.BuildTools.Version)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))[:32], nil
}

func cacheDir() (string, error) {
	root, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, "antuiapk", "res")
	return dir, os.MkdirAll(dir, 0o755)
}

// save copies a file into the cache, atomically. Every failure is ignored:
// the caller has what it needs and the cache is only ever an optimisation.
func save(from, to string) {
	body, err := os.ReadFile(from)
	if err != nil {
		return
	}
	tmp := to + ".tmp"
	if os.WriteFile(tmp, body, 0o644) != nil {
		return
	}
	if os.Rename(tmp, to) != nil {
		os.Remove(tmp)
	}
}
