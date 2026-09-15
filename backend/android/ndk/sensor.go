//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

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
*/
import "C"

import (
	"errors"
	"unsafe"
)

// SensorManager is the platform's list of sensors and the maker of queues.
type SensorManager struct {
	ptr *C.ASensorManager
}

// Sensors opens the sensor manager for a package.
//
// The package name matters: since API 26 the platform uses it to decide
// which sensors an app may see, and the older call that took no name is
// deprecated. Passing an empty name uses the old one deliberately; passing a
// name uses the new one where it exists and the old one where it does not.
func Sensors(packageName string) (*SensorManager, error) {
	var m *C.ASensorManager
	if packageName == "" {
		m = C.ASensorManager_getInstance()
	} else {
		name := C.CString(packageName)
		defer C.free(unsafe.Pointer(name))
		// Falls back to the nameless call on a platform too old for the
		// other, which is what the C helper does.
		m = C.mw_manager_for(name)
	}
	if m == nil {
		return nil, errors.New("ndk: this device has no sensor manager")
	}
	return &SensorManager{ptr: m}, nil
}

// Sensor is one physical or computed sensor.
type Sensor struct {
	ptr *C.ASensor
}

// IsNil reports whether the device has no such sensor.
func (s Sensor) IsNil() bool { return s.ptr == nil }

// Default is the sensor of a kind that the platform considers the best one,
// or a nil Sensor when the device has none.
func (m *SensorManager) Default(kind int) Sensor {
	return Sensor{ptr: C.ASensorManager_getDefaultSensor(m.ptr, C.int(kind))}
}

// DefaultWake is Default, asking specifically for a sensor that can wake the
// device or specifically for one that cannot. A wake-up sensor keeps the
// processor alive to deliver its readings; a non-wake-up one is silent while
// the device sleeps and reports what it collected when it wakes.
func (m *SensorManager) DefaultWake(kind int, wake bool) Sensor {
	var w C.bool
	if wake {
		w = true
	}
	return Sensor{ptr: C.ASensorManager_getDefaultSensorEx(m.ptr, C.int(kind), w)}
}

// List is every sensor the device has.
func (m *SensorManager) List() []Sensor {
	var list C.ASensorList
	n := int(C.ASensorManager_getSensorList(m.ptr, &list))
	if n <= 0 || list == nil {
		return nil
	}
	raw := unsafe.Slice((**C.ASensor)(unsafe.Pointer(list)), n)
	out := make([]Sensor, n)
	for i := range raw {
		out[i] = Sensor{ptr: raw[i]}
	}
	return out
}

// Name, Vendor, Type, Resolution and MinDelay describe one sensor.
func (s Sensor) Name() string { return C.GoString(C.ASensor_getName(s.ptr)) }

// Vendor is who made the sensor, which is the only way to tell two of
// the same kind apart.
func (s Sensor) Vendor() string { return C.GoString(C.ASensor_getVendor(s.ptr)) }

// Type is what the sensor measures, as one of the ASENSOR_TYPE values.
func (s Sensor) Type() int { return int(C.ASensor_getType(s.ptr)) }

// Resolution is the smallest change the sensor can report, in its own units.
func (s Sensor) Resolution() float32 { return float32(C.ASensor_getResolution(s.ptr)) }

// MinDelay is the shortest interval between readings, in microseconds. It is
// 0 for a sensor that only reports when something changes — a step detector,
// a proximity sensor — and asking such a sensor for a rate does nothing.
func (s Sensor) MinDelay() int { return int(C.ASensor_getMinDelay(s.ptr)) }

// SensorQueue is a stream of readings, delivered to a looper.
type SensorQueue struct {
	ptr     *C.ASensorEventQueue
	manager *SensorManager
}

// SensorID is the identifier the queue is attached to the looper with. It is
// distinct from [InputID] and from nothing else.
const SensorID = 2

// Queue makes a stream on the given looper. Like the input queue, it is
// thread-local: readings arrive only on the thread the looper belongs to.
func (m *SensorManager) Queue(l *Looper) (*SensorQueue, error) {
	q := C.ASensorManager_createEventQueue(m.ptr, l.ptr, C.int(SensorID), nil, nil)
	if q == nil {
		return nil, errors.New("ndk: cannot make a sensor queue")
	}
	return &SensorQueue{ptr: q, manager: m}, nil
}

// Close destroys the queue. Every sensor on it stops.
func (q *SensorQueue) Close() {
	if q.ptr != nil {
		C.ASensorManager_destroyEventQueue(q.manager.ptr, q.ptr)
		q.ptr = nil
	}
}

// Enable starts a sensor at its default rate.
func (q *SensorQueue) Enable(s Sensor) error {
	if C.ASensorEventQueue_enableSensor(q.ptr, s.ptr) < 0 {
		return errors.New("ndk: cannot enable " + s.Name())
	}
	return nil
}

// Disable stops one.
func (q *SensorQueue) Disable(s Sensor) error {
	if C.ASensorEventQueue_disableSensor(q.ptr, s.ptr) < 0 {
		return errors.New("ndk: cannot disable " + s.Name())
	}
	return nil
}

// SetRate asks for readings no more often than every micros microseconds.
// It is a request: the platform delivers at that rate or slower, never
// faster, and a rate below the sensor's own minimum is clamped to it.
func (q *SensorQueue) SetRate(s Sensor, micros int) error {
	if C.ASensorEventQueue_setEventRate(q.ptr, s.ptr, C.int32_t(micros)) < 0 {
		return errors.New("ndk: cannot set the rate of " + s.Name())
	}
	return nil
}

// ErrNoBatching is what [SensorQueue.Register] reports on a platform too old
// to have it.
var ErrNoBatching = errors.New("ndk: this platform has no sensor batching; " +
	"it arrived in API 26")

// Register starts a sensor with a rate and a batching latency, both in
// microseconds.
//
// Latency is what makes batching: the platform is allowed to hold readings
// for that long and deliver them together, which lets the processor sleep in
// between. For a step counter running all day that is the difference between
// a percent of the battery and ten. A latency of 0 means deliver each one as
// it happens.
func (q *SensorQueue) Register(s Sensor, rateMicros int, latencyMicros int64) error {
	if C.mw_register(q.ptr, s.ptr, C.int32_t(rateMicros), C.int64_t(latencyMicros)) < 0 {
		return ErrNoBatching
	}
	return nil
}

// SensorEvent is one reading.
type SensorEvent struct {
	// Type is which kind of sensor it came from.
	Type int
	// ID identifies the sensor itself, which matters when two of a kind are
	// enabled.
	ID int
	// At is the reading's own timestamp, in nanoseconds on a clock that runs
	// while the device is awake. It is not the time the app read it, and the
	// difference is the whole point of batching.
	At int64
	// V is what was measured. The platform's array is sixteen long and it is
	// kept whole: three numbers for anything that measures a direction, one
	// for anything that measures a quantity, four or five for a rotation.
	V [16]float32
}

// Read takes up to len(into) readings, and reports how many there were.
func (q *SensorQueue) Read(into []SensorEvent) int {
	if len(into) == 0 {
		return 0
	}
	buf := make([]C.ASensorEvent, len(into))
	n := int(C.ASensorEventQueue_getEvents(q.ptr, &buf[0], C.size_t(len(into))))
	if n <= 0 {
		return 0
	}
	for i := range n {
		e := &buf[i]
		into[i].Type = int(C.mw_sensor_type(e))
		into[i].ID = int(C.mw_sensor_id(e))
		into[i].At = int64(C.mw_sensor_time(e))
		for j := range into[i].V {
			into[i].V[j] = float32(C.mw_sensor_data(e, C.int(j)))
		}
	}
	return n
}

// HasEvents reports whether anything is waiting, without taking it.
func (q *SensorQueue) HasEvents() bool {
	return C.ASensorEventQueue_hasEvents(q.ptr) > 0
}
