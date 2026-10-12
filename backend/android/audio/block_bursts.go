//go:build android

package audio

// blockBursts is how many of the device's own bursts go in one read, for the
// same reason it is what it is on the way out.
const blockBursts = 4
