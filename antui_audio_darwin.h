#ifndef ANTUI_AUDIO_DARWIN_H
#define ANTUI_AUDIO_DARWIN_H

#include <AudioToolbox/AudioToolbox.h>
#include <string.h>
#include <stdlib.h>

// The Go side, called once for each buffer the queue hands back.
extern void antuiAudioFill(int id, float *samples, int frames);

OSStatus antuiAudioOpen(double rate, int id, int frames, int count, AudioQueueRef *out);
void antuiAudioClose(AudioQueueRef queue);

#endif
