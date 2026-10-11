package svg

import "github.com/gabrielluizsf/antui/canvas"

// walkRuns carries every run of writing to where the drawing puts it. Measuring
// and painting both take this walk, so that the box a fill is measured against
// and the place the letters land cannot come out different.
func walkRuns(runs []TextRun, m canvas.Matrix) []piece {
	pieces := make([]piece, 0, len(runs))
	penX, penY := 0.0, 0.0
	for _, run := range runs {
		st := run.Style
		if run.HasX {
			penX = run.X
		}
		if run.HasY {
			penY = run.Y
		}
		penX, penY = penX+run.DX, penY+run.DY
		if st.Hidden || run.Text == "" {
			continue
		}
		at := m.Mul(st.Transform)
		scale := scaleOf(at)
		if scale <= 0 {
			continue
		}
		size := st.FontSize
		if size <= 0 {
			size = defaultFontSize
		}
		// The face is asked for at the size in pixels this writing covers on
		// the canvas, so a drawing painted big is written big rather than a
		// small piece of writing blown up afterwards. The face itself is the
		// one the canvas draws with: a font-family list is read, but a drawing
		// carries no files of its own, so the names run into the face the
		// program already had — see [Style].
		face := canvas.DefaultFace().AtSize(size * scale)
		width := (float64(face.Width(run.Text)) +
			st.LetterSpacing*float64(spaced(st, run.Text))) / scale
		// An anchor hangs the writing on the pen, and a piece hung on an
		// outline is hung along the outline instead — the pen has no part in
		// where it lands, and does not carry the pen on afterwards either.
		along := run.Path != nil
		if !along {
			switch st.Anchor {
			case AnchorMiddle:
				penX -= width / 2
			case AnchorEnd:
				penX -= width
			}
		}
		pieces = append(pieces, piece{
			run: run, face: face, scale: scale, width: width,
			penX: penX, penY: penY, at: at,
		})
		if !along {
			penX += width
		}
	}
	return pieces
}
