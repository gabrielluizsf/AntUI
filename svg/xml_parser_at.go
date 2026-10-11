package svg

import "strings"

// at reports whether the text at the cursor starts with what.
func (p *xmlParser) at(what string) bool { return strings.HasPrefix(p.src[p.pos:], what) }
