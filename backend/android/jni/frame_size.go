//go:build android

package jni

// frameSize is how many local references a call may make before it has to
// manage them itself. The default table is 512 and overflowing it aborts the
// process, so a frame is opened for every Do; 32 is enough for anything that
// is not looping over an array, and the ones that do open their own.
const frameSize = 32
