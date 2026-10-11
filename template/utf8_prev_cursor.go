package template

// utf8PrevCursor and utf8NextCursor step a cursor by whole characters.
func utf8PrevCursor(text string, index int) int {
	if index <= 0 {
		return 0
	}
	index--
	for index > 0 && text[index]&0xC0 == 0x80 {
		index--
	}
	return index
}
