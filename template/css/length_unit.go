package css

// unit is the kind of measure a Length carries: reference pixels,
// percentages, font- and viewport-relative units and the fixed physical CSS
// units.
type unit uint8

const (
	unitPx unit = iota
	unitPct
	unitAuto
	unitNone
	unitEm
	unitRem
	unitVw
	unitVh
	unitVmin
	unitVmax
	unitCh
	unitEx
	unitCm
	unitMm
	unitIn
	unitPt
	unitPc
	unitQ
)

// The fixed physical units resolve against the CSS reference density of 96
// pixels per inch: 1in = 96px, 1cm = 96/2.54px and so on. These are reference
// pixels, scaled with the window like every other fixed length at draw time.
const (
	pxPerIn = 96.0
	pxPerCm = pxPerIn / 2.54
	pxPerMm = pxPerIn / 25.4
	pxPerQ  = pxPerMm / 4
	pxPerPt = pxPerIn / 72
	pxPerPc = 16.0
)
