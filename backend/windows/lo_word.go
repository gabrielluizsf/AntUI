//go:build windows

package windows

func loWord(v uintptr) int { return int(int16(uint16(v))) }
