package template

// dateCellSize is the side of one day cell: a two-digit day number in the
// built-in face is 16u wide, so a cell needs at least that much plus a
// little bearing to keep the digits off its neighbours.
func dateCellSize(u int) int { return 18 * u }
