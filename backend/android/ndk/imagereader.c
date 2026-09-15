//go:build android

#include "imagereader.h"

#include <dlfcn.h>
#include <string.h>

// AImageReader is API 24, and that half of libmediandk.so is not in the stub
// libraries a build against API 21 links against. The rest of libmediandk —
// the extractor and the codec — is API 21 and is linked normally; only this
// is looked up.
static struct {
	void *lib;
	media_status_t (*readerNew)(int32_t, int32_t, int32_t, int32_t, AImageReader **);
	void (*readerDelete)(AImageReader *);
	media_status_t (*readerWindow)(AImageReader *, ANativeWindow **);
	media_status_t (*readerLatest)(AImageReader *, AImage **);
	void (*imageDelete)(AImage *);
	media_status_t (*imageWidth)(const AImage *, int32_t *);
	media_status_t (*imageHeight)(const AImage *, int32_t *);
	media_status_t (*planeData)(const AImage *, int, uint8_t **, int *);
	media_status_t (*rowStride)(const AImage *, int, int32_t *);
	media_status_t (*pixelStride)(const AImage *, int, int32_t *);
} R;

static int mw_image_state;

#define MW_IMG(field, name)                                       \
	R.field = dlsym(R.lib, name);                                 \
	if (R.field == NULL) { mw_image_state = -1; return -1; }

int mw_image_load(void) {
	if (mw_image_state != 0) {
		return mw_image_state;
	}
	R.lib = dlopen("libmediandk.so", RTLD_NOW);
	if (R.lib == NULL) {
		mw_image_state = -1;
		return -1;
	}
	MW_IMG(readerNew, "AImageReader_new")
	MW_IMG(readerDelete, "AImageReader_delete")
	MW_IMG(readerWindow, "AImageReader_getWindow")
	MW_IMG(readerLatest, "AImageReader_acquireLatestImage")
	MW_IMG(imageDelete, "AImage_delete")
	MW_IMG(imageWidth, "AImage_getWidth")
	MW_IMG(imageHeight, "AImage_getHeight")
	MW_IMG(planeData, "AImage_getPlaneData")
	MW_IMG(rowStride, "AImage_getPlaneRowStride")
	MW_IMG(pixelStride, "AImage_getPlanePixelStride")
	mw_image_state = 1;
	return 1;
}

media_status_t mw_image_reader_new(int32_t width, int32_t height, int32_t format,
                                   int32_t maxImages, AImageReader **out) {
	if (mw_image_load() != 1) {
		return AMEDIA_ERROR_UNSUPPORTED;
	}
	return R.readerNew(width, height, format, maxImages, out);
}

void mw_image_reader_delete(AImageReader *r) {
	if (r != NULL) {
		R.readerDelete(r);
	}
}

media_status_t mw_image_reader_window(AImageReader *r, ANativeWindow **out) {
	return R.readerWindow(r, out);
}

// The conversion visits every pixel of every frame — two million of them at
// 1080p, thirty times a second — so it is here rather than in Go, and the
// arithmetic is integer BT.601, which is what YUV_420_888 is.
//
// The rotation costs nothing: it is a different index for the same write.
// A camera sensor is mounted turned on almost every phone and a video may
// carry a rotation of its own, so a frame that was not turned here would
// have to be turned in a second pass over all two million.
int mw_image_to_argb(AImageReader *r, uint32_t *dst, int dstStride, int rotate,
                     int32_t *outW, int32_t *outH) {
	AImage *img = NULL;
	if (R.readerLatest(r, &img) != AMEDIA_OK || img == NULL) {
		return 0;
	}
	int32_t w = 0, h = 0;
	R.imageWidth(img, &w);
	R.imageHeight(img, &h);
	if (rotate == 90 || rotate == 270) {
		*outW = h;
		*outH = w;
	} else {
		*outW = w;
		*outH = h;
	}

	uint8_t *yData = NULL, *uData = NULL, *vData = NULL;
	int yLen = 0, uLen = 0, vLen = 0;
	int32_t yRow = 0, uRow = 0, vRow = 0, uPix = 1, vPix = 1;
	if (R.planeData(img, 0, &yData, &yLen) != AMEDIA_OK ||
	    R.planeData(img, 1, &uData, &uLen) != AMEDIA_OK ||
	    R.planeData(img, 2, &vData, &vLen) != AMEDIA_OK) {
		R.imageDelete(img);
		return -1;
	}
	R.rowStride(img, 0, &yRow);
	R.rowStride(img, 1, &uRow);
	R.rowStride(img, 2, &vRow);
	R.pixelStride(img, 1, &uPix);
	R.pixelStride(img, 2, &vPix);

	for (int32_t y = 0; y < h; y++) {
		const uint8_t *yr = yData + (int64_t)y * yRow;
		const uint8_t *ur = uData + (int64_t)(y >> 1) * uRow;
		const uint8_t *vr = vData + (int64_t)(y >> 1) * vRow;
		for (int32_t x = 0; x < w; x++) {
			int Y = yr[x];
			int U = ur[(x >> 1) * uPix] - 128;
			int V = vr[(x >> 1) * vPix] - 128;
			int R_ = (65536 * Y + 91881 * V) >> 16;
			int G_ = (65536 * Y - 22554 * U - 46802 * V) >> 16;
			int B_ = (65536 * Y + 116130 * U) >> 16;
			if (R_ < 0) R_ = 0; else if (R_ > 255) R_ = 255;
			if (G_ < 0) G_ = 0; else if (G_ > 255) G_ = 255;
			if (B_ < 0) B_ = 0; else if (B_ > 255) B_ = 255;
			uint32_t argb = 0xFF000000u | ((uint32_t)R_ << 16) |
			                ((uint32_t)G_ << 8) | (uint32_t)B_;
			int64_t at;
			switch (rotate) {
			case 90:
				at = (int64_t)x * dstStride + (h - 1 - y);
				break;
			case 180:
				at = (int64_t)(h - 1 - y) * dstStride + (w - 1 - x);
				break;
			case 270:
				at = (int64_t)(w - 1 - x) * dstStride + y;
				break;
			default:
				at = (int64_t)y * dstStride + x;
			}
			dst[at] = argb;
		}
	}
	R.imageDelete(img);
	return 1;
}
