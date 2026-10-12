//go:build android

package notify

// Importance is how loudly a channel announces itself. The user can change
// it afterwards and the app cannot change it back — a channel's importance
// is a request made once, when the channel is created, and never again.
type Importance int32
