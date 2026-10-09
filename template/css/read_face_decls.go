package css

import "strings"

// faceDecls collects the four properties this engine reads out of a
// @font-face block, warning on anything else and passing it over.
func (p *parser) faceDecls(sh *Sheet) (family, src, weight, style string, err error) {
	err = p.declBlock(func(text string) error {
		d, ok := parseDeclaration(text)
		if !ok {
			return nil
		}
		switch prop := strings.ToLower(strings.TrimSpace(d.Prop)); prop {
		case "font-family":
			family = d.Raw
		case "src":
			src = d.Raw
		case "font-weight":
			weight = d.Raw
		case "font-style":
			style = d.Raw
		case "font-display":
			// How long a stylesheet waits for a face it did not read is
			// answered here: this canvas reads the file now or never, so
			// there is nothing for the hint to time.
		default:
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %q in @font-face", prop).Error())
		}
		return nil
	})
	return
}
