package apk

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// manifestShim writes the receiver a notification's action buttons broadcast
// to and the provider that hands a file to another app. The receiver is not
// exported, and from API 31 a receiver with an intent filter must say which
// it is. The provider is not exported either — nothing may query it
// uninvited — and grants permission per address instead: an app that is
// *given* a content:// address may read that one file, and only until it is
// revoked.
func manifestShim(b *strings.Builder, c *Config) {
	if c.Shim == nil || !*c.Shim {
		return
	}
	fmt.Fprintf(b, "    <receiver android:name=%s android:exported=\"false\" />\n",
		attr(shim.Receiver))
	fmt.Fprintf(b, "    <provider android:name=%s\n", attr(shim.Provider))
	fmt.Fprintf(b, "              android:authorities=%s\n",
		attr(c.Package+shim.ProviderSuffix))
	b.WriteString("              android:exported=\"false\"\n")
	b.WriteString("              android:grantUriPermissions=\"true\" />\n")
}
