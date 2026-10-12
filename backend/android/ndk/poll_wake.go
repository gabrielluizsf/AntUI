//go:build android

package ndk

// What Poll returns when it is not returning an identifier.
const (
	// PollWake means someone called Wake.
	PollWake = -1
	// PollCallback means the looper ran a callback itself; there is nothing
	// for the caller to do.
	PollCallback = -2
	// PollTimeout means the timeout ran out with nothing to report.
	PollTimeout = -3
	// PollError means the looper is broken.
	PollError = -4
)
