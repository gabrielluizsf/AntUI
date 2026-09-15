//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <aaudio/AAudio.h>
#include <dlfcn.h>
#include <stdlib.h>

// AAudio arrived in API 26, and libaaudio.so is not in the stub libraries a
// build against API 21 links against — so it cannot be linked at all, only
// opened. Every call goes through a pointer looked up once.
//
// The alternative below API 26 is OpenSL ES, which Google has itself
// deprecated in favour of this. See the note in antui/audio_android.go.
typedef struct {
	void *lib;
	aaudio_result_t (*createBuilder)(AAudioStreamBuilder **);
	void (*setSampleRate)(AAudioStreamBuilder *, int32_t);
	void (*setChannelCount)(AAudioStreamBuilder *, int32_t);
	void (*setFormat)(AAudioStreamBuilder *, int32_t);
	void (*setDirection)(AAudioStreamBuilder *, int32_t);
	void (*setPerformanceMode)(AAudioStreamBuilder *, int32_t);
	void (*setBufferCapacity)(AAudioStreamBuilder *, int32_t);
	aaudio_result_t (*openStream)(AAudioStreamBuilder *, AAudioStream **);
	aaudio_result_t (*deleteBuilder)(AAudioStreamBuilder *);
	aaudio_result_t (*start)(AAudioStream *);
	aaudio_result_t (*stop)(AAudioStream *);
	aaudio_result_t (*closeStream)(AAudioStream *);
	aaudio_result_t (*write)(AAudioStream *, const void *, int32_t, int64_t);
	aaudio_result_t (*read)(AAudioStream *, void *, int32_t, int64_t);
	int32_t (*getSampleRate)(AAudioStream *);
	int32_t (*getChannelCount)(AAudioStream *);
	int32_t (*getFramesPerBurst)(AAudioStream *);
	int32_t (*getBufferSize)(AAudioStream *);
	aaudio_result_t (*setBufferSize)(AAudioStream *, int32_t);
	const char *(*resultText)(aaudio_result_t);
} mw_aaudio;

static mw_aaudio A;
static int mw_aaudio_state; // 0 not tried, 1 there, -1 not

#define MW_SYM(field, name)                                  \
	A.field = dlsym(A.lib, name);                            \
	if (A.field == NULL) { dlclose(A.lib); A.lib = NULL;     \
		mw_aaudio_state = -1; return -1; }

static int mw_aaudio_load(void) {
	if (mw_aaudio_state != 0) {
		return mw_aaudio_state;
	}
	A.lib = dlopen("libaaudio.so", RTLD_NOW);
	if (A.lib == NULL) {
		mw_aaudio_state = -1;
		return -1;
	}
	MW_SYM(createBuilder, "AAudio_createStreamBuilder")
	MW_SYM(setSampleRate, "AAudioStreamBuilder_setSampleRate")
	MW_SYM(setChannelCount, "AAudioStreamBuilder_setChannelCount")
	MW_SYM(setFormat, "AAudioStreamBuilder_setFormat")
	MW_SYM(setDirection, "AAudioStreamBuilder_setDirection")
	MW_SYM(setPerformanceMode, "AAudioStreamBuilder_setPerformanceMode")
	MW_SYM(setBufferCapacity, "AAudioStreamBuilder_setBufferCapacityInFrames")
	MW_SYM(openStream, "AAudioStreamBuilder_openStream")
	MW_SYM(deleteBuilder, "AAudioStreamBuilder_delete")
	MW_SYM(start, "AAudioStream_requestStart")
	MW_SYM(stop, "AAudioStream_requestStop")
	MW_SYM(closeStream, "AAudioStream_close")
	MW_SYM(write, "AAudioStream_write")
	MW_SYM(read, "AAudioStream_read")
	MW_SYM(getSampleRate, "AAudioStream_getSampleRate")
	MW_SYM(getChannelCount, "AAudioStream_getChannelCount")
	MW_SYM(getFramesPerBurst, "AAudioStream_getFramesPerBurst")
	MW_SYM(getBufferSize, "AAudioStream_getBufferSizeInFrames")
	MW_SYM(setBufferSize, "AAudioStream_setBufferSizeInFrames")
	MW_SYM(resultText, "AAudio_convertResultToText")
	mw_aaudio_state = 1;
	return 1;
}

// mw_audio_open builds and starts a stream, and reports what the device
// actually settled on.
static aaudio_result_t mw_audio_open(int32_t rate, int32_t channels, int32_t capacity,
                                     int32_t input, AAudioStream **out,
                                     int32_t *gotRate, int32_t *gotChannels,
                                     int32_t *burst) {
	AAudioStreamBuilder *b = NULL;
	aaudio_result_t r = A.createBuilder(&b);
	if (r != AAUDIO_OK) {
		return r;
	}
	A.setDirection(b, input ? AAUDIO_DIRECTION_INPUT : AAUDIO_DIRECTION_OUTPUT);
	A.setFormat(b, AAUDIO_FORMAT_PCM_FLOAT);
	A.setSampleRate(b, rate);
	A.setChannelCount(b, channels);
	A.setPerformanceMode(b, AAUDIO_PERFORMANCE_MODE_LOW_LATENCY);
	if (capacity > 0) {
		A.setBufferCapacity(b, capacity);
	}
	r = A.openStream(b, out);
	A.deleteBuilder(b);
	if (r != AAUDIO_OK) {
		return r;
	}
	*gotRate = A.getSampleRate(*out);
	*gotChannels = A.getChannelCount(*out);
	*burst = A.getFramesPerBurst(*out);
	return AAUDIO_OK;
}

static aaudio_result_t mw_audio_start(AAudioStream *s)  { return A.start(s); }
static aaudio_result_t mw_audio_stop(AAudioStream *s)   { return A.stop(s); }
static aaudio_result_t mw_audio_close(AAudioStream *s)  { return A.closeStream(s); }
static int32_t mw_audio_buffer(AAudioStream *s)         { return A.getBufferSize(s); }
static aaudio_result_t mw_audio_set_buffer(AAudioStream *s, int32_t n) {
	return A.setBufferSize(s, n);
}
static aaudio_result_t mw_audio_write(AAudioStream *s, const float *data,
                                      int32_t frames, int64_t timeout) {
	return A.write(s, data, frames, timeout);
}
static aaudio_result_t mw_audio_read(AAudioStream *s, float *data,
                                     int32_t frames, int64_t timeout) {
	return A.read(s, data, frames, timeout);
}
static const char *mw_audio_text(aaudio_result_t r) { return A.resultText(r); }
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

// ErrNoAAudio is a device too old for AAudio, which is anything below
// Android 8.
var ErrNoAAudio = errors.New("ndk: this device has no AAudio; it arrived in API 26")

// AudioAvailable reports whether AAudio can be used at all. It opens the
// library the first time and remembers.
func AudioAvailable() bool { return C.mw_aaudio_load() == 1 }

// AudioStream is one stream of sound going out.
//
// It is written to rather than pulling: a goroutine calls [AudioStream.Write]
// and the call blocks until there is room. AAudio also has a callback mode,
// which is lower latency and runs on a real-time thread — a thread Go's
// collector may stop, which is exactly what a real-time thread must not
// have happen to it.
type AudioStream struct {
	ptr *C.AAudioStream
	// Rate, Channels and Burst are what the device settled on, which is not
	// always what was asked for.
	Rate     int
	Channels int
	// Burst is how many frames the device moves at a time. Writing in
	// multiples of it is what keeps a stream from stuttering; writing less
	// than one is a write that cannot be satisfied without waiting.
	Burst int
}

// OpenAudio starts a stream going out to the speakers. capacity is how many
// frames the device should buffer, and 0 lets it choose.
func OpenAudio(rate, channels, capacity int) (*AudioStream, error) {
	return openAudio(rate, channels, capacity, false)
}

// OpenAudioInput starts a stream coming in from the microphone.
//
// It needs the RECORD_AUDIO permission, and without it the open fails rather
// than returning silence — which is the right way round, and not what every
// platform does.
func OpenAudioInput(rate, channels, capacity int) (*AudioStream, error) {
	return openAudio(rate, channels, capacity, true)
}

func openAudio(rate, channels, capacity int, input bool) (*AudioStream, error) {
	if C.mw_aaudio_load() != 1 {
		return nil, ErrNoAAudio
	}
	var in C.int32_t
	if input {
		in = 1
	}
	var (
		stream      *C.AAudioStream
		gotRate     C.int32_t
		gotChannels C.int32_t
		burst       C.int32_t
	)
	r := C.mw_audio_open(C.int32_t(rate), C.int32_t(channels), C.int32_t(capacity),
		in, &stream, &gotRate, &gotChannels, &burst)
	if r != C.AAUDIO_OK {
		return nil, audioError("opening the stream", r)
	}
	s := &AudioStream{
		ptr:      stream,
		Rate:     int(gotRate),
		Channels: int(gotChannels),
		Burst:    int(burst),
	}
	// Two bursts is the smallest buffer that does not underrun on the first
	// hiccup, and is what every AAudio example starts from.
	if s.Burst > 0 {
		C.mw_audio_set_buffer(stream, C.int32_t(s.Burst*2))
	}
	if r := C.mw_audio_start(stream); r != C.AAUDIO_OK {
		C.mw_audio_close(stream)
		return nil, audioError("starting the stream", r)
	}
	return s, nil
}

// Write sends frames and waits for room. data is interleaved: one float per
// channel per frame.
func (s *AudioStream) Write(data []float32, timeout time.Duration) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("ndk: the audio stream is closed")
	}
	frames := len(data) / s.Channels
	if frames == 0 {
		return 0, nil
	}
	r := C.mw_audio_write(s.ptr, (*C.float)(unsafe.Pointer(&data[0])),
		C.int32_t(frames), C.int64_t(timeout.Nanoseconds()))
	if r < 0 {
		return 0, audioError("writing", C.aaudio_result_t(r))
	}
	return int(r), nil
}

// Read takes frames from the microphone and waits for them. data is filled
// interleaved, one float per channel per frame, and the count of frames read
// comes back — which may be fewer than asked for.
func (s *AudioStream) Read(data []float32, timeout time.Duration) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("ndk: the audio stream is closed")
	}
	frames := len(data) / s.Channels
	if frames == 0 {
		return 0, nil
	}
	r := C.mw_audio_read(s.ptr, (*C.float)(unsafe.Pointer(&data[0])),
		C.int32_t(frames), C.int64_t(timeout.Nanoseconds()))
	if r < 0 {
		return 0, audioError("reading", C.aaudio_result_t(r))
	}
	return int(r), nil
}

// BufferFrames is how much the device is holding, which is what latency is
// made of.
func (s *AudioStream) BufferFrames() int {
	if s.ptr == nil {
		return 0
	}
	return int(C.mw_audio_buffer(s.ptr))
}

// Latency is roughly how long it is between writing a sample and hearing it:
// the buffer over the rate. It is what the device was talked into rather
// than a measurement of the whole path, which no API reports.
func (s *AudioStream) Latency() time.Duration {
	if s.Rate <= 0 {
		return 0
	}
	return time.Duration(s.BufferFrames()) * time.Second / time.Duration(s.Rate)
}

// Close stops the stream and gives it back.
func (s *AudioStream) Close() error {
	if s.ptr == nil {
		return nil
	}
	C.mw_audio_stop(s.ptr)
	r := C.mw_audio_close(s.ptr)
	s.ptr = nil
	if r != C.AAUDIO_OK {
		return audioError("closing", r)
	}
	return nil
}

// audioError turns AAudio's numbered result into its own words, which are
// better than anything this package would invent.
func audioError(what string, r C.aaudio_result_t) error {
	text := C.GoString(C.mw_audio_text(r))
	if text == "" {
		return errors.New("ndk: audio: " + what + " failed")
	}
	return errors.New("ndk: audio: " + what + ": " + text)
}
