package template

// multiPiece is one fragment: the column it lands in, how far down that column
// it starts and how tall the run of it is.
type multiPiece struct {
	col int
	y   int
	h   int
}
