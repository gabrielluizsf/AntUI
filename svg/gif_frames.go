package svg

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
)

// gifFrames is a GIF read into the frames it animates with, each in the state
// it is on screen: the file stores only the change one frame makes to the one
// before, and the disposal byte says what of it is left for the next — piling
// the changes onto one sheet in order is what turns them into pictures.
func gifFrames(data []byte) (*pictureFrames, error) {
	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("it is not a GIF: %w", err)
	}
	if g.Config.Width <= 0 || g.Config.Height <= 0 || len(g.Image) == 0 {
		return nil, fmt.Errorf("it has no size to it")
	}
	out := &pictureFrames{show: make([]float64, len(g.Image))}
	sheet := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))
	var undo *image.RGBA
	for i, f := range g.Image {
		if undo != nil {
			// The frame before asked for what was there before it.
			draw.Draw(sheet, sheet.Bounds(), undo, image.Point{}, draw.Src)
			undo = nil
		}
		disposal := byte(gif.DisposalNone)
		if i < len(g.Disposal) {
			disposal = g.Disposal[i]
		}
		if disposal == gif.DisposalPrevious {
			undo = image.NewRGBA(sheet.Bounds())
			copy(undo.Pix, sheet.Pix)
		}
		draw.Draw(sheet, f.Bounds(), f, f.Bounds().Min, draw.Over)
		cv, err := pictureCanvas(sheet)
		if err != nil {
			return nil, err
		}
		out.frames = append(out.frames, cv)
		if i < len(g.Delay) {
			// shortcut: a loop of zero delays stays on the first frame, while
			// browsers give a zero a hundred milliseconds; upgrade when one of
			// those GIFs has to move.
			out.show[i] = float64(g.Delay[i]) / 100 // the file counts in hundredths
		}
		// shortcut: the spec restores the background colour here, while a GIF
		// written with transparency means what stands behind it — the rect goes
		// to nothing; upgrade if one painted behind shows that square wrong.
		if disposal == gif.DisposalBackground {
			draw.Draw(sheet, f.Bounds(), image.Transparent, image.Point{}, draw.Src)
		}
	}
	return out, nil
}
