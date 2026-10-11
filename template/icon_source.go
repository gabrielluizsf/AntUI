package template

// IconSource is the SVG text a named icon was registered from, which is the
// drawing to start from when changing it. The second answer says whether there
// is an icon of that name at all.
func IconSource(name string) (string, bool) {
	iconMu.RLock()
	defer iconMu.RUnlock()
	e := icons[name]
	if e == nil {
		return "", false
	}
	return e.src, true
}
