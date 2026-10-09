package css

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
	"os"
)

// fontSource walks an @font-face src list the way a browser does: the first
// source this canvas can open is the face. local() names a font this machine
// already has, and a woff file is a compressed container this engine does not
// unpack — both are reported and passed over, and so is a source that opens
// as no font at all.
func (p *parser) fontSource(sh *Sheet, family, src string) *canvas.Face {
	for _, part := range splitTopLevel(src) {
		kind, arg := sourceKind(part)
		what := fmt.Sprintf("%q in @font-face for %q", arg, family)
		switch kind {
		case "local":
			sh.Warn = append(sh.Warn, fmtErrf("ignoring local() %q in @font-face for %q: a font this machine already has is not read", arg, family).Error())
			continue
		case "url":
		default:
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %q in @font-face for %q: not a url() or local()", part, family).Error())
			continue
		}
		if format := sourceFormat(part); format == "woff" || format == "woff2" {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: woff is not read", what).Error())
			continue
		}
		path, warn := p.localFile(arg, what)
		if warn != "" {
			sh.Warn = append(sh.Warn, warn)
			continue
		}
		data, err := os.ReadFile(path)
		if err == nil {
			var face *canvas.Face
			if face, err = canvas.ParseFace(data, float64(DefaultFontSize)); err == nil {
				return face
			}
		}
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: %v", what, err).Error())
	}
	return nil
}
