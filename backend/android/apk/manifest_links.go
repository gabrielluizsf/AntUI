package apk

import (
	"fmt"
	"strings"
)

// manifestLinks writes one intent filter per deep link and closes the
// activity. One filter per link rather than one filter with every link in
// it: a single filter matches the *cross product* of its schemes, hosts and
// paths — two links in one filter would also claim the two combinations
// nobody asked for.
func manifestLinks(b *strings.Builder, c *Config) {
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
		fmt.Fprintf(b, "        <data android:scheme=%s", attr(l.Scheme))
		if l.Host != "" {
			fmt.Fprintf(b, " android:host=%s", attr(l.Host))
		}
		if l.Path != "" {
			fmt.Fprintf(b, " android:pathPrefix=%s", attr(l.Path))
		}
		b.WriteString(" />\n")
		b.WriteString("      </intent-filter>\n")
	}
	b.WriteString("    </activity>\n")
}
