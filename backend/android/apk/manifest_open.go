package apk

import (
	"fmt"
	"strings"
)

// manifestOpen writes the manifest element, the permissions and features it
// asks for, and the queries block that says what else installed the app is
// allowed to see — without which an app targeting API 30 finds nothing.
func manifestOpen(b *strings.Builder, c *Config) {
	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>` + "\n")
	b.WriteString(`<manifest xmlns:android="http://schemas.android.com/apk/res/android"` + "\n")
	fmt.Fprintf(b, "    package=%s>\n", attr(c.Package))

	for _, p := range c.Permissions {
		fmt.Fprintf(b, "  <uses-permission android:name=%s />\n", attr(p))
	}
	for _, f := range c.Features {
		fmt.Fprintf(b, "  <uses-feature android:name=%s android:required=%s />\n",
			attr(f.Name), attr(fmt.Sprint(f.Required)))
	}

	if len(c.Queries) > 0 {
		b.WriteString("  <queries>\n")
		for _, action := range c.Queries {
			b.WriteString("    <intent>\n")
			fmt.Fprintf(b, "      <action android:name=%s />\n", attr(action))
			b.WriteString("    </intent>\n")
		}
		b.WriteString("  </queries>\n")
	}
}
