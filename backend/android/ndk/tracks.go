//go:build android

package ndk

// Tracks is what the file holds.
func (m *Media) Tracks() []Track { return m.tracks }
