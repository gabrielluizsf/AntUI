//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

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

static int mw_camera_load(void) {
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

static int mw_camera_lost(void) { return mw_cam_lost; }
static void mw_camera_reset(void) { mw_cam_lost = 0; mw_cam_error_code = 0; }

// The wrappers Go calls. Each one is here rather than in Go because the
// pointer table cannot be reached from there.
static ACameraManager *mw_manager_create(void) { return K.managerCreate(); }
static void mw_manager_delete(ACameraManager *m) { K.managerDelete(m); }

static camera_status_t mw_id_list(ACameraManager *m, ACameraIdList **l) {
	return K.idList(m, l);
}
static void mw_id_list_delete(ACameraIdList *l) { K.idListDelete(l); }
static int mw_id_count(ACameraIdList *l) { return l->numCameras; }
static const char *mw_id_at(ACameraIdList *l, int i) { return l->cameraIds[i]; }

// One entry out of a camera's characteristics, as a single integer — which
// is what facing and orientation both are.
static int mw_camera_int(ACameraManager *m, const char *id, uint32_t tag, int fallback) {
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

static camera_status_t mw_open(ACameraManager *m, const char *id, ACameraDevice **d) {
	mw_camera_reset();
	return K.openCamera(m, id, &mw_device_callbacks, d);
}
static void mw_device_close(ACameraDevice *d) { K.deviceClose(d); }

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

static camera_status_t mw_preview_start(ACameraDevice *dev, int32_t width, int32_t height,
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

static void mw_preview_stop(mw_preview *p) {
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
static int mw_preview_frame(mw_preview *p, uint32_t *dst, int dstStride, int rotate,
                            int32_t *outW, int32_t *outH) {
	return mw_image_to_argb(p->reader, dst, dstStride, rotate, outW, outH);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// ErrNoCamera2 is a device too old for the NDK's camera API, which arrived
// in API 24.
var ErrNoCamera2 = errors.New("ndk: this device has no NDK camera; it arrived in API 24")

// CameraAvailable reports whether the camera API can be used at all.
func CameraAvailable() bool { return C.mw_camera_load() == 1 }

// Which way a camera points.
const (
	CameraBack     = 1
	CameraFront    = 0
	CameraExternal = 2
)

// The image format a preview asks for. YUV420 is the one every camera
// supports and the only one worth asking for.
const FormatYUV420 = 0x23

// CameraManager is the list of cameras.
type CameraManager struct{ ptr *C.ACameraManager }

// Cameras opens the camera manager.
func Cameras() (*CameraManager, error) {
	if C.mw_camera_load() != 1 {
		return nil, ErrNoCamera2
	}
	m := C.mw_manager_create()
	if m == nil {
		return nil, errors.New("ndk: cannot open the camera manager")
	}
	return &CameraManager{ptr: m}, nil
}

// Close gives the manager back.
func (m *CameraManager) Close() {
	if m.ptr != nil {
		C.mw_manager_delete(m.ptr)
		m.ptr = nil
	}
}

// IDs is every camera the app may open. They are opaque strings, not
// indices: "0" and "1" on most phones, and anything at all on a device with
// an external one.
func (m *CameraManager) IDs() ([]string, error) {
	var list *C.ACameraIdList
	if r := C.mw_id_list(m.ptr, &list); r != C.ACAMERA_OK {
		return nil, cameraError("listing the cameras", r)
	}
	defer C.mw_id_list_delete(list)
	n := int(C.mw_id_count(list))
	out := make([]string, n)
	for i := range n {
		out[i] = C.GoString(C.mw_id_at(list, C.int(i)))
	}
	return out, nil
}

// Facing is which way a camera points: [CameraBack], [CameraFront] or
// [CameraExternal].
func (m *CameraManager) Facing(id string) int {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	return int(C.mw_camera_int(m.ptr, cid, C.ACAMERA_LENS_FACING, C.int(CameraBack)))
}

// Orientation is how far the camera's own image is rotated from the device's
// natural orientation, in degrees clockwise: 90 on almost every phone.
//
// It is why photographs come out sideways. The sensor is mounted turned, the
// image arrives turned, and nothing rotates it — the app has to.
func (m *CameraManager) Orientation(id string) int {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	return int(C.mw_camera_int(m.ptr, cid, C.ACAMERA_SENSOR_ORIENTATION, 90))
}

// CameraDevice is an open camera.
type CameraDevice struct {
	ptr     *C.ACameraDevice
	preview C.mw_preview
	running bool
}

// Open opens a camera. It needs the CAMERA permission, and fails rather than
// returning a black picture without it.
func (m *CameraManager) Open(id string) (*CameraDevice, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	var dev *C.ACameraDevice
	if r := C.mw_open(m.ptr, cid, &dev); r != C.ACAMERA_OK {
		return nil, cameraError("opening camera "+id, r)
	}
	return &CameraDevice{ptr: dev}, nil
}

// StartPreview begins a stream of frames at the given size.
func (d *CameraDevice) StartPreview(width, height int) error {
	if d.running {
		return errors.New("ndk: the preview is already running")
	}
	if r := C.mw_preview_start(d.ptr, C.int32_t(width), C.int32_t(height),
		C.int32_t(FormatYUV420), &d.preview); r != C.ACAMERA_OK {
		C.mw_preview_stop(&d.preview)
		return cameraError("starting the preview", r)
	}
	d.running = true
	return nil
}

// Frame converts the newest frame into dst, which holds 0xAARRGGBB words
// with stride words between rows, turning it by rotate degrees clockwise on
// the way — 0, 90, 180 or 270. The width and height that come back are of
// the *rotated* image, so they are swapped for a quarter turn.
//
// It reports whether there was a frame: a preview produces them on its own
// schedule and asking more often than it produces them is normal. Frames
// that arrived and were not read are dropped, because what is wanted from a
// camera is always the newest.
func (d *CameraDevice) Frame(dst []uint32, stride, rotate int) (w, h int, ok bool, err error) {
	if !d.running || len(dst) == 0 {
		return 0, 0, false, nil
	}
	var cw, ch C.int32_t
	r := C.mw_preview_frame(&d.preview, (*C.uint32_t)(unsafe.Pointer(&dst[0])),
		C.int(stride), C.int(rotate), &cw, &ch)
	switch r {
	case 0:
		return 0, 0, false, nil
	case 1:
		return int(cw), int(ch), true, nil
	}
	return 0, 0, false, errors.New("ndk: cannot read the camera's image")
}

// Lost reports whether the camera was taken away — unplugged, or claimed by
// another app, which the platform allows and does not ask about.
func (d *CameraDevice) Lost() bool { return C.mw_camera_lost() != 0 }

// Close stops the preview and gives the camera back.
func (d *CameraDevice) Close() {
	if d.running {
		C.mw_preview_stop(&d.preview)
		d.running = false
	}
	if d.ptr != nil {
		C.mw_device_close(d.ptr)
		d.ptr = nil
	}
}

func cameraError(what string, r C.camera_status_t) error {
	return fmt.Errorf("ndk: camera: %s failed (%d)", what, int(r))
}
