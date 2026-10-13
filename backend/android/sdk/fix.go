package sdk

import "strings"

// Fix is the one sdkmanager command that installs everything missing, or the
// empty string when nothing missing can be installed that way.
func (m *Missing) Fix() string {
	var pkgs []string
	for _, p := range m.Problems {
		if p.Install != "" {
			pkgs = append(pkgs, `"`+p.Install+`"`)
		}
	}
	if len(pkgs) == 0 {
		return ""
	}
	return "sdkmanager " + strings.Join(pkgs, " ")
}
