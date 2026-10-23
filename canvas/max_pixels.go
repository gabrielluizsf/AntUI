package canvas

// MaxPixels is the most pixels one canvas may hold: sixteen thousand on a
// side, a gigabyte of colour, many times any screen or texture a program
// draws. Past it the size is refused rather than allocated, because a width
// and a height read out of a file nobody vetted multiply into whatever the
// file wants — into more memory than the machine has, or past the top of an
// int and round to a few bytes that every pixel written after it overruns.
const MaxPixels = 1 << 28

// sizeFits says whether a canvas width by height may be made: both sides
// above zero and the two of them no more than [MaxPixels] together, worked out
// by division so the product never gets the chance to wrap.
func sizeFits(width, height int) bool {
	return width > 0 && height > 0 && width <= MaxPixels/height
}
