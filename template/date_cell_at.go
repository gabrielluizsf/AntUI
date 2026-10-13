package template

// dateCell is where one day of a calendar lands: the square the day is painted
// in, and where its number goes inside it. Gathered so the painter that draws
// the cell and the one that draws the number agree on the same place without
// passing five numbers plus the day between them.
type dateCell struct {
	x, y, w, h int // the square the day is painted in
	tx, ty     int // where the number sits inside it, centred
}

// dateCellAt measures the cell of day `day` of a month laid out at p. The
// number's top-left is the square's, centred on the text at the popup's scale.
func dateCellAt(p datePaint, day, textW, textH int) dateCell {
	pos := p.firstWD + day - 1
	x, y, w, h := dayCellRect(p.x, p.y, p.u, pos%7, pos/7)
	return dateCell{x: x, y: y, w: w, h: h, tx: x + (w-textW)/2, ty: y + (h-textH)/2}
}
