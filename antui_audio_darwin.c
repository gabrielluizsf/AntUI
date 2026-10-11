#include "antui_audio_darwin.h"

// One buffer's worth: fill it and give it straight back. An audio queue holds
// several of these in a ring, so while this one is being played the next is
// already being asked for.
static void antuiAudioCallback(void *user, AudioQueueRef queue, AudioQueueBufferRef buffer)
{
    UInt32 bytes = buffer->mAudioDataBytesCapacity;
    buffer->mAudioDataByteSize = bytes;
    antuiAudioFill((int)(intptr_t)user, (float *)buffer->mAudioData,
                   (int)(bytes / (2 * sizeof(float))));
    AudioQueueEnqueueBuffer(queue, buffer, 0, NULL);
}

// antuiAudioOpen starts a stereo float stream at a rate and hands back the
// queue. The buffers are primed here — a queue started with nothing in it
// stops again immediately.
OSStatus antuiAudioOpen(double rate, int id, int frames, int count, AudioQueueRef *out)
{
    AudioStreamBasicDescription format;
    AudioQueueRef queue = NULL;
    OSStatus status;
    int i;

    memset(&format, 0, sizeof format);
    format.mSampleRate = rate;
    format.mFormatID = kAudioFormatLinearPCM;
    format.mFormatFlags = kAudioFormatFlagIsFloat | kAudioFormatFlagIsPacked;
    format.mFramesPerPacket = 1;
    format.mChannelsPerFrame = 2;
    format.mBitsPerChannel = 32;
    format.mBytesPerFrame = 2 * sizeof(float);
    format.mBytesPerPacket = format.mBytesPerFrame;

    status = AudioQueueNewOutput(&format, antuiAudioCallback, (void *)(intptr_t)id,
                                 NULL, NULL, 0, &queue);
    if (status != noErr) return status;

    for (i = 0; i < count; ++i) {
        AudioQueueBufferRef buffer = NULL;
        status = AudioQueueAllocateBuffer(queue, (UInt32)(frames * 2 * sizeof(float)), &buffer);
        if (status != noErr) {
            AudioQueueDispose(queue, true);
            return status;
        }
        // Filled here rather than left empty: the first thing anyone hears
        // should be the game, not the click of an empty buffer.
        antuiAudioCallback((void *)(intptr_t)id, queue, buffer);
    }

    status = AudioQueueStart(queue, NULL);
    if (status != noErr) {
        AudioQueueDispose(queue, true);
        return status;
    }
    *out = queue;
    return noErr;
}

void antuiAudioClose(AudioQueueRef queue)
{
    if (queue != NULL) AudioQueueDispose(queue, true);
}
