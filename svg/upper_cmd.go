package svg

// upperCmd is the command without its case, since the case only says whether
// the numbers are measured from where the line is.
func upperCmd(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}
