//go:build windows

package windows

func hiWord(v uintptr) int { return int(int16(uint16(v >> 16))) }
