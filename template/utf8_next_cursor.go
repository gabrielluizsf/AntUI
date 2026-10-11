package template

func utf8NextCursor(text string, index int) int {
	if index >= len(text) {
		return len(text)
	}
	index++
	for index < len(text) && text[index]&0xC0 == 0x80 {
		index++
	}
	return index
}
