package svg

// cameFromNetwork says whether this drawing, or any drawing it was read
// through, came from a server. Such a drawing is somebody else's text, and a
// browser lets it name nothing on this machine; reading a `file:` address or a
// bare path for it would hand a stranger every file the program can open, to
// probe for or to draw — the escape librsvg's CVE-2023-38633 made through an
// href. Where a drawing read from a disk points is the program's own business,
// as it always was.
func (img *Image) cameFromNetwork() bool {
	for _, from := range img.chain {
		if isNetworkScheme(hrefScheme(from)) {
			return true
		}
	}
	return false
}

// isNetworkScheme says whether a scheme names a server rather than this
// machine or bytes written in place.
func isNetworkScheme(scheme string) bool {
	return scheme == "http" || scheme == "https"
}

// readsThisMachine says whether an address, already resolved, is read off this
// machine's disk: a `file:` address, or one with no scheme at all.
func readsThisMachine(target string) bool {
	scheme := hrefScheme(target)
	return scheme == "" || scheme == "file"
}
