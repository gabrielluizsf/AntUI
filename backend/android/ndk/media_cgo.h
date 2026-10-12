//go:build android

// The media pump, shared by every file in this package that drives a
// codec: the loop is C because it runs per frame and per audio block
// and does nothing a caller would want to see.
#include <media/NdkMediaExtractor.h>
#include <media/NdkMediaCodec.h>
#include <media/NdkMediaFormat.h>
#include <android/native_window.h>
#include <stdlib.h>
#include <string.h>

#include "imagereader.h"

// The extractor, the codec and the format are all API 21, so unlike the
// camera and the image reader they are linked rather than looked up.
//
// A file is pulled apart by the extractor into compressed samples, fed to a
// codec, and comes out the other side either as pixels on a surface or as
// PCM in a buffer. The pump below is that loop, in C because it runs per
// frame and per audio block and does nothing a caller would want to see.

typedef struct {
	AMediaExtractor *ex;
	AMediaCodec *codec;
	AImageReader *reader;   // video only
	ANativeWindow *window;  // video only
	int64_t pts;            // the presentation time of the last output
	int eof;                // the extractor has no more samples
	int drained;            // the codec has produced its last output
} mw_media;

// What the pump reports.
enum {
	MW_MEDIA_AGAIN = 0, // nothing yet; call again
	MW_MEDIA_FRAME = 1, // an output was produced
	MW_MEDIA_END   = 2, // the stream is over
	MW_MEDIA_ERROR = -1,
};

// mw_media_feed hands the codec one compressed sample, if it wants one.
static void mw_media_feed(mw_media *m) {
	if (m->eof) {
		return;
	}
	ssize_t in = AMediaCodec_dequeueInputBuffer(m->codec, 0);
	if (in < 0) {
		return;
	}
	size_t capacity = 0;
	uint8_t *buf = AMediaCodec_getInputBuffer(m->codec, in, &capacity);
	if (buf == NULL) {
		return;
	}
	ssize_t n = AMediaExtractor_readSampleData(m->ex, buf, capacity);
	if (n < 0) {
		// The end. The codec is told so with an empty buffer carrying the
		// flag, and goes on producing until it has emptied itself.
		AMediaCodec_queueInputBuffer(m->codec, in, 0, 0, 0,
		                             AMEDIACODEC_BUFFER_FLAG_END_OF_STREAM);
		m->eof = 1;
		return;
	}
	int64_t time = AMediaExtractor_getSampleTime(m->ex);
	AMediaCodec_queueInputBuffer(m->codec, in, 0, (size_t)n, (uint64_t)time, 0);
	AMediaExtractor_advance(m->ex);
}

// mw_video_pump feeds the codec and sends whatever came out to the surface,
// which is where the image reader picks it up.
static int mw_video_pump(mw_media *m) {
	mw_media_feed(m);
	if (m->drained) {
		return MW_MEDIA_END;
	}
	AMediaCodecBufferInfo info;
	memset(&info, 0, sizeof(info));
	ssize_t out = AMediaCodec_dequeueOutputBuffer(m->codec, &info, 0);
	if (out == AMEDIACODEC_INFO_TRY_AGAIN_LATER ||
	    out == AMEDIACODEC_INFO_OUTPUT_FORMAT_CHANGED ||
	    out == AMEDIACODEC_INFO_OUTPUT_BUFFERS_CHANGED) {
		return MW_MEDIA_AGAIN;
	}
	if (out < 0) {
		return MW_MEDIA_ERROR;
	}
	m->pts = info.presentationTimeUs;
	if (info.flags & AMEDIACODEC_BUFFER_FLAG_END_OF_STREAM) {
		m->drained = 1;
	}
	// Rendering is what puts it on the surface. Releasing without it throws
	// the frame away, which is how a decoder is made to run at full speed
	// without drawing anything.
	AMediaCodec_releaseOutputBuffer(m->codec, out, true);
	return m->drained ? MW_MEDIA_END : MW_MEDIA_FRAME;
}

// mw_audio_pump feeds the codec and copies whatever came out into dst as
// floats. PCM from a decoder is signed 16-bit; anything else is refused
// rather than guessed at.
static int mw_audio_pump(mw_media *m, float *dst, int cap, int *wrote) {
	*wrote = 0;
	mw_media_feed(m);
	if (m->drained) {
		return MW_MEDIA_END;
	}
	AMediaCodecBufferInfo info;
	memset(&info, 0, sizeof(info));
	ssize_t out = AMediaCodec_dequeueOutputBuffer(m->codec, &info, 0);
	if (out == AMEDIACODEC_INFO_TRY_AGAIN_LATER ||
	    out == AMEDIACODEC_INFO_OUTPUT_FORMAT_CHANGED ||
	    out == AMEDIACODEC_INFO_OUTPUT_BUFFERS_CHANGED) {
		return MW_MEDIA_AGAIN;
	}
	if (out < 0) {
		return MW_MEDIA_ERROR;
	}
	m->pts = info.presentationTimeUs;
	size_t size = 0;
	uint8_t *buf = AMediaCodec_getOutputBuffer(m->codec, out, &size);
	if (buf != NULL && info.size > 0) {
		const int16_t *pcm = (const int16_t *)(buf + info.offset);
		int samples = info.size / 2;
		if (samples > cap) {
			samples = cap;
		}
		for (int i = 0; i < samples; i++) {
			dst[i] = (float)pcm[i] / 32768.0f;
		}
		*wrote = samples;
	}
	if (info.flags & AMEDIACODEC_BUFFER_FLAG_END_OF_STREAM) {
		m->drained = 1;
	}
	AMediaCodec_releaseOutputBuffer(m->codec, out, false);
	return m->drained ? MW_MEDIA_END : MW_MEDIA_FRAME;
}

static int mw_media_frame(mw_media *m, uint32_t *dst, int stride, int rotate,
                          int32_t *w, int32_t *h) {
	return mw_image_to_argb(m->reader, dst, stride, rotate, w, h);
}

static void mw_media_close(mw_media *m) {
	if (m->codec) { AMediaCodec_stop(m->codec); AMediaCodec_delete(m->codec); }
	if (m->reader) mw_image_reader_delete(m->reader);
	if (m->window) ANativeWindow_release(m->window);
	if (m->ex) AMediaExtractor_delete(m->ex);
	memset(m, 0, sizeof(*m));
}
