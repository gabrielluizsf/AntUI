//go:build android


#include <camera/NdkCameraManager.h>
#include <camera/NdkCameraDevice.h>
#include <camera/NdkCameraCaptureSession.h>
#include <camera/NdkCameraMetadataTags.h>
#include <android/native_window.h>
#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

#include "imagereader.h"

// libcamera2ndk.so is API 24 and is not in the stub libraries a build
// against API 21 links against, so it is opened and every call goes through
// a pointer — the same way AAudio and the sensor batching do. The image
// reader the frames arrive in is in imagereader.c, shared with the video
// decoder.
typedef struct {
	void *cam;

	ACameraManager *(*managerCreate)(void);
	void (*managerDelete)(ACameraManager *);
	camera_status_t (*idList)(ACameraManager *, ACameraIdList **);
	void (*idListDelete)(ACameraIdList *);
	camera_status_t (*characteristics)(ACameraManager *, const char *, ACameraMetadata **);
	void (*metadataFree)(ACameraMetadata *);
	camera_status_t (*metadataEntry)(const ACameraMetadata *, uint32_t, ACameraMetadata_const_entry *);
	camera_status_t (*openCamera)(ACameraManager *, const char *,
	                              ACameraDevice_StateCallbacks *, ACameraDevice **);
	camera_status_t (*deviceClose)(ACameraDevice *);
	camera_status_t (*createRequest)(const ACameraDevice *, ACameraDevice_request_template,
	                                 ACaptureRequest **);
	void (*requestFree)(ACaptureRequest *);
	camera_status_t (*targetCreate)(ANativeWindow *, ACameraOutputTarget **);
	void (*targetFree)(ACameraOutputTarget *);
	camera_status_t (*requestAddTarget)(ACaptureRequest *, const ACameraOutputTarget *);
	camera_status_t (*outputCreate)(ANativeWindow *, ACaptureSessionOutput **);
	void (*outputFree)(ACaptureSessionOutput *);
	camera_status_t (*containerCreate)(ACaptureSessionOutputContainer **);
	void (*containerFree)(ACaptureSessionOutputContainer *);
	camera_status_t (*containerAdd)(ACaptureSessionOutputContainer *,
	                                const ACaptureSessionOutput *);
	camera_status_t (*createSession)(ACameraDevice *, const ACaptureSessionOutputContainer *,
	                                 const ACameraCaptureSession_stateCallbacks *,
	                                 ACameraCaptureSession **);
	void (*sessionClose)(ACameraCaptureSession *);
	camera_status_t (*repeating)(ACameraCaptureSession *,
	                             ACameraCaptureSession_captureCallbacks *, int,
	                             ACaptureRequest **, int *);
	camera_status_t (*stopRepeating)(ACameraCaptureSession *);

} mw_camera_api;

static mw_camera_api K;
static int mw_camera_state;

#define MW_CAM(field, lib, name)                                    \
	K.field = dlsym(lib, name);                                     \
	if (K.field == NULL) { mw_camera_state = -1; return -1; }

int mw_camera_load(void) {
	if (mw_camera_state != 0) {
		return mw_camera_state;
	}
	K.cam = dlopen("libcamera2ndk.so", RTLD_NOW);
	if (K.cam == NULL || mw_image_load() != 1) {
		mw_camera_state = -1;
		return -1;
	}
	MW_CAM(managerCreate, K.cam, "ACameraManager_create")
	MW_CAM(managerDelete, K.cam, "ACameraManager_delete")
	MW_CAM(idList, K.cam, "ACameraManager_getCameraIdList")
	MW_CAM(idListDelete, K.cam, "ACameraManager_deleteCameraIdList")
	MW_CAM(characteristics, K.cam, "ACameraManager_getCameraCharacteristics")
	MW_CAM(metadataFree, K.cam, "ACameraMetadata_free")
	MW_CAM(metadataEntry, K.cam, "ACameraMetadata_getConstEntry")
	MW_CAM(openCamera, K.cam, "ACameraManager_openCamera")
	MW_CAM(deviceClose, K.cam, "ACameraDevice_close")
	MW_CAM(createRequest, K.cam, "ACameraDevice_createCaptureRequest")
	MW_CAM(requestFree, K.cam, "ACaptureRequest_free")
	MW_CAM(targetCreate, K.cam, "ACameraOutputTarget_create")
	MW_CAM(targetFree, K.cam, "ACameraOutputTarget_free")
	MW_CAM(requestAddTarget, K.cam, "ACaptureRequest_addTarget")
	MW_CAM(outputCreate, K.cam, "ACaptureSessionOutput_create")
	MW_CAM(outputFree, K.cam, "ACaptureSessionOutput_free")
	MW_CAM(containerCreate, K.cam, "ACaptureSessionOutputContainer_create")
	MW_CAM(containerFree, K.cam, "ACaptureSessionOutputContainer_free")
	MW_CAM(containerAdd, K.cam, "ACaptureSessionOutputContainer_add")
	MW_CAM(createSession, K.cam, "ACameraDevice_createCaptureSession")
	MW_CAM(sessionClose, K.cam, "ACameraCaptureSession_close")
	MW_CAM(repeating, K.cam, "ACameraCaptureSession_setRepeatingRequest")
	MW_CAM(stopRepeating, K.cam, "ACameraCaptureSession_stopRepeating")


	mw_camera_state = 1;
	return 1;
}

// The callbacks the platform insists on. It refuses to open a camera with a
// null table, so they exist; what they record is whether the camera went
// away, which Go asks about rather than being told.
static volatile int mw_cam_lost;
static volatile int mw_cam_error_code;

static void mw_on_disconnected(void *ctx, ACameraDevice *d) { mw_cam_lost = 1; }
static void mw_on_error(void *ctx, ACameraDevice *d, int err) {
	mw_cam_lost = 1;
	mw_cam_error_code = err;
}
static void mw_on_session(void *ctx, ACameraCaptureSession *s) {}

static ACameraDevice_StateCallbacks mw_device_callbacks = {
	NULL, mw_on_disconnected, mw_on_error,
};
static ACameraCaptureSession_stateCallbacks mw_session_callbacks = {
	NULL, mw_on_session, mw_on_session, mw_on_session,
};

int mw_camera_lost(void) { return mw_cam_lost; }
static void mw_camera_reset(void) { mw_cam_lost = 0; mw_cam_error_code = 0; }

// The wrappers Go calls. Each one is here rather than in Go because the
// pointer table cannot be reached from there.
ACameraManager *mw_manager_create(void) { return K.managerCreate(); }
void mw_manager_delete(ACameraManager *m) { K.managerDelete(m); }

camera_status_t mw_id_list(ACameraManager *m, ACameraIdList **l) {
	return K.idList(m, l);
}
void mw_id_list_delete(ACameraIdList *l) { K.idListDelete(l); }
int mw_id_count(ACameraIdList *l) { return l->numCameras; }
const char *mw_id_at(ACameraIdList *l, int i) { return l->cameraIds[i]; }

// One entry out of a camera's characteristics, as a single integer — which
// is what facing and orientation both are.
int mw_camera_int(ACameraManager *m, const char *id, uint32_t tag, int fallback) {
	ACameraMetadata *meta = NULL;
	if (K.characteristics(m, id, &meta) != ACAMERA_OK || meta == NULL) {
		return fallback;
	}
	ACameraMetadata_const_entry entry;
	memset(&entry, 0, sizeof(entry));
	int out = fallback;
	if (K.metadataEntry(meta, tag, &entry) == ACAMERA_OK && entry.count > 0) {
		if (entry.type == ACAMERA_TYPE_BYTE) {
			out = entry.data.u8[0];
		} else if (entry.type == ACAMERA_TYPE_INT32) {
			out = entry.data.i32[0];
		}
	}
	K.metadataFree(meta);
	return out;
}

camera_status_t mw_open(ACameraManager *m, const char *id, ACameraDevice **d) {
	mw_camera_reset();
	return K.openCamera(m, id, &mw_device_callbacks, d);
}
void mw_device_close(ACameraDevice *d) { K.deviceClose(d); }

// A whole preview in one call: a reader, a session and a repeating request.
// Doing it piecemeal from Go would mean nine more wrappers and nine more
// places to leak one of the handles.
typedef struct {
	AImageReader *reader;
	ANativeWindow *window;
	ACaptureSessionOutput *output;
	ACaptureSessionOutputContainer *container;
	ACameraOutputTarget *target;
	ACaptureRequest *request;
	ACameraCaptureSession *session;
} mw_preview;

camera_status_t mw_preview_start(ACameraDevice *dev, int32_t width, int32_t height,
                                        int32_t format, mw_preview *p) {
	memset(p, 0, sizeof(*p));
	// Two images in flight: one being read and one being filled. More only
	// adds delay between what the lens sees and what is drawn.
	if (mw_image_reader_new(width, height, format, 2, &p->reader) != AMEDIA_OK) {
		return ACAMERA_ERROR_UNKNOWN;
	}
	if (mw_image_reader_window(p->reader, &p->window) != AMEDIA_OK) {
		return ACAMERA_ERROR_UNKNOWN;
	}
	ANativeWindow_acquire(p->window);

	camera_status_t r = K.outputCreate(p->window, &p->output);
	if (r != ACAMERA_OK) return r;
	r = K.containerCreate(&p->container);
	if (r != ACAMERA_OK) return r;
	r = K.containerAdd(p->container, p->output);
	if (r != ACAMERA_OK) return r;
	r = K.createRequest(dev, TEMPLATE_PREVIEW, &p->request);
	if (r != ACAMERA_OK) return r;
	r = K.targetCreate(p->window, &p->target);
	if (r != ACAMERA_OK) return r;
	r = K.requestAddTarget(p->request, p->target);
	if (r != ACAMERA_OK) return r;
	r = K.createSession(dev, p->container, &mw_session_callbacks, &p->session);
	if (r != ACAMERA_OK) return r;
	return K.repeating(p->session, NULL, 1, &p->request, NULL);
}

void mw_preview_stop(mw_preview *p) {
	if (p->session) { K.stopRepeating(p->session); K.sessionClose(p->session); }
	if (p->request) K.requestFree(p->request);
	if (p->target) K.targetFree(p->target);
	if (p->container) K.containerFree(p->container);
	if (p->output) K.outputFree(p->output);
	if (p->window) ANativeWindow_release(p->window);
	if (p->reader) mw_image_reader_delete(p->reader);
	memset(p, 0, sizeof(*p));
}

// mw_preview_frame is mw_image_to_argb, named for what the camera calls it.
int mw_preview_frame(mw_preview *p, uint32_t *dst, int dstStride, int rotate,
                            int32_t *outW, int32_t *outH) {
	return mw_image_to_argb(p->reader, dst, dstStride, rotate, outW, outH);
}
