// Package gallery puts pictures, videos and sounds into the device's own
// library, where the gallery and every other app can see them.
//
// It is MediaStore, and the shape of it changed with Android 10. An app no
// longer writes a file somewhere and asks for it to be noticed: it asks the
// library for somewhere to write, writes there, and the entry is published
// when it is closed. That needs no permission — an app owns what it creates —
// and it is why an app that saves a photograph should no longer be asking
// for access to all of someone's photographs.
//
// [Scan] is the old way, kept for a file that already exists.
package gallery
