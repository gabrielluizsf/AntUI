package svg

// pictureHref is where the picture an `<image>` draws comes from: the plain
// href, or the older xlink:href spelling when there is none — the same
// fallback a `<use>` takes, and the one a drawing written for SVG 1.1 carries.
func pictureHref(e *element) string {
	if href := e.attr("href"); href != "" {
		return href
	}
	return e.attr("xlink:href")
}
