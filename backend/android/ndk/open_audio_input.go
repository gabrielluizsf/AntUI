//go:build android

package ndk

// OpenAudioInput starts a stream coming in from the microphone.
//
// It needs the RECORD_AUDIO permission, and without it the open fails rather
// than returning silence — which is the right way round, and not what every
// platform does.
func OpenAudioInput(rate, channels, capacity int) (*AudioStream, error) {
	return openAudio(rate, channels, capacity, true)
}
