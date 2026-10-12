//go:build android

package app

// All reports whether every permission asked for was granted.
func (r PermissionResult) All() bool {
	if len(r.Granted) == 0 {
		return false
	}
	for _, ok := range r.Granted {
		if !ok {
			return false
		}
	}
	return true
}
