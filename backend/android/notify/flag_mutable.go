//go:build android

package notify

// flagMutable lets the system add to the intent, which a reply needs.
const flagMutable int32 = 1 << 25
