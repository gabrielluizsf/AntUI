package template

import "github.com/gabrielluizsf/antui/canvas"

// cyberPalette is the whole range of the futuristic look.
var cyberPalette = struct {
	bg       canvas.Color
	panel    canvas.Color
	panelHov canvas.Color
	panelAct canvas.Color
	edge     canvas.Color
	neon     canvas.Color
	text     canvas.Color
	textDim  canvas.Color
}{
	bg:       0xFF0B0E1A,
	panel:    0xFF12182B,
	panelHov: 0xFF17213D,
	panelAct: 0xFF0D1222,
	edge:     0xFF24365C,
	neon:     0xFF3DE8FF,
	text:     0xFFE8F0FF,
	textDim:  0xFF8FA2C0,
}
