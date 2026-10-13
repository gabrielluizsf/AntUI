package template

// datePaint is one calendar's geometry and its month, gathered so the painters
// that draw its parts agree on where they are and what day they are marking
// without passing a dozen numbers between them.
//
// x, y, w and h are the popup box, u the scale everything in it is measured at
// (its own, once it grew to fill the window), and the rest is the month: which
// year and month it is, which weekday its first day falls on, how many days it
// has, and the day numbers to mark as chosen, today and hovered — 0 when not
// applicable.
//
// Every template's calendar is painted from one of these — the built-in one,
// the neon one and the CSS one — which is why the shape of a month is stated
// once here rather than in each painter's parameter list.
type datePaint struct {
	x, y, w, h  int
	u           int
	month, year int
	firstWD     int
	days        int
	selected    int
	today       int
	hover       int
}

// header is the height of the title row the month and the arrows sit in.
func (p datePaint) header() int { return dateHeaderHeight(p.u) }

// headerTextY is the top of the title row's text, centred in it around a
// line `th` tall.
func (p datePaint) headerTextY(th int) int { return p.y + (p.header()-th)/2 }

// weekHeadY is the top of the weekday-name row under the header.
func (p datePaint) weekHeadY() int { return weekdayHeadY(p.x, p.y, p.u) }

// rows is how many week rows the month actually needs, up to six.
func (p datePaint) rows() int { return min((p.firstWD+p.days+6)/7, 6) }
