// Package apk builds an Android package out of a Go program.
//
// It runs on the machine doing the building, not on the phone: no cgo, no
// build tag, no device. What it needs is an SDK — see
// [antui/backend/android/sdk] — and a Go package that registers an app with
// [antui/backend/android/app].
//
// The pipeline is four steps, and each one is a separate function so that a
// caller who wants only part of it can have it:
//
//	BuildLibs   compile the Go package once per ABI, with the NDK's clang
//	Link        turn a manifest into a base APK, with aapt2
//	Assemble    put the libraries in the APK, uncompressed
//	Sign        align it and sign it
//
// [Build] is all four.
package apk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// Config is everything about the app being built. Only Package, Label and
// Dir have no sensible default.
type Config struct {
	// Package is the application id — "com.example.hello". It is what
	// Android identifies the app by forever: two apps with the same one
	// cannot be installed side by side, and it cannot be changed after a
	// release without the store treating it as a different app.
	Package string `json:"package,omitempty"`
	// Label is the name under the icon.
	Label string `json:"label,omitempty"`

	// Dir is the Go package to build. It must be a main package that
	// registers an app.
	Dir string `json:"dir,omitempty"`
	// Lib is the name of the library inside the APK, without the "lib"
	// prefix or the ".so". It ends up in the manifest as
	// android.app.lib_name. Defaults to "antui".
	Lib string `json:"lib,omitempty"`

	// VersionCode is what the store orders releases by, and must go up with
	// every upload. VersionName is what the user is shown and can be
	// anything.
	VersionCode int    `json:"versionCode,omitempty"`
	VersionName string `json:"versionName,omitempty"`

	// MinSDK is the oldest Android that may install this. TargetSDK is the
	// newest whose behaviour the app has been written for, and the store has
	// a floor for it — see [sdk.PlayTarget].
	MinSDK    int `json:"minSdk,omitempty"`
	TargetSDK int `json:"targetSdk,omitempty"`

	// ABIs to build. Defaults to [sdk.Ship], which is what a store upload
	// wants; add sdk.X64 to run on the emulator.
	ABIs []sdk.ABI `json:"abis,omitempty"`

	// Permissions are uses-permission entries, by their full name:
	// "android.permission.INTERNET".
	Permissions []string `json:"permissions,omitempty"`
	// Features are uses-feature entries. A feature listed here without
	// android:required="false" stops the app being offered to devices that
	// lack it, which is usually not what is wanted.
	Features []Feature `json:"features,omitempty"`

	// Icon is a PNG or a JPEG to make the launcher icon from — one square
	// picture, 512 across for preference, and every size and shape the
	// platform wants is generated from it. Empty leaves the app with the
	// platform's own blank icon.
	Icon string `json:"icon,omitempty"`
	// IconBackground is what shows behind the icon on a device that draws
	// adaptive ones, which is every device since Android 8: the picture is
	// the foreground layer and this is the layer under it. Zero means white.
	IconBackground canvas.Color `json:"iconBackground,omitempty"`

	// Assets is a directory whose contents go into the package, readable
	// through [antui/backend/android/assets]. Empty means none.
	//
	// They are not files on the device: they stay inside the APK, which is
	// never unpacked, so nothing on disk corresponds to one.
	Assets string `json:"assets,omitempty"`
	// Store lists the file extensions to leave uncompressed in the package,
	// without the dot. The default is the formats that are already
	// compressed — squeezing a PNG again gains nothing, costs the processor
	// on every read, and stops the file being memory-mapped, which is the
	// difference between opening a large asset instantly and copying it.
	//
	// Set it to a slice with one empty string to store everything.
	Store []string `json:"store,omitempty"`

	// Queries are the intent actions this app needs to be able to *see*
	// other apps for.
	//
	// Since Android 11 an app cannot see what else is installed. Not "cannot
	// ask" — cannot see: a query returns nothing, a service is not there, a
	// speech engine that is plainly installed does not exist. Nothing fails
	// and nothing says why. An app that binds to another app's service has
	// to name the action here.
	//
	//	Queries: []string{"android.intent.action.TTS_SERVICE"}
	Queries []string `json:"queries,omitempty"`

	// Links are the addresses this app opens: a web address it claims, or a
	// scheme of its own. Each becomes an intent filter, and an address that
	// matches one starts the app with it.
	Links []DeepLink `json:"links,omitempty"`

	// Orientation locks the screen. Empty leaves it to the user's device.
	Orientation string `json:"orientation,omitempty"`

	// Shim puts the Java shim in the package, and points the manifest at it.
	// It is what makes a permission request, an activity result and an
	// arriving intent reach the app at all — Android delivers those three by
	// calling a method on the activity, and native code cannot override one.
	//
	// On by default. It costs about 1.5 KB, and the failure mode of leaving
	// it out is not an error but silence: the request goes out and the
	// answer never comes.
	Shim *bool `json:"shim,omitempty"`

	// Backup says whether the system may copy the app's data to the user's
	// cloud account and put it back on their next device.
	//
	// nil leaves the platform's own default, which is on. Turning it off is
	// for an app holding something that must not travel — a key, a session,
	// anything the user would not expect to appear on a phone they have just
	// bought. Opting out one *file* rather than all of them needs an XML
	// resource, which is Phase 11's work.
	Backup *bool `json:"backup,omitempty"`

	// Debuggable lets a debugger attach and makes the app installable over
	// an existing release build. The store refuses an upload with it set.
	Debuggable bool `json:"debuggable,omitempty"`
	// Strip drops the symbol table and DWARF from the shipped libraries,
	// keeping a symbol file beside the package so a crash from a user's
	// phone can still be read. On by default. Turning it off leaves the
	// symbols in the library itself and writes no separate file; either way
	// nothing is thrown away. See [Symbols].
	Strip *bool `json:"strip,omitempty"`
	// Symbols is where the symbol archive goes. Empty puts it beside the
	// package: hello.apk gives hello-symbols.zip. It is what Play asks for
	// on the release page, and [Bundle] also carries a copy inside the
	// bundle so that one upload does both.
	Symbols string `json:"symbols,omitempty"`

	// Out is the APK to write. Defaults to <Label or last path element>.apk
	// in the working directory.
	Out string `json:"out,omitempty"`

	// Keystore signs the result. The zero value means the debug keystore,
	// created on first use — fine for installing, never for the store.
	Keystore Keystore `json:"keystore,omitzero"`

	// Tags and Ldflags are passed through to the Go build.
	Tags    []string `json:"tags,omitempty"`
	Ldflags string   `json:"ldflags,omitempty"`
}

// DefaultStore is what [Config.Store] leaves uncompressed: everything that
// is already compressed, so that the package does not compress it twice.
var DefaultStore = []string{
	"png", "jpg", "jpeg", "webp", "gif",
	"ogg", "mp3", "m4a", "opus", "wav",
	"mp4", "m4v", "webm", "mkv",
	"zip", "gz", "br", "woff", "woff2",
}

// Feature is a uses-feature entry.
type Feature struct {
	Name     string `json:"name"`
	Required bool   `json:"required,omitempty"`
}

// DeepLink is an address this app opens.
//
//	DeepLink{Scheme: "myapp"}                  myapp://anything
//	DeepLink{Scheme: "https", Host: "ex.com"}  every page on ex.com
//	DeepLink{Scheme: "https", Host: "ex.com", Path: "/open"}
type DeepLink struct {
	// Scheme is "https", or a scheme of the app's own. A custom scheme is
	// simpler and is claimed by whoever asks — any other app may claim the
	// same one, and the user is asked which. A web address cannot be stolen
	// that way, but has to be proved.
	Scheme string `json:"scheme"`
	// Host is the domain, for a web address.
	Host string `json:"host,omitempty"`
	// Path is a prefix, and empty means the whole host.
	Path string `json:"path,omitempty"`
	// Verify makes it an App Link: Android checks
	// https://<host>/.well-known/assetlinks.json for this app's signature
	// and, if it finds it, opens the address here **without asking the
	// user**. Without that file the filter still works and the user is asked.
	Verify bool `json:"verify,omitempty"`
}

// valid reports what is wrong with a link, if anything.
func (l DeepLink) valid() error {
	if l.Scheme == "" {
		return fmt.Errorf("apk: a link needs a scheme")
	}
	if l.Host == "" && (l.Scheme == "http" || l.Scheme == "https") {
		return fmt.Errorf("apk: a %s link needs a host, or it claims the whole web", l.Scheme)
	}
	if l.Verify && l.Host == "" {
		return fmt.Errorf("apk: only a link with a host can be verified")
	}
	return nil
}

// The defaults, applied by [Config.check].
const (
	defaultLib     = "antui"
	defaultVersion = "1.0"
)

// check fills in the defaults and reports what is still missing. It is
// called by every entry point, so a caller who skips Build still gets it.
func (c *Config) check() error {
	if c.Package == "" {
		return fmt.Errorf("apk: no package name; it looks like \"com.example.app\"")
	}
	if err := validPackage(c.Package); err != nil {
		return err
	}
	if c.Dir == "" {
		return fmt.Errorf("apk: no Go package to build")
	}
	abs, err := filepath.Abs(c.Dir)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	c.Dir = abs
	if st, err := os.Stat(c.Dir); err != nil || !st.IsDir() {
		return fmt.Errorf("apk: %s is not a directory", c.Dir)
	}

	if c.Label == "" {
		c.Label = filepath.Base(c.Dir)
	}
	if c.Lib == "" {
		c.Lib = defaultLib
	}
	if c.VersionCode == 0 {
		c.VersionCode = 1
	}
	if c.VersionName == "" {
		c.VersionName = defaultVersion
	}
	if c.MinSDK == 0 {
		c.MinSDK = sdk.MinSDK
	}
	if c.TargetSDK == 0 {
		c.TargetSDK = sdk.PlayTarget
	}
	if c.MinSDK > c.TargetSDK {
		return fmt.Errorf("apk: minSdk %d is above targetSdk %d", c.MinSDK, c.TargetSDK)
	}
	if len(c.ABIs) == 0 {
		c.ABIs = sdk.Ship
	}
	for _, a := range c.ABIs {
		if !a.Valid() {
			return fmt.Errorf("apk: %q is not an Android ABI", a)
		}
	}
	for _, l := range c.Links {
		if err := l.valid(); err != nil {
			return err
		}
	}
	if c.Strip == nil {
		on := true
		c.Strip = &on
	}
	if c.Shim == nil {
		on := true
		c.Shim = &on
	}
	if c.Store == nil {
		c.Store = DefaultStore
	}
	if c.Icon != "" {
		abs, err := filepath.Abs(c.Icon)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		if _, err := os.Stat(abs); err != nil {
			return fmt.Errorf("apk: there is no icon at %s", c.Icon)
		}
		c.Icon = abs
	}
	if c.IconBackground == 0 {
		c.IconBackground = canvas.White
	}
	if c.Assets != "" {
		abs, err := filepath.Abs(c.Assets)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		if st, err := os.Stat(abs); err != nil || !st.IsDir() {
			return fmt.Errorf("apk: %s is not a directory of assets", c.Assets)
		}
		c.Assets = abs
	}
	if c.Out == "" {
		c.Out = c.Label + ".apk"
	}
	c.Keystore.fill()
	return nil
}

// validPackage checks the application id against Android's rules, which are
// stricter than they look: at least two parts, each starting with a letter,
// and no Java keyword anywhere — the id becomes a package name in generated
// code, so "com.example.class" is refused by the platform, not by us.
func validPackage(p string) error {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return fmt.Errorf("apk: package %q needs at least two parts, like \"com.example.app\"", p)
	}
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("apk: package %q has an empty part", p)
		}
		if !isLetter(part[0]) && part[0] != '_' {
			return fmt.Errorf("apk: package part %q must start with a letter", part)
		}
		for i := range len(part) {
			ch := part[i]
			if !isLetter(ch) && !isDigit(ch) && ch != '_' {
				return fmt.Errorf("apk: package part %q may only hold letters, digits and _", part)
			}
		}
		if javaKeywords[part] {
			return fmt.Errorf("apk: package part %q is a Java keyword, which Android refuses", part)
		}
	}
	return nil
}

func isLetter(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }
func isDigit(c byte) bool  { return c >= '0' && c <= '9' }

var javaKeywords = map[string]bool{
	"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
	"case": true, "catch": true, "char": true, "class": true, "const": true,
	"continue": true, "default": true, "do": true, "double": true, "else": true,
	"enum": true, "extends": true, "final": true, "finally": true, "float": true,
	"for": true, "goto": true, "if": true, "implements": true, "import": true,
	"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
	"new": true, "package": true, "private": true, "protected": true, "public": true,
	"return": true, "short": true, "static": true, "strictfp": true, "super": true,
	"switch": true, "synchronized": true, "this": true, "throw": true, "throws": true,
	"transient": true, "try": true, "void": true, "volatile": true, "while": true,
	"true": true, "false": true, "null": true, "_": true,
}
