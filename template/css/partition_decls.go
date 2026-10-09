package css

import "strings"

// partitionDecls splits the cascade's declarations into normal and !important
// buckets, skipping custom properties and keeping source order.
func partitionDecls(got []candidate) (normal, imp []Declaration) {
	for _, cand := range got {
		for _, d := range cand.decls {
			prop := strings.ToLower(strings.TrimSpace(d.Prop))
			if strings.HasPrefix(prop, "--") {
				continue
			}
			if d.Important {
				imp = append(imp, d)
			} else {
				normal = append(normal, d)
			}
		}
	}
	return normal, imp
}
