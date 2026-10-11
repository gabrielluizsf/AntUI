package svg

// PictureAlign is which of the nine points of the box the picture is pinned to
// when the fit leaves something over: an edge or the middle of each axis,
// spelled in the file as `xMinYMid` and its eight siblings. It is read as the
// one number it maps to — the axis positions come out of it with [PictureAlign.offsets].
type PictureAlign uint8
