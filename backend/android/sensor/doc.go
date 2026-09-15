// Package sensor reads the device's sensors.
//
// It goes through the NDK rather than through Java: Android exposes the whole
// sensor API to native code, so nothing here crosses the bridge except
// looking up the app's own package name once.
//
// The one thing worth reading twice is [Reading.ForDisplay]. A sensor reports
// in a frame fixed to the hardware, and the screen rotates without it.
package sensor
