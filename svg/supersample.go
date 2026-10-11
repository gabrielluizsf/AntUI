package svg

// supersample is how many pixels to either side of the one being asked for a
// drawing is painted at before being brought back down. Drawing an icon at
// twenty-four pixels from a drawing of curves and diagonals leaves the edges
// ragged at one pixel; drawn at forty-eight they are smooth, and the halving
// averages four of them into one, which is what the antialiasing of a fill
// already does and the same cost.
const supersample = 2
