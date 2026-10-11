//go:build android

package antui

// fontsXML is the platform's own font configuration, in the order the
// platform has kept it.
//
// /system/etc/fonts.xml is where it has been since API 21. The one before it
// is listed because a device that old still runs this library, and /product
// is where a manufacturer puts its own — a phone whose maker replaced the
// interface font says so there and nowhere else.
var fontsXML = []string{
	"/system/etc/fonts.xml",
	"/product/etc/fonts.xml",
	"/system/etc/system_fonts.xml",
}
