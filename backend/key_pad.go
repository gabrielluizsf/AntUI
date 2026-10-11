package backend

// A game controller's buttons, named the way the platform names them rather
// than the way any one pad is printed — what is drawn on the button varies by
// manufacturer and the code does not.
const (
	KeyPadA Key = 256 + 40 + iota
	KeyPadB
	KeyPadX
	KeyPadY
	KeyPadL1
	KeyPadR1
	KeyPadL2
	KeyPadR2
	KeyPadStart
	KeyPadSelect
	KeyPadThumbL
	KeyPadThumbR
	KeyPadMode
)
