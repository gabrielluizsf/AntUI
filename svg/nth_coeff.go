package svg

import "strconv"

// nthCoeff reads the coefficient written before an `n`: nothing means one, a
// bare sign is that sign, and otherwise it is the integer written.
func nthCoeff(s string) (int, bool) {
	switch s {
	case "", "+":
		return 1, true
	case "-":
		return -1, true
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return v, true
}
