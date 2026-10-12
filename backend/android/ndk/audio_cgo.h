//go:build android

// Shared by every file in this package that plays or records sound: the
// library is opened once and the pointer table lives in one place, so the C
// stays in one translation unit — see audio_cgo.c. A static helper in a cgo
// preamble would belong to one Go file and its state to that file alone.

#include <aaudio/AAudio.h>
#include <stdint.h>

int mw_aaudio_load(void);
aaudio_result_t mw_audio_open(int32_t rate, int32_t channels, int32_t capacity,
                              int32_t input, AAudioStream **out,
                              int32_t *gotRate, int32_t *gotChannels,
                              int32_t *burst);
aaudio_result_t mw_audio_start(AAudioStream *s);
aaudio_result_t mw_audio_stop(AAudioStream *s);
aaudio_result_t mw_audio_close(AAudioStream *s);
int32_t mw_audio_buffer(AAudioStream *s);
aaudio_result_t mw_audio_set_buffer(AAudioStream *s, int32_t n);
aaudio_result_t mw_audio_write(AAudioStream *s, const float *data,
                               int32_t frames, int64_t timeout);
aaudio_result_t mw_audio_read(AAudioStream *s, float *data,
                              int32_t frames, int64_t timeout);
const char *mw_audio_text(aaudio_result_t r);
