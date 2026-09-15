// Package permission asks the user for the things Android will not give an
// app for the asking: the camera, the microphone, where the device is.
//
// A permission has to be in the manifest as well as requested here — the
// request for one the manifest does not declare is refused immediately and
// silently, which looks exactly like the user saying no. antuiapk's -perms
// flag is what puts them in the manifest.
package permission
