package svg

import "strings"

// matchAttr says whether an element's attribute meets the condition a selector
// wrote: that it is there at all, or how its value compares. The comparison is
// by case unless the selector asked for `i`.
func matchAttr(e *element, s simpleSelector) bool {
	v, ok := attrUnder(e, s)
	if !ok {
		return false
	}
	val := s.val
	if s.caseless {
		v, val = strings.ToLower(v), strings.ToLower(val)
	}
	switch s.op {
	case attrExists:
		return true
	case attrEquals:
		return v == val
	case attrIncludes:
		for _, f := range strings.Fields(v) {
			if f == val {
				return true
			}
		}
	case attrDash:
		return val != "" && (v == val || strings.HasPrefix(v, val+"-"))
	case attrPrefix:
		return val != "" && strings.HasPrefix(v, val)
	case attrSuffix:
		return val != "" && strings.HasSuffix(v, val)
	case attrSubstr:
		return val != "" && strings.Contains(v, val)
	}
	return false
}
