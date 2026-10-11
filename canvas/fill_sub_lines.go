package canvas

// How many lines a pixel row is cut into before its coverage is worked out, and
// with them the finest step along an edge a fill can describe. Sixteen lines
// put the error of a shallow edge at a fraction of a pixel, which is as fine as
// a screen can show; the coverage of a pixel is not spread out, it is gone, so
// going finer buys nothing and costs work in fours.
const fillSubLines = 16
