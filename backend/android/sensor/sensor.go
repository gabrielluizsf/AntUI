//go:build android

package sensor

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Kind is a kind of sensor, numbered as the platform numbers it.
type Kind int

// The sensors worth naming. The three groups are worth telling apart.
//
// The first measure something directly and report it raw, noise and all.
// The second are computed by the platform out of the first — a rotation
// vector is the accelerometer, the gyroscope and the compass fused, and is
// steadier than any of them and slower to react. The third report an event
// rather than a value.
const (
	// Measured directly.
	Accelerometer Kind = 1  // metres per second squared, gravity included
	MagneticField Kind = 2  // microtesla
	Gyroscope     Kind = 4  // radians per second
	Light         Kind = 5  // lux
	Pressure      Kind = 6  // hectopascal
	Proximity     Kind = 8  // centimetres, and on most phones only 0 or "far"
	Humidity      Kind = 12 // per cent
	Temperature   Kind = 13 // degrees celsius

	// Computed by the platform.
	Gravity            Kind = 9  // which way is down, without the movement
	LinearAcceleration Kind = 10 // the movement, without gravity
	RotationVector     Kind = 11 // orientation as a quaternion, using the compass
	GameRotationVector Kind = 15 // the same without the compass: no true north,
	// and no lurching when a magnet goes past
	GeomagneticRotation Kind = 20 // the same using the compass and no gyroscope

	// Events rather than values.
	StepDetector      Kind = 18 // one reading per step
	StepCounter       Kind = 19 // steps since the device booted, in V[0]
	SignificantMotion Kind = 17
)

var kindNames = map[Kind]string{
	Accelerometer: "accelerometer", MagneticField: "magnetic field",
	Gyroscope: "gyroscope", Light: "light", Pressure: "pressure",
	Proximity: "proximity", Humidity: "humidity", Temperature: "temperature",
	Gravity: "gravity", LinearAcceleration: "linear acceleration",
	RotationVector: "rotation vector", GameRotationVector: "game rotation vector",
	GeomagneticRotation: "geomagnetic rotation", StepDetector: "step detector",
	StepCounter: "step counter", SignificantMotion: "significant motion",
}

// String names the sensor, and says the number for one this package has
// no name for.
func (k Kind) String() string {
	if n, ok := kindNames[k]; ok {
		return n
	}
	return fmt.Sprintf("sensor(%d)", int(k))
}

// Reading is one measurement.
type Reading struct {
	Kind Kind
	// At is when the sensor took it, not when the app read it — on a clock
	// that runs only while the device is awake. The difference between the
	// two is the whole point of batching.
	At time.Duration
	// V is what was measured, kept the full width the platform reports.
	// Three numbers for anything with a direction, one for anything with a
	// quantity, four or five for a rotation.
	V [16]float32
}

// X, Y and Z are the first three numbers, for a sensor that measures a
// direction. On a phone held upright facing the user, X points right, Y
// points up and Z points out of the screen — and that is the **device's**
// frame, which does not turn when the screen does. See [Reading.ForDisplay].
func (r Reading) X() float32 { return r.V[0] }

// Y is the second number. See [Reading.X] for which way it points.
func (r Reading) Y() float32 { return r.V[1] }

// Z is the third number. See [Reading.X] for which way it points.
func (r Reading) Z() float32 { return r.V[2] }

// Value is the first number, for a sensor that measures a quantity: lux,
// hectopascal, degrees, or a count of steps.
func (r Reading) Value() float32 { return r.V[0] }

// ForDisplay turns a reading from the device's frame into the one the app is
// drawing in.
//
// This is the thing everyone gets wrong once. A sensor reports in a frame
// fixed to the *hardware*: X across the device as it was built. The screen
// rotates and the sensor does not, so on a phone held sideways the
// accelerometer's X runs down the screen and a ball rolls the wrong way.
// The rotation to pass is [antui/backend/android/display.Rotation].
func (r Reading) ForDisplay(rotation int) Reading {
	r.V[0], r.V[1] = remap(r.V[0], r.V[1], rotation)
	return r
}

// Info describes a sensor the device has.
type Info struct {
	Kind       Kind
	Name       string
	Vendor     string
	Resolution float32
	// MinDelay is the shortest interval between readings. It is zero for a
	// sensor that reports only when something changes, and asking one of
	// those for a rate does nothing.
	MinDelay time.Duration
}

var (
	managerOnce sync.Once
	manager     *ndk.SensorManager
	managerErr  error
)

func sensors() (*ndk.SensorManager, error) {
	managerOnce.Do(func() {
		name := ""
		if a := app.Current(); a != nil {
			// The package name, which the platform uses from API 26 to
			// decide which sensors this app may see.
			name = packageName()
		}
		manager, managerErr = ndk.Sensors(name)
	})
	return manager, managerErr
}

// List is every sensor the device has.
func List() ([]Info, error) {
	m, err := sensors()
	if err != nil {
		return nil, err
	}
	all := m.List()
	out := make([]Info, 0, len(all))
	for _, s := range all {
		out = append(out, Info{
			Kind:       Kind(s.Type()),
			Name:       s.Name(),
			Vendor:     s.Vendor(),
			Resolution: s.Resolution(),
			MinDelay:   time.Duration(s.MinDelay()) * time.Microsecond,
		})
	}
	return out, nil
}

// Has reports whether the device has a sensor of a kind.
func Has(k Kind) bool {
	m, err := sensors()
	if err != nil {
		return false
	}
	return !m.Default(int(k)).IsNil()
}

// Options is how a stream is opened.
type Options struct {
	// Rate is how often to report. Zero asks for the sensor's own default.
	// It is a ceiling, not a promise: the platform delivers at that rate or
	// slower, and never faster than the sensor's own minimum.
	Rate time.Duration
	// Latency lets the platform hold readings and deliver them in a burst,
	// which is what lets the processor sleep between them. Zero delivers each
	// as it happens. It needs API 26; below that it is ignored.
	Latency time.Duration
	// Wake asks for a sensor that keeps the device awake to report. Almost
	// nothing should: a step counter that does not wake the device reports
	// everything it collected when the device next wakes, which is the same
	// answer for a fraction of the battery.
	Wake bool
	// Buffer is how many readings [Stream.Events] keeps when nothing is
	// reading them. The oldest are dropped. Zero means no channel at all,
	// and only [Stream.Latest].
	Buffer int
}

// Stream is a set of running sensors.
//
// It owns a thread of its own with a looper of its own, because that is what
// the platform's queue needs and because sensors run far faster than a frame
// loop — an accelerometer at its fastest is 200 readings a second, and
// pulling those through the frame loop would tie the two together for no
// reason.
type Stream struct {
	events chan Reading

	mu     sync.Mutex
	latest map[Kind]Reading
	fresh  map[Kind]bool

	stop     chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once
}

// Open starts the named sensors.
func Open(kinds []Kind, opt Options) (*Stream, error) {
	if len(kinds) == 0 {
		return nil, errors.New("sensor: no sensors asked for")
	}
	m, err := sensors()
	if err != nil {
		return nil, err
	}
	chosen := make([]ndk.Sensor, 0, len(kinds))
	for _, k := range kinds {
		s := m.Default(int(k))
		if opt.Wake {
			s = m.DefaultWake(int(k), true)
		}
		if s.IsNil() {
			return nil, fmt.Errorf("sensor: this device has no %s", k)
		}
		chosen = append(chosen, s)
	}

	st := &Stream{
		latest:  make(map[Kind]Reading, len(kinds)),
		fresh:   make(map[Kind]bool, len(kinds)),
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	if opt.Buffer > 0 {
		st.events = make(chan Reading, opt.Buffer)
	}

	ready := make(chan error, 1)
	go st.run(m, chosen, opt, ready)
	if err := <-ready; err != nil {
		return nil, err
	}
	return st, nil
}

// run is the stream's own thread: it makes the looper, the queue and the
// subscriptions, and stays there until Close.
func (s *Stream) run(m *ndk.SensorManager, chosen []ndk.Sensor, opt Options, ready chan<- error) {
	// A looper belongs to its thread, and so does everything delivered to
	// it, so the goroutine may never move.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(s.stopped)

	looper := ndk.PrepareLooper()
	if looper == nil {
		ready <- errors.New("sensor: cannot make a looper for the sensor thread")
		return
	}
	queue, err := m.Queue(looper)
	if err != nil {
		ready <- err
		return
	}
	defer queue.Close()

	rate := int(opt.Rate.Microseconds())
	latency := opt.Latency.Microseconds()
	for _, sn := range chosen {
		// Batching first, and the plain path when the platform is too old
		// for it or the caller asked for none.
		if latency > 0 {
			if err := queue.Register(sn, rate, latency); err == nil {
				continue
			}
		}
		if err := queue.Enable(sn); err != nil {
			ready <- err
			return
		}
		if rate > 0 {
			if err := queue.SetRate(sn, rate); err != nil {
				ready <- err
				return
			}
		}
	}
	ready <- nil

	buf := make([]ndk.SensorEvent, 32)
	for {
		select {
		case <-s.stop:
			return
		default:
		}
		// A timeout rather than a wait, so that a Close that races with the
		// poll is noticed even if the wake is missed.
		looper.Poll(250)
		for {
			n := queue.Read(buf)
			if n == 0 {
				break
			}
			for _, e := range buf[:n] {
				s.deliver(Reading{
					Kind: Kind(e.Type),
					At:   time.Duration(e.At),
					V:    e.V,
				})
			}
		}
	}
}

func (s *Stream) deliver(r Reading) {
	s.mu.Lock()
	s.latest[r.Kind] = r
	s.fresh[r.Kind] = true
	s.mu.Unlock()
	if s.events != nil {
		select {
		case s.events <- r:
		default:
			// The buffer is full and nobody is reading. Dropping the newest
			// would be wrong — for a sensor the newest is the one that
			// matters — so the oldest goes.
			select {
			case <-s.events:
			default:
			}
			select {
			case s.events <- r:
			default:
			}
		}
	}
}

// Latest is the most recent reading from one sensor, and whether there has
// been a new one since the last call. It is what a frame loop wants: the
// current value, with no queue to drain.
func (s *Stream) Latest(k Kind) (Reading, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.latest[k]
	fresh := s.fresh[k]
	s.fresh[k] = false
	return r, ok && fresh
}

// Current is the most recent reading, whether or not it is new.
func (s *Stream) Current(k Kind) (Reading, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.latest[k]
	return r, ok
}

// Events is every reading, for the sensors that report an event rather than
// a value — a step, a significant motion. It is nil unless [Options.Buffer]
// was set.
func (s *Stream) Events() <-chan Reading { return s.events }

// Close stops every sensor and gives the thread back. It waits for the
// thread to finish, so that nothing is still writing when it returns.
func (s *Stream) Close() error {
	s.stopOnce.Do(func() {
		close(s.stop)
		<-s.stopped
	})
	return nil
}
