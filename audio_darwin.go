//go:build darwin

package antui

/*
#cgo CFLAGS: -Wno-deprecated-declarations
#cgo LDFLAGS: -framework AudioToolbox -framework CoreFoundation
#include <AudioToolbox/AudioToolbox.h>
#include <string.h>
#include <stdlib.h>

// The Go side, called once for each buffer the queue hands back.
extern void antuiAudioFill(int id, float *samples, int frames);

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
static OSStatus antuiAudioOpen(double rate, int id, int frames, int count, AudioQueueRef *out)
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

static void antuiAudioClose(AudioQueueRef queue)
{
    if (queue != NULL) AudioQueueDispose(queue, true);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// CoreAudio through an audio queue, which is the part of it that takes
// samples in the format they are already in and deals with the hardware's own
// ideas about rate and layout on the way past.
//
// The lower AudioUnit is what a professional tool reaches for; the queue is a
// tenth of the code and its extra few milliseconds are under what a screen
// refresh costs anyway.
//
// NOT YET COMPILED ON macOS. Like the Cocoa half of the window backend, this
// is written from the framework headers and waits for a machine to run it —
// see audio_darwin_test.go, which only runs there.

// The queues that are open, by a number rather than by a pointer: a Go
// pointer cannot be handed to C and kept there, and this is the one thing C
// has to hold on to.
var (
	queuesMu  sync.Mutex
	queues    = map[int]*coreaudio{}
	nextQueue int
)

type coreaudio struct {
	audio *Audio
	id    int
	queue C.AudioQueueRef
	delay time.Duration
}

func (c *coreaudio) latency() time.Duration { return c.delay }

func (c *coreaudio) close() error {
	queuesMu.Lock()
	delete(queues, c.id)
	queuesMu.Unlock()

	// Disposing waits for the callback to finish, so after this nothing is
	// looking the stream up any more.
	C.antuiAudioClose(c.queue)
	c.queue = nil
	return nil
}

// openAudio is the macOS half of OpenAudio.
func openAudio(name string, a *Audio) (driver, error) {
	// Ten milliseconds a buffer and three of them: a third of a frame at
	// sixty, and short enough that a sound lands with the thing that made it.
	const buffers = 3
	frames := max(a.rate/100, 64)

	c := &coreaudio{audio: a}

	queuesMu.Lock()
	nextQueue++
	c.id = nextQueue
	queues[c.id] = c
	queuesMu.Unlock()

	var queue C.AudioQueueRef
	status := C.antuiAudioOpen(C.double(a.rate), C.int(c.id), C.int(frames),
		C.int(buffers), &queue)
	if status != 0 {
		queuesMu.Lock()
		delete(queues, c.id)
		queuesMu.Unlock()
		return nil, fmt.Errorf("%w: CoreAudio refused a stream (%d)", ErrNoAudio, int(status))
	}

	c.queue = queue
	c.delay = time.Duration(frames*buffers) * time.Second / time.Duration(a.rate)
	return c, nil
}

//export antuiAudioFill
func antuiAudioFill(id C.int, samples *C.float, frames C.int) {
	// The buffer belongs to CoreAudio, so it is written through rather than
	// copied into: a slice over C memory, which is exactly the shape the
	// mixer writes.
	fillQueue(int(id), unsafe.Slice((*float32)(unsafe.Pointer(samples)), int(frames)*2))
}

// fillQueue is the Go half of the callback, with the C types left at the door
// so that the part which can be got wrong is the part a test can call.
func fillQueue(id int, block []float32) {
	queuesMu.Lock()
	c := queues[id]
	queuesMu.Unlock()

	if c == nil {
		// The stream was closed between the queue asking and this answering.
		for i := range block {
			block[i] = 0
		}
		return
	}
	c.audio.pull(block)
}