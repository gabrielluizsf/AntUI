package css

// maxGridTracks caps how many tracks one axis may end up with, so a
// pathological repeat() cannot make the solver walk a million entries.
const maxGridTracks = 1000

// GridTemplate is a parsed grid-template value: one segment per track, in
// declaration order, each carrying the line names standing on the lines to its
// left and right. repeat(n, …) is already expanded here, where the count is
// written down; repeat(auto-fill | auto-fit, …) stays a single segment and
// [GridTemplate.Resolve] decides how many times it repeats.
type GridTemplate []GridSegment
