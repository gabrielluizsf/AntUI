//go:build android

// Shared for every file in this package that reads a sensor event: the
// platform's event payload is anonymous unions that cgo cannot name, so the
// reading and the API-26 calls are done here in C.

#include <android/sensor.h>
#include <dlfcn.h>
#include <stdlib.h>

// ASensorEvent's payload is a nest of anonymous unions, which cgo cannot
// name — so reading one has to happen in C. All of them are floats except
// the step counter, which is a 64-bit integer sharing the same bytes.
static float mw_sensor_data(const ASensorEvent *e, int i) { return e->data[i]; }
static int64_t mw_sensor_time(const ASensorEvent *e)      { return e->timestamp; }
static int32_t mw_sensor_type(const ASensorEvent *e)      { return e->type; }
static int32_t mw_sensor_id(const ASensorEvent *e)        { return e->sensor; }

// ASensorEventQueue_registerSensor is API 26 and is the only way to ask for
// batching — the platform collecting readings and delivering them in bursts,
// which is what lets the processor sleep between them. Below 26 there is
// enableSensor and a rate, and no batching.
//
// It is loaded by hand rather than called, because a build against API 21
// does not declare it at all. RTLD_DEFAULT finds it in libandroid, which is
// already open.
typedef int (*mw_register_fn)(ASensorEventQueue *, const ASensor *, int32_t, int64_t);
typedef ASensorManager *(*mw_for_package_fn)(const char *);

// getInstanceForPackage is API 26 as well, and is the call the platform
// wants an app to use — the one without a name is deprecated and, on some
// versions, sees fewer sensors.
static ASensorManager *mw_manager_for(const char *name) {
	static mw_for_package_fn fn;
	static int looked;
	if (!looked) {
		looked = 1;
		fn = (mw_for_package_fn)dlsym(RTLD_DEFAULT, "ASensorManager_getInstanceForPackage");
	}
	if (fn == NULL) {
		return ASensorManager_getInstance();
	}
	return fn(name);
}

static mw_register_fn mw_register_sensor(void) {
	static mw_register_fn fn;
	static int looked;
	if (!looked) {
		looked = 1;
		fn = (mw_register_fn)dlsym(RTLD_DEFAULT, "ASensorEventQueue_registerSensor");
	}
	return fn;
}

static int mw_register(ASensorEventQueue *q, const ASensor *s, int32_t rate, int64_t latency) {
	mw_register_fn fn = mw_register_sensor();
	if (fn == NULL) {
		return -1;
	}
	return fn(q, s, rate, latency);
}