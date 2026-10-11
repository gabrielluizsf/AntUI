package template

import "github.com/gabrielluizsf/antui/template/css"

// resolveGridAxis turns a parsed template into the explicit tracks of one axis,
// working out how many repetitions an auto-fill or auto-fit stands for and
// carrying the line names along with them. The rows of a container that does
// not declare its height have no space to fill, so an auto-repeat there keeps
// the single track the grammar promises.
func (b *gridBatch) resolveGridAxis(template css.GridTemplate, available, gap int) css.GridAxis {
	// A track minimum is a length in the sheet's own pixels, so the space it is
	// counted against has to be in those too. A window drawn at twice its size
	// measures twice as wide on the canvas, and without the conversion an
	// auto-repeat would fit twice as many tracks as the room holds.
	unit := max(Scale(b.css.win), 1)
	return template.Resolve(available/unit, gap/unit, b.gridInFlow())
}
