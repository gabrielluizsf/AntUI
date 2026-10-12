//go:build android

package app

// ShimClass is the Java class this library ships, if the app was built with
// it. An app built with -no-shim has none, and everything that depends on it
// reports [ErrNoShim].
const ShimClass = "dev/antui/AntuiActivity"
