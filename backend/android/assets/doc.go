// Package assets reads the files packed into the APK, as an io/fs.FS.
//
// Assets are not files. They are entries inside the package, which is a zip
// the system never unpacks — so no path leads to one and os.Open will not
// find one. What this package does is make that difference stop mattering:
// code that read its data with os.ReadFile on a desktop reads it with
// fs.ReadFile here, and the same code serves both.
//
// Put them in the APK with antuiapk's -assets flag.
//
// # What is in there besides your files
//
// The space an app reads from is not only the app's: the framework merges
// its own assets in, so a walk over the root finds a boot logo, a clock font
// and a web error page as well as whatever was packed. Read by name and it
// never matters; walk the tree and it does.
package assets
