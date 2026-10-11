package svg

// bolderWeight and lighterWeight are the steps the two words take from the
// weight in hand, indexed by the hundred that weight sits on, from 100 to
// 900. They are CSS's own table: only four weights take part in it — thin,
// normal, bold and heavy — so `bolder` on a normal weight is bold rather than
// a step of a hundred, which is also the only step that draws differently here.
var (
	bolderWeight  = [10]uint16{0, 400, 400, 400, 700, 700, 900, 900, 900, 900}
	lighterWeight = [10]uint16{0, 100, 100, 100, 100, 100, 400, 400, 700, 700}
)
