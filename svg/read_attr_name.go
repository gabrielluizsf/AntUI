package svg

// readAttrName reads the name of an attribute selector, with the namespace
// prefix one may carry — `[svg|href]`, `[*|href]`, `[|href]`. It answers the
// name, the prefix as written (`*` for the wildcard, empty for `|` or for no
// prefix at all), whether a prefix was written at all, and where it ends.
func readAttrName(s string, i int) (name, prefix string, named bool, next int) {
	if i < len(s) && s[i] == '*' && i+1 < len(s) && s[i+1] == '|' && (i+2 >= len(s) || s[i+2] != '=') {
		name, next = readLocalAttrName(s, i+2)
		return name, "*", true, next
	}
	if i < len(s) && s[i] == '|' && (i+1 >= len(s) || s[i+1] != '=') {
		name, next = readLocalAttrName(s, i+1)
		return name, "", true, next
	}
	name, next = readSelIdent(s, i)
	if next < len(s) && s[next] == '|' && (next+1 >= len(s) || s[next+1] != '=') {
		local, after := readLocalAttrName(s, next+1)
		return local, name, true, after
	}
	return name, "", false, next
}

// readLocalAttrName reads the name after an attribute selector's namespace
// prefix. `*` is not a name of its own, so `[*|*]` reads as none.
func readLocalAttrName(s string, i int) (string, int) {
	if i < len(s) && s[i] == '*' {
		return "", i + 1
	}
	return readSelIdent(s, i)
}
