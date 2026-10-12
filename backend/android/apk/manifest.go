package apk

import "strings"

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
	manifestOpen(&b, c)
	manifestApplication(&b, c)
	manifestLinks(&b, c)
	manifestShim(&b, c)
	b.WriteString("  </application>\n")
	b.WriteString("</manifest>\n")
	return b.String()
}
