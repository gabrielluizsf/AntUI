package css

import "github.com/gabrielluizsf/antui/canvas"

// CurrentColor is the sentinel a style field holds while its declaration said
// currentColor. The cascade resolves it to the computed color property once
// every declaration has applied, because a browser reads currentColor after
// the whole cascade rather than mid-way. The value happens to spell
// "transparent black" if a sheet literally writes rgba(0 0 0 / 0); the engine
// accepts that quirk.
const CurrentColor canvas.Color = 0x00000001
