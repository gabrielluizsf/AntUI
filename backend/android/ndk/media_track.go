//go:build android

package ndk

// Track finds the first track of a kind, and reports whether there was one.
func (m *Media) Track(video bool) (Track, bool) {
	for _, t := range m.tracks {
		if t.Video() == video && t.Audio() != video {
			return t, true
		}
	}
	return Track{}, false
}
