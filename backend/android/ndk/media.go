//go:build android

package ndk

// Media is an open media file, taken apart into its tracks.
type Media struct {
	path   string
	tracks []Track
}
