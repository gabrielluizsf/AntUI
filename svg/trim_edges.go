package svg

import "strings"

// trimEdges takes the whitespace off the two ends of the writing, which is
// what the line breaks around a `<text>` in the file are there for — they make
// the file readable and say nothing about where the writing starts and ends.
// A piece left with nothing in it is dropped, and where it would have started
// goes onto the writing after it: the place was in the file and still counts.
//
// preserve is `xml:space="preserve"`, and where it is on there is nothing to
// take off: the writing is drawn the way the file wrote it, spaces at the two
// ends of it included, because that is what the drawing asked for.
func trimEdges(runs []TextRun, preserve bool) []TextRun {
	if len(runs) == 0 {
		return runs
	}
	if !preserve {
		// The places a drawing gave one per letter are kept in step with the
		// letters they belong to: the values the dropped spaces were sitting
		// on go with them, so the letter that is now first is still at the
		// value the list counted out for it.
		n := len([]rune(runs[0].Text))
		runs[0].Text = strings.TrimLeft(runs[0].Text, " ")
		runs[0].XList = dropFront(runs[0].XList, n-len([]rune(runs[0].Text)))
		runs[0].YList = dropFront(runs[0].YList, n-len([]rune(runs[0].Text)))
		last := len(runs) - 1
		n = len([]rune(runs[last].Text))
		runs[last].Text = strings.TrimRight(runs[last].Text, " ")
		keep := len([]rune(runs[last].Text))
		runs[last].XList = dropPast(runs[last].XList, keep, n)
		runs[last].YList = dropPast(runs[last].YList, keep, n)
	}
	if runs[0].Text == "" && len(runs) > 1 {
		head := runs[0]
		if !runs[1].HasX {
			runs[1].X, runs[1].HasX = head.X, head.HasX
		}
		if !runs[1].HasY {
			runs[1].Y, runs[1].HasY = head.Y, head.HasY
		}
		runs[1].DX += head.DX
		runs[1].DY += head.DY
	}
	out := runs[:0]
	for _, r := range runs {
		if r.Text != "" {
			out = append(out, r)
		}
	}
	return out
}
