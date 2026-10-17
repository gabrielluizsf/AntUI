package svg

import (
	"fmt"
	"strings"
)

// declaration reads the `<?xml … ?>` that may open a document and holds to
// what it says. The version has to be one this reads and the encoding has to
// be one this can turn into characters: the text is read as UTF-8 or UTF-16
// and nothing else, so a document declaring another would come out as
// characters other than the ones it wrote — a tree different from the
// declared one, which this package does not draw quietly.
func (p *xmlParser) declaration() error {
	p.pos += len("<?xml")
	i := strings.Index(p.src[p.pos:], "?>")
	if i < 0 {
		return fmt.Errorf("antui/svg: the xml declaration is never closed")
	}
	body := p.src[p.pos : p.pos+i]
	p.pos += i + len("?>")
	return checkDecl(body)
}

// checkDecl holds the declaration's own words to them: a version of 1.0 or
// 1.1, which a declaration has to carry, and an encoding this reads, where
// none written means UTF-8.
func checkDecl(decl string) error {
	said := declAttrs(decl)
	if v := said["version"]; v != "1.0" && v != "1.1" {
		return fmt.Errorf("antui/svg: the xml declaration says version %q, which is not one this reads", v)
	}
	if enc := said["encoding"]; enc != "" && !readableEncoding(enc) {
		return fmt.Errorf("antui/svg: the drawing is written in %s, which is not a reading this package has", enc)
	}
	return nil
}

// readableEncoding says the declared encoding is one this package reads: the
// text itself, in the spellings it is written by.
func readableEncoding(enc string) bool {
	switch strings.ToLower(strings.ReplaceAll(enc, "_", "-")) {
	case "utf-8", "utf8", "utf-16", "utf16":
		return true
	}
	return false
}
