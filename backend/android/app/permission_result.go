//go:build android

package app

// PermissionResult is the answer to a permission request.
type PermissionResult struct {
	// Permissions is what was asked for, in the order it was asked.
	Permissions []string
	// Granted is the answer for each, in the same order.
	Granted []bool
}
