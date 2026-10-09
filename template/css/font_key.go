package css

// fontKey is one question the sheet answers: a font-family list at a weight
// and a slant.
type fontKey struct {
	family  string
	weight  uint16
	slanted bool
}
