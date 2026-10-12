package template

import "github.com/gabrielluizsf/antui/template/css"

// flexBatch is one flex layout: the container's style and box on the page,
// and the items collected while its children drew. The silent pass fills
// items; solve turns them into boxes; the draw pass hands those boxes back.
// Nested [CSS.Flex] blocks flatten into the enclosing batch, so a container
// is always drawn by the Flex that started it.
type flexBatch struct {
	css      *CSS
	st       css.Style
	x, y, w  int // container border box
	h        int // solved after content has a size
	originX  int // content-box top-left in window coordinates
	originY  int
	contentW int // content-box width, fixed by the flow
	contentH int // content-box height, solved by the distribution
	frameW   int // border plus padding, left and right
	frameH   int // border plus padding, top and bottom
	frameTop int // border plus padding above the content
	items    []flexItem
	next     int // the item the draw pass is handing back
}
