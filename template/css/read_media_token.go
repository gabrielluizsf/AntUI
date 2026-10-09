package css

import "strings"

// readMediaToken folds one media condition "(key: value)" or type word into
// the query, returning the warnings it produced.
func readMediaToken(cur *MediaQuery, t string) []string {
	if t[0] != '(' {
		switch t {
		case "screen", "all":
			// The window kind the canvas always is.
			return nil
		}
		return []string{fmtErrf("ignoring media type %q", t).Error()}
	}
	body := strings.TrimSuffix(strings.TrimPrefix(t, "("), ")")
	kv := strings.SplitN(body, ":", 2)
	if len(kv) != 2 {
		return []string{fmtErrf("ignoring media condition %q", t).Error()}
	}
	key := strings.TrimSpace(kv[0])
	val := strings.TrimSpace(kv[1])
	return applyMediaFeature(cur, mediaFeat{key: key, val: val, tok: t})
}
