// Package audio is the part of sound that is Android's rather than AntUI's:
// asking the system for the right to be heard, noticing the headphones
// coming out, and listening through the microphone.
//
// Playing goes through [antui.OpenAudio] like everywhere else — on Android it
// lands on AAudio. This package is what an app should use *around* that.
package audio
