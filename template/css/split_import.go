package css

import (
	"strings"
)

// splitImport takes an @import prelude apart: the URL it names, the media
// query that gates the whole of it, and the layer() and supports() clauses
// between the two — which this engine does not fold, and reports rather than
// drop the import along with them.
func splitImport(prelude string) (ref, media string, warns []string) {
	s := strings.TrimSpace(prelude)
	switch {
	case len(s) >= 4 && strings.EqualFold(s[:4], "url("):
		end := parenEnd(s, 3)
		if end < 0 {
			return "", "", []string{fmtErrf("ignoring an @import whose url( never ends").Error()}
		}
		ref = unquote(strings.TrimSpace(s[4:end]))
		s = strings.TrimSpace(s[end+1:])
	case len(s) > 0 && (s[0] == '"' || s[0] == '\''):
		end := quoteEnd(s, 0)
		if end < 0 {
			return "", "", []string{fmtErrf("ignoring an @import with an unterminated url").Error()}
		}
		ref = s[1:end]
		s = strings.TrimSpace(s[end+1:])
	default:
		return "", "", []string{fmtErrf("ignoring an @import with no file name").Error()}
	}
	// A query and a fragment are for a server to answer; localFile strips
	// them when it comes to open the file.
	if strings.TrimSpace(ref) == "" {
		return "", "", []string{fmtErrf("ignoring an @import with no file name").Error()}
	}
	for s != "" {
		s = strings.TrimLeft(s, " \t\r\n")
		clause := clauseWord(s)
		if clause == "" {
			break
		}
		s = s[len(clause):]
		if strings.HasPrefix(s, "(") {
			end := parenEnd(s, 0)
			if end < 0 {
				warns = append(warns, fmtErrf("ignoring an @import whose %s( never ends", clause).Error())
				return ref, "", warns
			}
			s = s[end+1:]
		}
		warns = append(warns, fmtErrf("ignoring the %s() clause of @import", clause).Error())
	}
	return ref, strings.TrimSpace(s), warns
}
