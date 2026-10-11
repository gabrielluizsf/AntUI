//go:build windows

package windows

// GMEM_MOVEABLE. The clipboard owns what it is given and moves it about, so
// it will not take a fixed block.
const gmemMoveable = 0x0002
