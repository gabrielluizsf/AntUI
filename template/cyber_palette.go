package template

import "github.com/gabrielluizsf/antui/canvas"

// cyberPalette is the whole range of the futuristic look: the Hydra skin, a
// single blue accent on near-black panels, no chromatic split.
var cyberPalette = struct {
	bg       canvas.Color
	panel    canvas.Color
	panelHov canvas.Color
	panelAct canvas.Color
	field    canvas.Color
	edge     canvas.Color
	neon     canvas.Color
	text     canvas.Color
	textDim  canvas.Color
}{
	bg:       0xFF04070F,
	panel:    0xFF050B19,
	panelHov: 0xFF0B1530,
	panelAct: 0xFF030710,
	field:    0xFF020306,
	edge:     0xFF224B6A,
	neon:     0xFF57C1FF,
	text:     0xFFFFFFFF,
	textDim:  0xFF8FB4D9,
}
