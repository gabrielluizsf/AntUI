//go:build android

package audio

// Kind is how long the app expects to want the sound for. Asking for the
// right one is not politeness: it is what tells the music app whether to
// stop or to duck, and getting it wrong is why some apps kill a podcast to
// play a notification.
type Kind int32
