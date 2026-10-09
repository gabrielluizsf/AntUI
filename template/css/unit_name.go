package css

// unitName maps every unit to the CSS suffix that spells it.
var unitName = map[unit]string{
	unitPx:   "px",
	unitPct:  "%",
	unitEm:   "em",
	unitRem:  "rem",
	unitVw:   "vw",
	unitVh:   "vh",
	unitVmin: "vmin",
	unitVmax: "vmax",
	unitCh:   "ch",
	unitEx:   "ex",
	unitCm:   "cm",
	unitMm:   "mm",
	unitIn:   "in",
	unitPt:   "pt",
	unitPc:   "pc",
	unitQ:    "q",
	unitAuto: "",
	unitNone: "",
}
