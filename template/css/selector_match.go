package css

import (
	"strings"
)

func (s Selector) match(tag string, classes []string, state State) bool {
	if s.Unsupported != "" {
		return false
	}
	tag = strings.ToLower(tag)
	if s.All {
		if s.Tag != "" && s.Tag != tag {
			return false
		}
	} else if s.Tag != "" && s.Tag != tag {
		return false
	}
	for _, c := range s.Classes {
		found := false
		for _, have := range classes {
			if have == strings.ToLower(c) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return s.State&state == s.State
}
