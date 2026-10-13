package apk

import "strings"

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
