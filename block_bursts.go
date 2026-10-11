//go:build android

package antui

// blockBursts is how many of the device's own bursts go in one write.
//
// One is the smallest the device can take and leaves no room for a late
// goroutine; four is about ten milliseconds on a typical phone, which is
// under the threshold where a person hears the delay and far enough above
// the burst to survive a garbage collection.
const blockBursts = 4
