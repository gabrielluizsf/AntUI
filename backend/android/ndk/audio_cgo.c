//go:build android


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

int mw_aaudio_load(void) {
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
aaudio_result_t mw_audio_open(int32_t rate, int32_t channels, int32_t capacity,
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

aaudio_result_t mw_audio_start(AAudioStream *s)  { return A.start(s); }
aaudio_result_t mw_audio_stop(AAudioStream *s)   { return A.stop(s); }
aaudio_result_t mw_audio_close(AAudioStream *s)  { return A.closeStream(s); }
int32_t mw_audio_buffer(AAudioStream *s)         { return A.getBufferSize(s); }
aaudio_result_t mw_audio_set_buffer(AAudioStream *s, int32_t n) {
	return A.setBufferSize(s, n);
}
aaudio_result_t mw_audio_write(AAudioStream *s, const float *data,
                                      int32_t frames, int64_t timeout) {
	return A.write(s, data, frames, timeout);
}
aaudio_result_t mw_audio_read(AAudioStream *s, float *data,
                                     int32_t frames, int64_t timeout) {
	return A.read(s, data, frames, timeout);
}
const char *mw_audio_text(aaudio_result_t r) { return A.resultText(r); }
