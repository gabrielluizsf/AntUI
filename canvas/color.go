package canvas

// Color is a pixel in 0xAARRGGBB order: alpha in the top eight bits, then
// red, green and blue. It is the one colour representation the whole library
// speaks, on every surface format — the narrow formats convert on the way in
// and on the way out, so callers never see a 565 word or a palette index.
type Color uint32
