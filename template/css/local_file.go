package css

import (
	"path/filepath"
	"strings"
)

// localFile resolves a URL from the file being read to a path on this
// machine, or returns the warning that says why it is not one. what names
// where the URL came from for that warning.
func (p *parser) localFile(ref, what string) (string, string) {
	if isRemoteRef(ref) {
		return "", fmtErrf("ignoring %s: only local files are read", what).Error()
	}
	// A query and a fragment are for a server to answer; the file beside
	// this one is the file beside this one.
	if i := strings.IndexAny(ref, "?#"); i >= 0 {
		ref = ref[:i]
	}
	if strings.TrimSpace(ref) == "" {
		return "", fmtErrf("ignoring %s: no file name", what).Error()
	}
	if filepath.IsAbs(ref) {
		return filepath.Clean(ref), ""
	}
	if p.path == "" {
		return "", fmtErrf("ignoring %s: no file to resolve it against", what).Error()
	}
	return filepath.Join(filepath.Dir(p.path), ref), ""
}
