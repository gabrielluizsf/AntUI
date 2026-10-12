//go:build android

package app

// Event is one thing that happened to the app. They arrive in the order
// Android raised them, one at a time, from [App.Next].
type Event int
