package apk

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// An Android App Bundle is not an APK.
//
// It is what the Play Store has taken instead of one since August 2021, and
// it is a different format rather than a different name: the manifest and
// the resource table inside it are **protocol buffers**, not the binary XML
// and arsc an APK holds, and everything sits under a module directory. What
// the user installs is built from it by the store, one APK per device, with
// only the ABI and the densities that device needs — which is where the
// download saving comes from and why a bundle carries every ABI and an APK
// should not.
//
// aapt2 can produce the proto forms with --proto-format; the rearranging and
// the BundleConfig are done here. bundletool is not needed to *build* one —
// only to take one apart again, which is what [VerifyBundle] is for when it
// is there.

// Bundle builds an .aab, signed and ready to upload.
func Bundle(tc *sdk.Toolchain, cfg Config) (*Result, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	if tc.NDK.Dir == "" {
		return nil, fmt.Errorf("apk: no NDK; a Go app is native code")
	}
	if tc.JDK.Dir == "" {
		return nil, fmt.Errorf("apk: signing a bundle needs a JDK: apksigner " +
			"cannot sign one, and jarsigner is what does")
	}
	if cfg.Out == cfg.Label+".apk" {
		cfg.Out = cfg.Label + ".aab"
	}
	// What the store would say, said now. A bundle is built to be uploaded,
	// so it is worth hearing before the two minutes of building rather than
	// after — but it is said and not enforced: a bundle built to be tried on
	// an emulator wants x86_64, which no phone has and no store accepts, and
	// that is a perfectly good reason to build one.
	for _, p := range cfg.PlayProblems() {
		fmt.Fprintln(os.Stderr, "antuiapk: "+p.String())
	}

	work, err := os.MkdirTemp("", "antuiaab-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)

	libs, err := BuildLibs(tc, &cfg, work)
	if err != nil {
		return nil, err
	}
	proto := filepath.Join(work, "proto.apk")
	if err := LinkProto(tc, &cfg, proto); err != nil {
		return nil, err
	}
	unsigned := filepath.Join(work, "unsigned.aab")
	if err := AssembleBundle(proto, libs, &cfg, unsigned); err != nil {
		return nil, err
	}
	if err := SignBundle(tc, &cfg, unsigned, cfg.Out); err != nil {
		return nil, err
	}

	st, err := os.Stat(cfg.Out)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	res := &Result{APK: cfg.Out, Libs: libs, Size: st.Size()}
	// The bundle already carries the symbols inside it, so a copy beside it
	// would be the same three megabytes twice. One is written only when the
	// caller asked for one by name.
	if cfg.Symbols != "" {
		if wrote, err := WriteSymbols(libs, cfg.Lib, cfg.Symbols); err != nil {
			return nil, err
		} else if wrote {
			res.Symbols = cfg.Symbols
		}
	}
	return res, nil
}

func carries64(abis []sdk.ABI) bool {
	for _, a := range abis {
		if a.Bits() == 64 {
			return true
		}
	}
	return false
}

// LinkProto is [Link] asking aapt2 for the protocol-buffer forms, which is
// what a bundle holds instead of binary XML and an arsc.
func LinkProto(tc *sdk.Toolchain, cfg *Config, out string) error {
	return link(tc, cfg, out, true)
}

// AssembleBundle rearranges what aapt2 produced into the layout a bundle
// has, and adds the libraries, the dex and the configuration.
func AssembleBundle(proto string, libs map[sdk.ABI]string, cfg *Config, out string) error {
	zr, err := zip.OpenReader(proto)
	if err != nil {
		return fmt.Errorf("apk: reading %s: %w", proto, err)
	}
	defer zr.Close()

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)

	// BundleConfig first, at the root and outside every module.
	if err := add(zw, "BundleConfig.pb", bundleConfig()); err != nil {
		return err
	}

	for _, e := range zr.File {
		name := bundlePath(e.Name)
		if name == "" {
			continue
		}
		rc, err := e.Open()
		if err != nil {
			return fmt.Errorf("apk: %s: %w", e.Name, err)
		}
		body, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("apk: %s: %w", e.Name, err)
		}
		if err := add(zw, name, body); err != nil {
			return err
		}
	}

	for abi, lib := range libs {
		body, err := os.ReadFile(lib)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		name := path.Join("base", "lib", string(abi), "lib"+cfg.Lib+".so")
		if err := add(zw, name, body); err != nil {
			return err
		}
	}
	if cfg.Shim != nil && *cfg.Shim {
		if err := add(zw, "base/dex/classes.dex", shim.Dex); err != nil {
			return err
		}
	}
	// The symbols ride along, under the path the store looks for them at.
	// They are metadata and not part of any module, so nothing about them
	// reaches a device: the bundle is where they are kept, not shipped.
	if err := addSymbols(zw, libs, cfg.Lib); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return f.Close()
}

// bundlePath is where something from the proto APK belongs in the bundle, or
// empty for something that does not belong in one at all.
func bundlePath(name string) string {
	switch {
	case name == "AndroidManifest.xml":
		// The one file that moves rather than being prefixed. A bundle keeps
		// the manifest in a directory of its own because a module may have
		// more than one thing describing it.
		return "base/manifest/AndroidManifest.xml"
	case name == "resources.pb":
		return "base/resources.pb"
	case strings.HasPrefix(name, "res/"), strings.HasPrefix(name, "assets/"):
		return "base/" + name
	case strings.HasPrefix(name, "META-INF/"):
		// Whatever signed the intermediate is not what signs this.
		return ""
	case name == "resources.arsc":
		// Only there if --proto-format was not passed, which would be a bug
		// rather than something to carry.
		return ""
	}
	// Anything else an app packed goes under root/, which is where a bundle
	// keeps files that are not resources, assets, dex or libraries.
	return "base/root/" + name
}

func add(zw *zip.Writer, name string, body []byte) error {
	w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
	if err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("apk: %s: %w", name, err)
	}
	return nil
}

// SignBundle signs with jarsigner.
//
// Not apksigner, which refuses a bundle: the signature schemes it writes go
// in an APK's own structures, and a bundle is signed the old way, as a JAR.
// What the user installs is signed later — by the store, with the key it
// holds — so this signature only proves who uploaded it.
func SignBundle(tc *sdk.Toolchain, cfg *Config, in, out string) error {
	if err := cfg.Keystore.Ensure(tc); err != nil {
		return err
	}
	k := cfg.Keystore
	jarsigner := filepath.Join(tc.JDK.Dir, "bin", "jarsigner")
	_, err := run(jarsigner,
		"-keystore", k.Path,
		"-storepass", k.StorePass,
		"-keypass", k.KeyPass,
		"-signedjar", out,
		"-digestalg", "SHA-256",
		"-sigalg", "SHA256withRSA",
		in, k.Alias,
	)
	if err != nil {
		return fmt.Errorf("apk: signing the bundle: %w", err)
	}
	return nil
}

// BundleConfig.pb, encoded by hand.
//
// It is a protocol buffer and there is no generated code here to write one
// with, but the message needed is small enough to spell out: which tool
// built it, and which files the store should leave uncompressed in the APKs
// it generates. Everything else is a default, and a default is what an
// ordinary app wants.
func bundleConfig() []byte {
	var out bytes.Buffer

	// Bundletool { version = 2 }. The field number is **2** and not 1, which
	// is the sort of thing that cannot be reasoned out: a wrong number
	// produces a perfectly well-formed protocol buffer, and bundletool
	// answers "Version must match the format '<major>.<minor>.<revision>',
	// but found ''" — an empty string, because it read a field that was not
	// there. It was settled by having bundletool write one and comparing the
	// bytes, which is the only way to settle it.
	out.Write(protoMessage(1, protoString(2, bundleFormat)))

	// Compression { uncompressed_glob = ... } — field 3, inner field 1,
	// checked the same way. The native libraries have to be stored rather
	// than deflated in the generated APKs, because the manifest says the
	// loader maps them out of the package and a compressed entry cannot be
	// mapped.
	var compression bytes.Buffer
	for _, glob := range []string{"**.so"} {
		compression.Write(protoString(1, glob))
	}
	out.Write(protoMessage(3, compression.Bytes()))
	return out.Bytes()
}

// bundleFormat is the bundle format this writes. It is a bundletool version
// because that is what the field holds; nothing here runs bundletool. This
// one is the version whose own output the encoding above was checked against.
const bundleFormat = "1.18.3"

// The two pieces of protocol buffer encoding this needs. Both fields used
// above are length-delimited, which is wire type 2.
func protoString(field int, s string) []byte {
	return protoBytes(field, []byte(s))
}

func protoMessage(field int, body []byte) []byte {
	return protoBytes(field, body)
}

func protoBytes(field int, body []byte) []byte {
	var out []byte
	out = binary.AppendUvarint(out, uint64(field)<<3|2)
	out = binary.AppendUvarint(out, uint64(len(body)))
	return append(out, body...)
}

// ErrNoBundletool is a machine without the jar that takes a bundle apart.
//
// Nothing needs it to build one — this package writes the format itself —
// and it is not part of the SDK. It is one jar from Google's releases, and
// it is the only way to see what the store will actually make of an upload.
var ErrNoBundletool = fmt.Errorf("apk: bundletool is not on this machine; " +
	"get the jar from github.com/google/bundletool/releases and put it at " +
	"$ANDROID_HOME/bundletool.jar, or point $BUNDLETOOL at it")

// VerifyBundle takes a bundle apart the way the store will, and puts the
// result on a device.
//
// This is the only check that means anything about a bundle. Everything
// before it says the file is shaped right; this says the store can turn it
// into an app that installs and runs — which is a different question, and
// the one that matters at three in the morning before a release.
//
// serial names a device, or is empty for the only one attached.
func VerifyBundle(tc *sdk.Toolchain, cfg *Config, bundle, serial string) error {
	if err := cfg.check(); err != nil {
		return err
	}
	if tc.Bundletool == "" {
		return ErrNoBundletool
	}
	if tc.JDK.Dir == "" {
		return fmt.Errorf("apk: bundletool is a jar and needs a JDK to run")
	}

	work, err := os.MkdirTemp("", "antuiverify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	apks := filepath.Join(work, "out.apks")

	k := cfg.Keystore
	args := []string{"-jar", tc.Bundletool, "build-apks",
		"--bundle=" + bundle,
		"--output=" + apks,
		"--ks=" + k.Path,
		"--ks-key-alias=" + k.Alias,
		"--ks-pass=pass:" + k.StorePass,
		"--key-pass=pass:" + k.KeyPass,
		// Only the splits this device needs, which is what the store does
		// and is the whole point of the format.
		"--connected-device",
	}
	if serial != "" {
		args = append(args, "--device-id="+serial)
	}
	java := filepath.Join(tc.JDK.Dir, "bin", "java")
	if _, err := run(java, args...); err != nil {
		return fmt.Errorf("apk: bundletool could not build APKs from the bundle: %w", err)
	}

	install := []string{"-jar", tc.Bundletool, "install-apks", "--apks=" + apks}
	if serial != "" {
		install = append(install, "--device-id="+serial)
	}
	if _, err := run(java, install...); err != nil {
		return fmt.Errorf("apk: the APKs the bundle produced would not install: %w", err)
	}
	return nil
}
