package css

// mathParser is a small recursive-descent evaluator over + - * / ( ) ,
// numbers-with-units and the four math functions.
type mathParser struct {
	t   *mathTokens
	ctx Units
}
