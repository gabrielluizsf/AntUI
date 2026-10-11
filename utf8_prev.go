package antui

// utf8Prev and utf8Next step the cursor by whole characters.
func utf8Prev(text string, index int) int {
	if index <= 0 {
		return 0
	}
	index--
	for index > 0 && text[index]&0xC0 == 0x80 {
		index--
	}
	return index
}
