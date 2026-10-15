package svg

import (
	"path/filepath"
	"strings"
)

// joinLocal completes a relative href against a base that is a path on this
// machine. A base that ends in a separator is a directory; any other base is
// the file the drawing was read from, so its directory is used. The trailing
// separator of the href is kept, since it is what makes an xml:base a folder.
func joinLocal(base, href string) string {
	if filepath.IsAbs(href) || strings.HasPrefix(href, "/") || strings.HasPrefix(href, `\`) {
		return href
	}
	dir := base
	if last := base[len(base)-1]; last != '/' && last != '\\' {
		dir = filepath.Dir(base)
	}
	out := filepath.Join(dir, filepath.FromSlash(href))
	if last := href[len(href)-1]; last == '/' || last == '\\' {
		out += string(filepath.Separator)
	}
	return out
}
