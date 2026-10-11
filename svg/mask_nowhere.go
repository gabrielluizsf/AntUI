package svg

// maskNowhere is the mask a reference the drawing cannot follow stands for: a
// picture with nothing in it, which keeps nothing, the same as a clip that
// names a clipPath that is not there. It is shared because it holds nothing
// that belongs to any one drawing.
var maskNowhere = &maskDef{}
