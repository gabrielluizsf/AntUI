package css

import (
	"path/filepath"
)

// Parse turns CSS source text into a Sheet. Unknown declarations are skipped
// and reported through Warn; a malformed selector or an unbalanced block is an
// error.
//
// An @import read this way has no file of its own for a relative URL to
// resolve against, so it is reported and skipped; [ParseFile] reads a file,
// and reads its imports from beside it.
func Parse(src string) (*Sheet, error) {
	return parseText(src, "")
}

// parseText reads source text that came from path — "" when it came from no
// file at all — into a Sheet, with the directory of path standing for where a
// relative @import resolves to.
func parseText(src, path string) (*Sheet, error) {
	sh := &Sheet{}
	p := &parser{src: src}
	if path != "" {
		p.path = filepath.Clean(path)
		p.seen = map[string]bool{p.path: true}
	}
	return sh, p.parseAll(sh, Media{})
}
