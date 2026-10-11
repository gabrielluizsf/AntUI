package canvas

// bayer is the classic 4x4 threshold matrix.
var bayer = [16]uint32{
	0, 8, 2, 10,
	12, 4, 14, 6,
	3, 11, 1, 9,
	15, 7, 13, 5,
}
