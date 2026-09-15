// Package biometric asks the user to prove who they are with a fingerprint
// or a face.
//
// It never sees the fingerprint. The platform's prompt does the asking and
// the matching in hardware the app cannot reach, and answers yes or no —
// which is the whole point, and is why an app must not build its own.
package biometric
