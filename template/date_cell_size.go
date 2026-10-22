package template

// dateCellSize is the side of one day cell: a two-digit day number in the
// built-in face is 16u wide, so a cell needs at least that much plus the
// bearing a round indicator wants on every side of it to read as its own
// shape rather than a tile glued to its neighbours.
func dateCellSize(u int) int { return 22 * u }
