//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// writeChunk bounds how much goes out between two reads. Large enough that a
// frame is a handful of writes, small enough that the server never waits long
// for us to listen.
const writeChunk = 64 << 10
