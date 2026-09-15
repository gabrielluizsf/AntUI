//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

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
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"time"
	"unsafe"
)

// newExtractor opens a file and hands the extractor its descriptor.
//
// Not setDataSource, which takes a *URI* rather than a path — a bare path is
// refused with AMEDIA_ERROR_UNSUPPORTED, which says nothing about why. A
// descriptor also means the same code will one day serve a content:// address
// from the picker, which arrives as a descriptor and never as a path.
//
// The file has to stay open for as long as the extractor reads from it, so
// it is handed back to be closed with it.
func newExtractor(path string) (*C.AMediaExtractor, *os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	ex := C.AMediaExtractor_new()
	if ex == nil {
		f.Close()
		return nil, nil, errors.New("ndk: cannot make a media extractor")
	}
	if r := C.AMediaExtractor_setDataSourceFd(ex, C.int(f.Fd()), 0,
		C.off64_t(st.Size())); r != C.AMEDIA_OK {
		C.AMediaExtractor_delete(ex)
		f.Close()
		return nil, nil, fmt.Errorf("ndk: nothing on this device can read %s (%d)",
			path, int(r))
	}
	return ex, f, nil
}

// Track is one stream inside a media file: a picture, a sound, a subtitle.
type Track struct {
	Index int
	// Mime is what it holds — "video/avc", "audio/mp4a-latm".
	Mime string
	// Duration is how long the whole file is, which the container reports
	// per track and which is usually the same for all of them.
	Duration time.Duration

	// Width and Height are set on a video track, and Rotation is how far it
	// has to be turned to look right — a video shot on a phone held sideways
	// carries 90 here and is stored as if it were not.
	Width, Height int
	Rotation      int
	FrameRate     int

	// SampleRate and Channels are set on an audio track.
	SampleRate int
	Channels   int
}

// Video reports whether this is a picture track.
func (t Track) Video() bool { return len(t.Mime) > 6 && t.Mime[:6] == "video/" }

// Audio reports whether this is a sound track.
func (t Track) Audio() bool { return len(t.Mime) > 6 && t.Mime[:6] == "audio/" }

// Media is an open media file, taken apart into its tracks.
type Media struct {
	path   string
	tracks []Track
}

// OpenMedia reads a file's table of contents. It decodes nothing: that is
// [Media.DecodeVideo] and [Media.DecodeAudio], each of which opens the file
// again for a track of its own.
func OpenMedia(path string) (*Media, error) {
	ex, f, err := newExtractor(path)
	if err != nil {
		return nil, err
	}
	defer C.AMediaExtractor_delete(ex)
	defer f.Close()

	n := int(C.AMediaExtractor_getTrackCount(ex))
	m := &Media{path: path, tracks: make([]Track, 0, n)}
	for i := range n {
		format := C.AMediaExtractor_getTrackFormat(ex, C.size_t(i))
		if format == nil {
			continue
		}
		t := Track{Index: i}
		t.Mime = formatString(format, "mime")
		if us, ok := formatInt64(format, "durationUs"); ok {
			t.Duration = time.Duration(us) * time.Microsecond
		}
		t.Width, _ = formatInt(format, "width")
		t.Height, _ = formatInt(format, "height")
		t.Rotation, _ = formatInt(format, "rotation-degrees")
		t.FrameRate, _ = formatInt(format, "frame-rate")
		t.SampleRate, _ = formatInt(format, "sample-rate")
		t.Channels, _ = formatInt(format, "channel-count")
		C.AMediaFormat_delete(format)
		m.tracks = append(m.tracks, t)
	}
	if len(m.tracks) == 0 {
		return nil, fmt.Errorf("ndk: %s has no tracks this device can read", path)
	}
	return m, nil
}

// Tracks is what the file holds.
func (m *Media) Tracks() []Track { return m.tracks }

// Track finds the first track of a kind, and reports whether there was one.
func (m *Media) Track(video bool) (Track, bool) {
	for _, t := range m.tracks {
		if t.Video() == video && t.Audio() != video {
			return t, true
		}
	}
	return Track{}, false
}

func formatString(f *C.AMediaFormat, key string) string {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out *C.char
	if !C.AMediaFormat_getString(f, k, &out) {
		return ""
	}
	return C.GoString(out)
}

func formatInt(f *C.AMediaFormat, key string) (int, bool) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out C.int32_t
	if !C.AMediaFormat_getInt32(f, k, &out) {
		return 0, false
	}
	return int(out), true
}

func formatInt64(f *C.AMediaFormat, key string) (int64, bool) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var out C.int64_t
	if !C.AMediaFormat_getInt64(f, k, &out) {
		return 0, false
	}
	return int64(out), true
}

// Decoder is one track being decoded.
type Decoder struct {
	media C.mw_media
	// file is the one the extractor reads from, kept open for its life.
	file   *os.File
	track  Track
	rotate int
	closed bool
	// ended is set once the codec has produced its last output.
	ended bool
}

// The states a pump can be in.
const (
	pumpAgain = 0
	pumpFrame = 1
	pumpEnd   = 2
	pumpError = -1
)

// DecodeVideo starts decoding a picture track into frames.
//
// rotate is how far to turn each frame, in degrees clockwise. Passing the
// track's own [Track.Rotation] is what makes a video shot on a phone held
// sideways come out the right way up; passing 0 gives it as stored.
//
// It needs API 24, because the frames arrive through an image reader.
func (m *Media) DecodeVideo(track Track, rotate int) (*Decoder, error) {
	if !track.Video() {
		return nil, fmt.Errorf("ndk: track %d is %s, not a picture", track.Index, track.Mime)
	}
	if track.Width <= 0 || track.Height <= 0 {
		return nil, fmt.Errorf("ndk: track %d does not say how big it is", track.Index)
	}
	d := &Decoder{track: track, rotate: rotate}
	if err := d.open(m.path, track); err != nil {
		return nil, err
	}
	// The reader is the surface the codec draws into, at the video's own
	// size — scaling here would be scaling before anything has been seen.
	if r := C.mw_image_reader_new(C.int32_t(track.Width), C.int32_t(track.Height),
		C.int32_t(FormatYUV420), 3, &d.media.reader); r != C.AMEDIA_OK {
		d.Close()
		return nil, ErrNoImageReader
	}
	if r := C.mw_image_reader_window(d.media.reader, &d.media.window); r != C.AMEDIA_OK {
		d.Close()
		return nil, errors.New("ndk: the image reader has no surface")
	}
	C.ANativeWindow_acquire(d.media.window)

	if err := d.configure(d.media.window); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// ErrNoImageReader is a device below API 24, where video frames cannot be
// got at from native code at all.
var ErrNoImageReader = errors.New("ndk: this device cannot hand decoded video " +
	"to native code; the image reader arrived in API 24")

// DecodeAudio starts decoding a sound track into PCM.
func (m *Media) DecodeAudio(track Track) (*Decoder, error) {
	if !track.Audio() {
		return nil, fmt.Errorf("ndk: track %d is %s, not a sound", track.Index, track.Mime)
	}
	d := &Decoder{track: track}
	if err := d.open(m.path, track); err != nil {
		return nil, err
	}
	if err := d.configure(nil); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// open makes an extractor of this decoder's own and selects its track. Each
// decoder gets one: an extractor has a single position, and two tracks read
// at different rates cannot share it.
func (d *Decoder) open(path string, track Track) error {
	ex, f, err := newExtractor(path)
	if err != nil {
		return err
	}
	d.media.ex = ex
	d.file = f
	if r := C.AMediaExtractor_selectTrack(d.media.ex, C.size_t(track.Index)); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot select track %d (%d)", track.Index, int(r))
	}
	return nil
}

func (d *Decoder) configure(window *C.ANativeWindow) error {
	mime := C.CString(d.track.Mime)
	defer C.free(unsafe.Pointer(mime))
	d.media.codec = C.AMediaCodec_createDecoderByType(mime)
	if d.media.codec == nil {
		return fmt.Errorf("ndk: this device has no decoder for %s", d.track.Mime)
	}
	format := C.AMediaExtractor_getTrackFormat(d.media.ex, C.size_t(d.track.Index))
	if format == nil {
		return errors.New("ndk: the track has no format")
	}
	defer C.AMediaFormat_delete(format)

	if r := C.AMediaCodec_configure(d.media.codec, format, window, nil, 0); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot set up the %s decoder (%d)", d.track.Mime, int(r))
	}
	if r := C.AMediaCodec_start(d.media.codec); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot start the %s decoder (%d)", d.track.Mime, int(r))
	}
	return nil
}

// Frame decodes until a picture is ready, and converts it into dst.
//
// It reports whether there is one: a decoder needs several compressed
// samples before it produces its first frame, so a call that gives nothing
// is normal and not the end. The end is [Decoder.Ended].
func (d *Decoder) Frame(dst []uint32, stride int) (w, h int, pts time.Duration, ok bool, err error) {
	if d.closed || len(dst) == 0 {
		return 0, 0, 0, false, nil
	}
	switch C.mw_video_pump(&d.media) {
	case pumpError:
		return 0, 0, 0, false, errors.New("ndk: the video decoder failed")
	case pumpEnd:
		d.ended = true
	}
	var cw, ch C.int32_t
	got := C.mw_media_frame(&d.media, (*C.uint32_t)(unsafe.Pointer(&dst[0])),
		C.int(stride), C.int(d.rotate), &cw, &ch)
	if got != 1 {
		return 0, 0, 0, false, nil
	}
	return int(cw), int(ch), time.Duration(d.media.pts) * time.Microsecond, true, nil
}

// Samples decodes sound into dst, which is filled interleaved, and reports
// how many floats were written.
func (d *Decoder) Samples(dst []float32) (n int, pts time.Duration, err error) {
	if d.closed || len(dst) == 0 {
		return 0, 0, nil
	}
	var wrote C.int
	switch C.mw_audio_pump(&d.media, (*C.float)(unsafe.Pointer(&dst[0])),
		C.int(len(dst)), &wrote) {
	case pumpError:
		return 0, 0, errors.New("ndk: the audio decoder failed")
	case pumpEnd:
		d.ended = true
	}
	return int(wrote), time.Duration(d.media.pts) * time.Microsecond, nil
}

// Ended reports whether the codec has produced everything it is going to.
func (d *Decoder) Ended() bool { return d.ended }

// Track is what is being decoded.
func (d *Decoder) Track() Track { return d.track }

// Seek moves to a position. The decoder jumps to the nearest key frame at or
// before it, which for video may be a second earlier — there is no other
// kind of seek in a compressed stream.
func (d *Decoder) Seek(to time.Duration) error {
	if d.closed {
		return errors.New("ndk: the decoder is closed")
	}
	if r := C.AMediaExtractor_seekTo(d.media.ex, C.int64_t(to.Microseconds()),
		C.AMEDIAEXTRACTOR_SEEK_PREVIOUS_SYNC); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot seek to %v (%d)", to, int(r))
	}
	C.AMediaCodec_flush(d.media.codec)
	d.media.eof = 0
	d.media.drained = 0
	d.ended = false
	return nil
}

// Close stops the decoder and gives everything back.
func (d *Decoder) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	C.mw_media_close(&d.media)
	if d.file != nil {
		d.file.Close()
		d.file = nil
	}
	return nil
}
