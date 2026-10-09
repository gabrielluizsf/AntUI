package css

import (
	"fmt"
	"os"
)

// importRef is one @import's target: the URL as written and the media query
// gating it.
type importRef struct {
	name  string
	media string
}

// importSheet opens the file an @import names, refusing what is remote, too
// deep, already read or unreadable, and parses it into the same sheet under
// the media the import carries.
func (p *parser) importSheet(sh *Sheet, outer Media, ref importRef) error {
	what := fmt.Sprintf("@import of %q", ref.name)
	path, warn := p.localFile(ref.name, what)
	if warn != "" {
		sh.Warn = append(sh.Warn, warn)
		return nil
	}
	if p.depth >= maxImportDepth {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: more than %d files deep", what, maxImportDepth).Error())
		return nil
	}
	if p.seen == nil {
		p.seen = make(map[string]bool)
	}
	if p.seen[path] {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: already read", what).Error())
		return nil
	}
	p.seen[path] = true
	data, err := os.ReadFile(path)
	if err != nil {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: %v", what, err).Error())
		return nil
	}
	m, mediaWarns := parseMedia(ref.media)
	sh.Warn = append(sh.Warn, mediaWarns...)
	sub := &parser{src: string(data), path: path, seen: p.seen, depth: p.depth + 1}
	return sub.parseAll(sh, outer.and(m))
}
