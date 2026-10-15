package svg

import (
	"net/url"
	"path/filepath"
)

// filePathFromURL is the path on this machine a `file:` URL names. On Windows
// a drive letter is written `file:///C:/dir/x.png`, so the URL's path is
// `/C:/dir/x.png` and the slash before the letter is not part of the path.
// On Unix a path never starts with `/X:`, so nothing changes there.
func filePathFromURL(u *url.URL) string {
	p := u.Path
	if len(p) >= 3 && p[0] == '/' && isASCIIAlpha(p[1]) && p[2] == ':' {
		p = p[1:]
	}
	return filepath.FromSlash(p)
}
