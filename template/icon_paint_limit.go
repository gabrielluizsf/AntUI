package template

// iconPaintLimit is how many painted icons are kept before the oldest are thrown
// out by starting over. An interface has a handful of icons at a handful of
// sizes in a handful of states, so the map is only ever as big as the set of
// things being drawn — and the bound is what keeps a program that names icons
// by something that changes per frame from growing it without end.
const iconPaintLimit = 64
