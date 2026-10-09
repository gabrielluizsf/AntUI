package css

// finishGridLine resolves a line that carried a span: a span by itself
// becomes an index-count line, a named line carries the count as its span.
func finishGridLine(line GridLine, count int) (GridLine, bool) {
	if count == 0 {
		count = 1
	}
	if line.Kind == GridLineAuto {
		return GridLine{Kind: GridLineSpan, Index: count}, true
	}
	line.Span = count
	return line, true
}
