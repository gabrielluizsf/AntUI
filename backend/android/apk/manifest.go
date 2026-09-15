package apk

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// configChanges is what the activity says it handles itself. Without it,
// Android destroys and recreates the activity on every rotation, keyboard
// change and theme switch — which for a native app means tearing the surface
// down and building it back up, losing everything the app had in memory.
// A NativeActivity is expected to cope, so it says so.
const configChanges = "orientation|screenSize|smallestScreenSize|screenLayout|" +
	"keyboard|keyboardHidden|navigation|uiMode|density|fontScale|locale|layoutDirection"

// Manifest is the AndroidManifest.xml for this app, as text. aapt2 turns it
// into the binary form an APK holds; writing that form directly is a
// separate path, for the day building without the SDK matters.
//
// Three attributes here are load-bearing and easy to get wrong:
//
//	hasCode                   whether there is a dex in this APK. It has to
//	                          match: false with a dex is code that never
//	                          loads, and true without one is an app that
//	                          dies at start looking for it.
//	extractNativeLibs="false" the library is loaded straight out of the APK
//	                          rather than unpacked into the data directory,
//	                          which halves the install size — and requires
//	                          the .so to be stored uncompressed and page
//	                          aligned. See [Assemble] and [Sign].
//	exported="true"           required on anything with an intent filter
//	                          since API 31, and an install-time error without.
func (c *Config) Manifest() string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>` + "\n")
	b.WriteString(`<manifest xmlns:android="http://schemas.android.com/apk/res/android"` + "\n")
	fmt.Fprintf(&b, "    package=%s>\n", attr(c.Package))

	for _, p := range c.Permissions {
		fmt.Fprintf(&b, "  <uses-permission android:name=%s />\n", attr(p))
	}
	for _, f := range c.Features {
		fmt.Fprintf(&b, "  <uses-feature android:name=%s android:required=%s />\n",
			attr(f.Name), attr(fmt.Sprint(f.Required)))
	}

	// What this app is allowed to see of what else is installed. Without a
	// matching entry here, an app targeting API 30 or above finds nothing —
	// see Config.Queries.
	if len(c.Queries) > 0 {
		b.WriteString("  <queries>\n")
		for _, action := range c.Queries {
			b.WriteString("    <intent>\n")
			fmt.Fprintf(&b, "      <action android:name=%s />\n", attr(action))
			b.WriteString("    </intent>\n")
		}
		b.WriteString("  </queries>\n")
	}

	fmt.Fprintf(&b, "  <application android:label=%s\n", attr(c.Label))
	if c.Icon != "" {
		// Both, because a launcher that asks for the round one will not
		// round the square one itself.
		b.WriteString("               android:icon=\"@mipmap/ic_launcher\"\n")
		b.WriteString("               android:roundIcon=\"@mipmap/ic_launcher_round\"\n")
	}
	fmt.Fprintf(&b, "               android:hasCode=%s\n", attr(fmt.Sprint(*c.Shim)))
	b.WriteString("               android:extractNativeLibs=\"false\"")
	if c.Backup != nil {
		fmt.Fprintf(&b, "\n               android:allowBackup=%s", attr(fmt.Sprint(*c.Backup)))
	}
	if c.Debuggable {
		b.WriteString("\n               android:debuggable=\"true\"")
	}
	b.WriteString(">\n")

	fmt.Fprintf(&b, "    <activity android:name=%s\n", attr(c.ActivityClass()))
	fmt.Fprintf(&b, "              android:label=%s\n", attr(c.Label))
	b.WriteString("              android:exported=\"true\"\n")
	b.WriteString("              android:theme=\"@android:style/Theme.NoTitleBar.Fullscreen\"\n")
	if c.Orientation != "" {
		fmt.Fprintf(&b, "              android:screenOrientation=%s\n", attr(c.Orientation))
	}
	fmt.Fprintf(&b, "              android:configChanges=%s>\n", attr(configChanges))

	// The name the framework passes to System.loadLibrary: no "lib", no
	// ".so". Getting this wrong is a crash at start with a message about a
	// library that does not exist, naming a file that plainly does.
	fmt.Fprintf(&b, "      <meta-data android:name=\"android.app.lib_name\" android:value=%s />\n", attr(c.Lib))

	b.WriteString("      <intent-filter>\n")
	b.WriteString("        <action android:name=\"android.intent.action.MAIN\" />\n")
	b.WriteString("        <category android:name=\"android.intent.category.LAUNCHER\" />\n")
	b.WriteString("      </intent-filter>\n")

	// One filter per link rather than one filter with every link in it. A
	// single filter matches the *cross product* of its schemes, hosts and
	// paths — two links in one filter would also claim the two combinations
	// nobody asked for.
	for _, l := range c.Links {
		if l.Verify {
			b.WriteString("      <intent-filter android:autoVerify=\"true\">\n")
		} else {
			b.WriteString("      <intent-filter>\n")
		}
		b.WriteString("        <action android:name=\"android.intent.action.VIEW\" />\n")
		b.WriteString("        <category android:name=\"android.intent.category.DEFAULT\" />\n")
		// BROWSABLE is what lets a link in a browser or a message reach the
		// app at all. Without it the filter only matches an intent another
		// app built by hand.
		b.WriteString("        <category android:name=\"android.intent.category.BROWSABLE\" />\n")
		fmt.Fprintf(&b, "        <data android:scheme=%s", attr(l.Scheme))
		if l.Host != "" {
			fmt.Fprintf(&b, " android:host=%s", attr(l.Host))
		}
		if l.Path != "" {
			fmt.Fprintf(&b, " android:pathPrefix=%s", attr(l.Path))
		}
		b.WriteString(" />\n")
		b.WriteString("      </intent-filter>\n")
	}
	b.WriteString("    </activity>\n")

	// The receiver a notification's action buttons broadcast to. It is not
	// exported: nothing outside this app has any business sending to it, and
	// from API 31 a receiver with an intent filter must say which it is.
	if c.Shim != nil && *c.Shim {
		fmt.Fprintf(&b, "    <receiver android:name=%s android:exported=\"false\" />\n",
			attr(shim.Receiver))

		// The provider that hands a file to another app. It is not exported
		// — nothing may query it uninvited — and grants permission per
		// address instead: an app that is *given* a content:// address may
		// read that one file, and only until it is revoked.
		fmt.Fprintf(&b, "    <provider android:name=%s\n", attr(shim.Provider))
		fmt.Fprintf(&b, "              android:authorities=%s\n",
			attr(c.Package+shim.ProviderSuffix))
		b.WriteString("              android:exported=\"false\"\n")
		b.WriteString("              android:grantUriPermissions=\"true\" />\n")
	}
	b.WriteString("  </application>\n")
	b.WriteString("</manifest>\n")
	return b.String()
}

// ActivityClass is the activity the manifest names: the shim when there is
// one, and the platform's own NativeActivity when there is not.
func (c *Config) ActivityClass() string {
	// nil is the default and the default is on, which is what [Config.check]
	// fills in. Answering "NativeActivity" for nil here made this disagree
	// with the manifest that check produces, so a caller who asked before
	// building — anything that installs and then starts the app — started an
	// activity the package does not declare.
	if c.Shim == nil || *c.Shim {
		return shim.Activity
	}
	return "android.app.NativeActivity"
}

// attr writes an XML attribute value, quotes and all, with everything
// escaped. It goes through encoding/xml rather than a handful of Replaces so
// that a label with an apostrophe or a permission from a config file cannot
// produce a manifest that is not XML.
func attr(s string) string {
	var buf bytes.Buffer
	buf.WriteByte('"')
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		// EscapeText only fails if the writer does, and a bytes.Buffer does
		// not.
		panic(err)
	}
	buf.WriteByte('"')
	return buf.String()
}
