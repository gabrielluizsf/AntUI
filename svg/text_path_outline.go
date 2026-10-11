package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// textPathOutline is the outline a `<textPath>` writes along: the href followed
// the way a `<use>` follows its own, to a `<path>` this drawing has. Anything
// else — no href, one into another drawing, an id the drawing does not have, a
// shape that is not a path — has nowhere to put the writing, and says so rather
// than leaving it to come out somewhere the file did not ask for.
func (img *Image) textPathOutline(e *element, warn func(string, ...any)) *canvas.Path {
	href := e.attr("href")
	if href == "" {
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is no outline to write along")
		return nil
	}
	id, ok := strings.CutPrefix(href, "#")
	if !ok || id == "" {
		warn("the href %q names something in another drawing, which is not read", href)
		return nil
	}
	target := img.ids[id]
	if target == nil {
		warn("the href names #%s, which the drawing does not have", id)
		return nil
	}
	if target.Name != "path" {
		warn("the href names a <%s>, and writing follows a <path>", target.Name)
		return nil
	}
	p := pathData(target.attr("d"), warn)
	if p == nil || p.Empty() {
		warn("the outline the writing follows has nothing in it")
		return nil
	}
	return p
}
