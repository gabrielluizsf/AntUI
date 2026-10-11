package svg

import "strings"

// useTarget is the element a `<use>` names and the id it named it by. Only the
// whole drawing can say what an id is, so the reference is followed here
// rather than where it is written. A reference that cannot be followed — none
// at all, one into another drawing, or one to an id this drawing does not have
// — says so and stands for nothing, since a `<use>` with nothing behind it is
// a hole in the picture either way.
func (img *Image) useTarget(e *element, warn func(string, ...any)) (*element, string) {
	href := e.attr("href")
	if href == "" {
		// The older spelling of the same attribute, which is the one a drawing
		// written for SVG 1.1 carries.
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is nothing for it to draw")
		return nil, ""
	}
	id, ok := strings.CutPrefix(href, "#")
	if !ok {
		warn("the href %q names something in another drawing, which is not read", href)
		return nil, ""
	}
	if id == "" {
		warn("the href %q names nothing at all", href)
		return nil, ""
	}
	target := img.ids[id]
	if target == nil {
		warn("the href names #%s, which the drawing does not have", id)
		return nil, ""
	}
	return target, id
}
