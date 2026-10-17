package svg

// entityAt reads the entity at the `&` of s[i] — one of the five XML defines,
// a character written as `&#…;`, or a name the drawing's own doctype defines
// — and answers what it stands for and how many bytes it took.
func (p *xmlParser) entityAt(s string, i int, over []string) (string, int, error) {
	if i+1 < len(s) && s[i+1] == '#' {
		return charAt(s, i)
	}
	name, size, err := refName(s, i)
	if err != nil {
		return "", 0, err
	}
	switch name {
	case "lt":
		return "<", size, nil
	case "gt":
		return ">", size, nil
	case "amp":
		return "&", size, nil
	case "apos":
		return "'", size, nil
	case "quot":
		return "\"", size, nil
	}
	return p.definedEntity(name, size, over)
}
