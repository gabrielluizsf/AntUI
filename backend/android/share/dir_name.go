//go:build android

package share

// The one directory the provider serves, under the app's cache. It matches
// what the Java side opens; the two have to agree and there is nowhere else
// for them to agree in.
const dirName = "antui-share"
