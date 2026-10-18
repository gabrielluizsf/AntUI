package svg

import "strings"

// matchLang says whether a language a drawing stands in meets one of the
// ranges `:lang(…)` named. A range of `en` stands for `en` and everything
// beneath it, `en-US` and `en-US-x` as well; a `*` in a range stands for
// whatever subtag sits there, which is how CSS writes a wider question. A
// drawing that gives no language at all is in none, so it meets no range.
func matchLang(lang string, ranges []string) bool {
	if lang == "" {
		return false
	}
	lang = strings.ToLower(lang)
	for _, r := range ranges {
		if langRangeMatches(strings.ToLower(r), lang) {
			return true
		}
	}
	return false
}

// langRangeMatches answers whether one language range covers one language.
func langRangeMatches(r, lang string) bool {
	if r == "*" || r == lang {
		return true
	}
	for r != "" {
		if lang == "" {
			return false
		}
		var rp, lp string
		rp, r = rangeNext(r)
		lp, lang = rangeNext(lang)
		if rp != "*" && rp != lp {
			return false
		}
	}
	return true
}

// rangeNext reads the first subtag before the next dash, and what is left of
// the range after it.
func rangeNext(s string) (string, string) {
	if i := strings.IndexByte(s, '-'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}
