//go:build android

// The image reader, and the conversion from what it produces into pixels.
//
// Both the camera and the video decoder end up here: each hands its frames
// to an AImageReader and each wants them as 0xAARRGGBB words. Keeping this
// in a file of its own rather than in one of their preambles is what lets
// them share it — a static function in a cgo preamble belongs to that one
// Go file and nothing else can see it.
#ifndef ANTUI_ANDROID_NDK_IMAGEREADER_H
#define ANTUI_ANDROID_NDK_IMAGEREADER_H

#include <media/NdkImageReader.h>
#include <android/native_window.h>
#include <stdint.h>

// mw_image_load opens libmediandk and finds the image reader. It returns 1
// when it is there, -1 when it is not — which is every device below API 24.
int mw_image_load(void);

// The image reader, wrapped so that Go never holds one of these types.
media_status_t mw_image_reader_new(int32_t width, int32_t height, int32_t format,
                                   int32_t maxImages, AImageReader **out);
void mw_image_reader_delete(AImageReader *r);
media_status_t mw_image_reader_window(AImageReader *r, ANativeWindow **out);

// mw_image_to_argb takes the newest image from the reader, converts it into
// 0xAARRGGBB words and turns it by rotate degrees clockwise on the way.
//
// Returns 1 for a frame, 0 for none waiting, -1 for a failure. The width and
// height written back are of the rotated image, so a quarter turn swaps them.
int mw_image_to_argb(AImageReader *r, uint32_t *dst, int dstStride, int rotate,
                     int32_t *outW, int32_t *outH);

#endif
