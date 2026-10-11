package svg

// maxTilePixels is how many pixels one tile may be drawn into. A tile the size
// of the drawing at the drawing's own scale is ordinary and worth every pixel
// of it; a tile written in fractions of a box on a shape drawn enormous is not
// worth running the machine out of memory for. What is worth it is what the
// drawing can reach: such a tile is drawn as the part of it the canvas can ask
// about, at the scale the pixels are sampled at, and comes out smaller and
// blurrier only where even that part is over the limit — the tiling is still
// the tiling, and a pattern is a texture rather than the thing itself.
const maxTilePixels = 2_000_000
