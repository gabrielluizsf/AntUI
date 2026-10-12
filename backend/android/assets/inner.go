//go:build android

package assets

// inner is the name the platform wants: no leading slash, and "" for the
// root rather than fs's ".".
func inner(name string) string {
	if name == "." {
		return ""
	}
	return name
}
