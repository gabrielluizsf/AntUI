//go:build android

package ndk

// OpenAudio starts a stream going out to the speakers. capacity is how many
// frames the device should buffer, and 0 lets it choose.
func OpenAudio(rate, channels, capacity int) (*AudioStream, error) {
	return openAudio(rate, channels, capacity, false)
}
