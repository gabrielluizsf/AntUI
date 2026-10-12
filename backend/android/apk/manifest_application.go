package apk

import (
	"fmt"
	"strings"
)

// manifestApplication writes the application element and the activity that
// hosts the app, up to its launcher intent filter.
func manifestApplication(b *strings.Builder, c *Config) {
	fmt.Fprintf(b, "  <application android:label=%s\n", attr(c.Label))
	if c.Icon != "" {
		// Both, because a launcher that asks for the round one will not
		// round the square one itself.
		b.WriteString("               android:icon=\"@mipmap/ic_launcher\"\n")
		b.WriteString("               android:roundIcon=\"@mipmap/ic_launcher_round\"\n")
	}
	fmt.Fprintf(b, "               android:hasCode=%s\n", attr(fmt.Sprint(*c.Shim)))
	b.WriteString("               android:extractNativeLibs=\"false\"")
	if c.Backup != nil {
		fmt.Fprintf(b, "\n               android:allowBackup=%s", attr(fmt.Sprint(*c.Backup)))
	}
	if c.Debuggable {
		b.WriteString("\n               android:debuggable=\"true\"")
	}
	b.WriteString(">\n")

	fmt.Fprintf(b, "    <activity android:name=%s\n", attr(c.ActivityClass()))
	fmt.Fprintf(b, "              android:label=%s\n", attr(c.Label))
	b.WriteString("              android:exported=\"true\"\n")
	b.WriteString("              android:theme=\"@android:style/Theme.NoTitleBar.Fullscreen\"\n")
	if c.Orientation != "" {
		fmt.Fprintf(b, "              android:screenOrientation=%s\n", attr(c.Orientation))
	}
	fmt.Fprintf(b, "              android:configChanges=%s>\n", attr(configChanges))

	// The name the framework passes to System.loadLibrary: no "lib", no
	// ".so". Getting this wrong is a crash at start with a message about a
	// library that does not exist, naming a file that plainly does.
	fmt.Fprintf(b, "      <meta-data android:name=\"android.app.lib_name\" android:value=%s />\n", attr(c.Lib))

	b.WriteString("      <intent-filter>\n")
	b.WriteString("        <action android:name=\"android.intent.action.MAIN\" />\n")
	b.WriteString("        <category android:name=\"android.intent.category.LAUNCHER\" />\n")
	b.WriteString("      </intent-filter>\n")
}
