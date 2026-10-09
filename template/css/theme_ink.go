package css

import "github.com/gabrielluizsf/antui/canvas"

// ThemeInk is what the engine stores in Style.Color when the cascade resolves
// colour to its CSS initial value — a `color: initial` (or revert/unset-less
// body) means the theme's ink, not a literal pixel. Templates must translate
// this sentinel into the window's theme text colour; the sentinel itself is an
// almost-invisible black that must never be produced by parseColor.
const ThemeInk = canvas.Color(0xFF000001)
