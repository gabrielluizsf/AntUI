//go:build android

package app

// Has reports whether one particular permission was granted.
func (r PermissionResult) Has(name string) bool {
	for i, p := range r.Permissions {
		if p == name && i < len(r.Granted) {
			return r.Granted[i]
		}
	}
	return false
}
