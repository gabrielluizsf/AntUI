package svg

import "fmt"

// readParen reads the argument a functional selector is given. i is the index
// of the `(`, and the answer is what stands between it and its matching `)`,
// with the index just past that.
func readParen(s string, i int) (string, int, error) {
	end := parenEnd(s, i)
	if end < 0 {
		return "", i, fmt.Errorf("a parenthesised selector argument is never closed")
	}
	return s[i+1 : end], end + 1, nil
}
