package svg

// attrOpOf reads the operator of an attribute selector and how many bytes it
// took, or a width of none.
func attrOpOf(s string) (attrOp, int) {
	if s == "" {
		return 0, 0
	}
	if s[0] == '=' {
		return attrEquals, 1
	}
	if len(s) > 1 && s[1] == '=' {
		switch s[0] {
		case '~':
			return attrIncludes, 2
		case '|':
			return attrDash, 2
		case '^':
			return attrPrefix, 2
		case '$':
			return attrSuffix, 2
		case '*':
			return attrSubstr, 2
		}
	}
	return 0, 0
}
