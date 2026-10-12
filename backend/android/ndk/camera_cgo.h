//go:build android

// Shared by every file in this package that drives a camera. The pointer
// table is loaded once and the state the callbacks record is read from the
// Go side, so the C has to live in one translation unit — see camera_cgo.c.
// A static helper in a cgo preamble would belong to one Go file only.

#include <camera/NdkCameraManager.h>
#include <camera/NdkCameraDevice.h>
#include <camera/NdkCameraCaptureSession.h>
#include <camera/NdkCameraMetadataTags.h>
#include <android/native_window.h>
#include <stdint.h>

#include "imagereader.h"

// A whole preview: a reader, a session and a repeating request.
typedef struct {
	AImageReader *reader;
	ANativeWindow *window;
	ACaptureSessionOutput *output;
	ACaptureSessionOutputContainer *container;
	ACameraOutputTarget *target;
	ACaptureRequest *request;
	ACameraCaptureSession *session;
} mw_preview;

int mw_camera_load(void);
ACameraManager *mw_manager_create(void);
void mw_manager_delete(ACameraManager *m);
camera_status_t mw_id_list(ACameraManager *m, ACameraIdList **l);
void mw_id_list_delete(ACameraIdList *l);
int mw_id_count(ACameraIdList *l);
const char *mw_id_at(ACameraIdList *l, int i);
int mw_camera_int(ACameraManager *m, const char *id, uint32_t tag, int fallback);
camera_status_t mw_open(ACameraManager *m, const char *id, ACameraDevice **d);
void mw_device_close(ACameraDevice *d);
camera_status_t mw_preview_start(ACameraDevice *dev, int32_t width, int32_t height,
                                 int32_t format, mw_preview *p);
void mw_preview_stop(mw_preview *p);
int mw_preview_frame(mw_preview *p, uint32_t *dst, int dstStride, int rotate,
                     int32_t *outW, int32_t *outH);
int mw_camera_lost(void);
