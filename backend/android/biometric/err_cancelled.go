//go:build android

package biometric

import "errors"

// The errors a prompt can end in. Cancelled and Refused are the two a caller
// has to handle; the rest are worth telling apart when reporting.
var (
	// ErrCancelled is the user backing out, or the app cancelling.
	ErrCancelled = errors.New("biometric: cancelled")
	// ErrRefused is the negative button — the user choosing another way in.
	ErrRefused = errors.New("biometric: the user chose the other button")
	// ErrLockedOut is too many failed attempts.
	ErrLockedOut = errors.New("biometric: too many attempts; locked out for now")
	// ErrNotDeclared is the manifest missing USE_BIOMETRIC.
	//
	// It is a *normal* permission: granted at install by being declared, with
	// nothing to ask the user. Without it every call here throws a
	// SecurityException, which names the permission and reads like a bug in
	// the app rather than a line missing from its manifest — so it is caught
	// and said plainly.
	ErrNotDeclared = errors.New("biometric: this app's manifest does not declare " +
		"android.permission.USE_BIOMETRIC; add it and rebuild — it is granted at " +
		"install and the user is never asked")
)
