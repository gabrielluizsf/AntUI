package css

// unitSuffixes lists every unit suffix the parser accepts, longest first so a
// suffix is not swallowed by a shorter one: "rem" must win over "em", and
// "vmax" over "vh". "px" is a suffix like any other here.
var unitSuffixes = []struct {
	name string
	u    unit
}{
	{"vmin", unitVmin},
	{"vmax", unitVmax},
	{"rem", unitRem},
	{"cm", unitCm},
	{"mm", unitMm},
	{"in", unitIn},
	{"pt", unitPt},
	{"pc", unitPc},
	{"px", unitPx},
	{"em", unitEm},
	{"vw", unitVw},
	{"vh", unitVh},
	{"ch", unitCh},
	{"ex", unitEx},
	{"q", unitQ},
}
