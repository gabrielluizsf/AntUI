//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// X keysyms this backend recognises by name.
const (
	symBackspace = 0xFF08
	symTab       = 0xFF09
	symReturn    = 0xFF0D
	symKPEnter   = 0xFF8D
	symEscape    = 0xFF1B
	symHome      = 0xFF50
	symLeft      = 0xFF51
	symUp        = 0xFF52
	symRight     = 0xFF53
	symDown      = 0xFF54
	symPageUp    = 0xFF55
	symPageDown  = 0xFF56
	symEnd       = 0xFF57
	symInsert    = 0xFF63
	symDelete    = 0xFFFF
	symISOLevel3 = 0xFF7E // the AltGr some layouts send instead
	symF1        = 0xFFBE
	symF12       = 0xFFC9
	symKP0       = 0xFFB0
	symKP9       = 0xFFB9

	symShiftL = 0xFFE1
	symShiftR = 0xFFE2
	symCtrlL  = 0xFFE3
	symCtrlR  = 0xFFE4
	symCaps   = 0xFFE5
	symAltL   = 0xFFE9
	symAltR   = 0xFFEA
	symSuperL = 0xFFEB
	symSuperR = 0xFFEC

	symDeadLo = 0xFE50 // the dead-accent block
	symDeadHi = 0xFE5F
)
