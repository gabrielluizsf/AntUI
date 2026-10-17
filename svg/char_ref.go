package svg

import (
	"fmt"
	"strconv"
	"strings"
)

// charAt reads the character written as `&#…;` at the `&` of s[i] — the
// number in decimal, or in hexadecimal behind an `x` — and answers the
// character it is and how many bytes it took.
func charAt(s string, i int) (string, int, error) {
	end := strings.IndexByte(s[i+2:], ';')
	if end < 0 {
		return "", 0, fmt.Errorf("the number after a `&#` is never closed with a `;`")
	}
	num := s[i+2 : i+2+end]
	r, ok := parseCharRef(num)
	if !ok {
		return "", 0, fmt.Errorf("`&#%s;` is not a character XML allows", num)
	}
	return string(r), end + 3, nil
}

// parseCharRef reads the number of a character reference and answers the
// character it is, or says the number is not one a document may hold.
func parseCharRef(num string) (rune, bool) {
	base := 10
	if strings.HasPrefix(num, "x") {
		base, num = 16, num[1:]
	}
	if num == "" {
		return 0, false
	}
	n, err := strconv.ParseUint(num, base, 32)
	if err != nil {
		return 0, false
	}
	r := rune(n)
	return r, validXMLChar(r)
}

// validXMLChar says a character may appear in a document at all: XML keeps the
// control characters and the unpaired surrogates out of the text a reader is
// handed, so a reference to one of them is an error rather than the character
// it names.
func validXMLChar(r rune) bool {
	return r == 0x9 || r == 0xA || r == 0xD ||
		0x20 <= r && r <= 0xD7FF ||
		0xE000 <= r && r <= 0xFFFD ||
		0x10000 <= r && r <= 0x10FFFF
}
