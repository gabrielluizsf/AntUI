package css

import "strings"

// resolveCustoms folds the cascade's custom-property declarations (--foo)
// into one map, in cascade order so a later declaration overwrites.
func resolveCustoms(got []candidate) map[string]string {
	customs := make(map[string]string, 8)
	for _, cand := range got {
		for _, d := range cand.decls {
			prop := strings.ToLower(strings.TrimSpace(d.Prop))
			if strings.HasPrefix(prop, "--") {
				customs[prop] = d.Raw
			}
		}
	}
	return customs
}
