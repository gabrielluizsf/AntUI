package svg

import "strings"

// matchAttr says whether an element's attribute meets the condition a selector
// wrote: that it is there at all, or how its value compares.
func matchAttr(e *element, s simpleSelector) bool {
	v, ok := e.Attr[s.name]
	if !ok {
		return false
	}
	switch s.op {
	case attrExists:
		return true
	case attrEquals:
		return v == s.val
	case attrIncludes:
		for _, f := range strings.Fields(v) {
			if f == s.val {
				return true
			}
		}
	case attrDash:
		return s.val != "" && (v == s.val || strings.HasPrefix(v, s.val+"-"))
	case attrPrefix:
		return s.val != "" && strings.HasPrefix(v, s.val)
	case attrSuffix:
		return s.val != "" && strings.HasSuffix(v, s.val)
	case attrSubstr:
		return s.val != "" && strings.Contains(v, s.val)
	}
	return false
}
