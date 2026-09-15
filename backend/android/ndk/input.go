//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

import "unsafe"

// InputQueue is the stream of key presses and touches the system has for
// this window. It only produces anything while it is attached to a looper,
// and only on the thread that looper belongs to.
type InputQueue struct {
	ptr *C.AInputQueue
}

// WrapQueue takes the pointer a native-activity callback was given.
func WrapQueue(p unsafe.Pointer) *InputQueue {
	if p == nil {
		return nil
	}
	return &InputQueue{ptr: (*C.AInputQueue)(p)}
}

// InputID is the identifier the input queue is attached to the looper with.
// It is arbitrary; it only has to be distinct from anything else attached.
const InputID = 1

// AttachLooper starts delivery to l. The callback argument the C function
// takes is deliberately null: a callback would be C calling into Go on a
// thread at a moment neither side chose, and there is nothing to gain from
// it when the frame loop is going to ask anyway.
func (q *InputQueue) AttachLooper(l *Looper) {
	C.AInputQueue_attachLooper(q.ptr, l.ptr, C.int(InputID), nil, nil)
}

// DetachLooper stops delivery. It must be called before the queue is
// destroyed, and from the looper's own thread.
func (q *InputQueue) DetachLooper() { C.AInputQueue_detachLooper(q.ptr) }

// HasEvents reports whether anything is waiting.
func (q *InputQueue) HasEvents() bool { return C.AInputQueue_hasEvents(q.ptr) > 0 }

// Next takes the next event, or reports false when there are none. Every
// event it hands out must be given back to [InputQueue.Finish] or
// [InputQueue.PreDispatch] — the system has a fixed number of them and an
// app that keeps one stops receiving input entirely.
func (q *InputQueue) Next() (*InputEvent, bool) {
	var ev *C.AInputEvent
	if C.AInputQueue_getEvent(q.ptr, &ev) < 0 {
		return nil, false
	}
	return &InputEvent{ptr: ev}, true
}

// PreDispatch offers the event to the system first, and reports whether the
// system took it. **Every event has to go through this before it is looked
// at**: it is how the on-screen keyboard sees a key at all, and an app that
// skips it appears to work until someone types.
//
// When it reports true the event has been consumed and must not be finished
// — the system will return it later if it did not want it after all.
func (q *InputQueue) PreDispatch(e *InputEvent) bool {
	return C.AInputQueue_preDispatchEvent(q.ptr, e.ptr) != 0
}

// Finish gives the event back, saying whether the app dealt with it. An
// unhandled key is passed on to the system, which is what makes the back
// button close the app when nothing else wanted it.
func (q *InputQueue) Finish(e *InputEvent, handled bool) {
	var h C.int
	if handled {
		h = 1
	}
	C.AInputQueue_finishEvent(q.ptr, e.ptr, h)
}

// InputEvent is one key press or one touch report. It is only valid until it
// is finished.
type InputEvent struct {
	ptr *C.AInputEvent
}

// The two kinds of event.
const (
	EventKey    = int(C.AINPUT_EVENT_TYPE_KEY)
	EventMotion = int(C.AINPUT_EVENT_TYPE_MOTION)
)

// Kind is whether this is a key or a motion.
func (e *InputEvent) Kind() int { return int(C.AInputEvent_getType(e.ptr)) }

// Source is what produced it: a touchscreen, a mouse, a joystick. The value
// is a set of bits, so it is tested with [Source.Is] rather than compared.
func (e *InputEvent) Source() Source { return Source(C.AInputEvent_getSource(e.ptr)) }

// Device is the input device's id, which is how two gamepads are told apart.
func (e *InputEvent) Device() int { return int(C.AInputEvent_getDeviceId(e.ptr)) }

// Source is where an event came from.
type Source int32

// The sources worth naming. Each is a class of device or'd with the kind of
// data it produces, which is why they overlap and are tested as bits.
const (
	SourceKeyboard    Source = C.AINPUT_SOURCE_KEYBOARD
	SourceDpad        Source = C.AINPUT_SOURCE_DPAD
	SourceGamepad     Source = C.AINPUT_SOURCE_GAMEPAD
	SourceTouchscreen Source = C.AINPUT_SOURCE_TOUCHSCREEN
	SourceMouse       Source = C.AINPUT_SOURCE_MOUSE
	SourceStylus      Source = C.AINPUT_SOURCE_STYLUS
	SourceTrackball   Source = C.AINPUT_SOURCE_TRACKBALL
	SourceTouchpad    Source = C.AINPUT_SOURCE_TOUCHPAD
	SourceJoystick    Source = C.AINPUT_SOURCE_JOYSTICK
)

// Is reports whether the source includes want. A touchscreen and a stylus
// both carry the touchscreen bits, so this is an "includes", not an "equals".
func (s Source) Is(want Source) bool { return s&want == want }

// Key events.

// The key actions.
const (
	KeyDown = int(C.AKEY_EVENT_ACTION_DOWN)
	KeyUp   = int(C.AKEY_EVENT_ACTION_UP)
)

// KeyAction is whether the key went down, came up, or is repeating.
func (e *InputEvent) KeyAction() int { return int(C.AKeyEvent_getAction(e.ptr)) }

// KeyCode is which key it was, as one of the AKEYCODE constants.
func (e *InputEvent) KeyCode() int { return int(C.AKeyEvent_getKeyCode(e.ptr)) }

// KeyScanCode is the hardware code behind the key, which differs by
// device and is almost never what a program wants.
func (e *InputEvent) KeyScanCode() int { return int(C.AKeyEvent_getScanCode(e.ptr)) }

// KeyRepeat is how many times the key has repeated. It is 0 on the first
// press, which is how a repeat is told from a new one.
func (e *InputEvent) KeyRepeat() int { return int(C.AKeyEvent_getRepeatCount(e.ptr)) }

// Meta is the modifier keys held, for a key or a motion.
func (e *InputEvent) Meta() Meta {
	if e.Kind() == EventKey {
		return Meta(C.AKeyEvent_getMetaState(e.ptr))
	}
	return Meta(C.AMotionEvent_getMetaState(e.ptr))
}

// Meta is a set of held modifiers.
type Meta int32

// The modifiers this library reads.
const (
	MetaShift Meta = C.AMETA_SHIFT_ON
	MetaAlt   Meta = C.AMETA_ALT_ON
	MetaCtrl  Meta = C.AMETA_CTRL_ON
	MetaSuper Meta = C.AMETA_META_ON
	MetaCaps  Meta = C.AMETA_CAPS_LOCK_ON
)

// Has reports whether every modifier in want is held.
func (m Meta) Has(want Meta) bool { return m&want == want }

// Motion events.

// The motion actions, after masking. DOWN and UP are the first finger going
// down and the last coming up; POINTER_DOWN and POINTER_UP are the ones in
// between, and they carry which finger in the top bits — see
// [InputEvent.MotionIndex].
const (
	MotionDown        = int(C.AMOTION_EVENT_ACTION_DOWN)
	MotionUp          = int(C.AMOTION_EVENT_ACTION_UP)
	MotionMove        = int(C.AMOTION_EVENT_ACTION_MOVE)
	MotionCancel      = int(C.AMOTION_EVENT_ACTION_CANCEL)
	MotionOutside     = int(C.AMOTION_EVENT_ACTION_OUTSIDE)
	MotionPointerDown = int(C.AMOTION_EVENT_ACTION_POINTER_DOWN)
	MotionPointerUp   = int(C.AMOTION_EVENT_ACTION_POINTER_UP)
	MotionHoverMove   = int(C.AMOTION_EVENT_ACTION_HOVER_MOVE)
	MotionScroll      = int(C.AMOTION_EVENT_ACTION_SCROLL)
	MotionHoverEnter  = int(C.AMOTION_EVENT_ACTION_HOVER_ENTER)
	MotionHoverExit   = int(C.AMOTION_EVENT_ACTION_HOVER_EXIT)
	MotionButtonPress = int(C.AMOTION_EVENT_ACTION_BUTTON_PRESS)
	MotionButtonUp    = int(C.AMOTION_EVENT_ACTION_BUTTON_RELEASE)
)

// MotionAction is what happened, with the pointer index masked off.
func (e *InputEvent) MotionAction() int {
	return int(C.AMotionEvent_getAction(e.ptr)) & C.AMOTION_EVENT_ACTION_MASK
}

// MotionIndex is *which* pointer a POINTER_DOWN or POINTER_UP is about, as
// an index into this event's pointers — **not** a pointer id.
//
// This is the single most common way to get multi-touch wrong. The index is
// a position in the list of pointers in this event and changes as fingers
// come and go; the id from [InputEvent.PointerID] is what stays with a
// finger for as long as it is down. Track ids; use indices only to read this
// event.
func (e *InputEvent) MotionIndex() int {
	a := int(C.AMotionEvent_getAction(e.ptr))
	return (a & C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >>
		C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
}

// PointerCount is how many pointers this event reports on.
func (e *InputEvent) PointerCount() int {
	return int(C.AMotionEvent_getPointerCount(e.ptr))
}

// PointerID is the identifier that stays with one finger from the moment it
// touches to the moment it lifts.
func (e *InputEvent) PointerID(index int) int {
	return int(C.AMotionEvent_getPointerId(e.ptr, C.size_t(index)))
}

// X and Y are in the window's pixels, and are floats because a touchscreen
// reports between pixels and a mouse on a scaled display does too.
func (e *InputEvent) X(index int) float32 {
	return float32(C.AMotionEvent_getX(e.ptr, C.size_t(index)))
}

// Y is where a pointer is, down the window, in pixels.
func (e *InputEvent) Y(index int) float32 {
	return float32(C.AMotionEvent_getY(e.ptr, C.size_t(index)))
}

// Pressure is 0 to 1 on a screen that measures it, and is 1 while down and 0
// while up on one that does not — so it cannot be used to tell whether the
// device reports pressure at all.
func (e *InputEvent) Pressure(index int) float32 {
	return float32(C.AMotionEvent_getPressure(e.ptr, C.size_t(index)))
}

// Size is how much of the screen the touch covers, normalised to 0..1. It is
// how a thumb is told from a fingertip.
func (e *InputEvent) Size(index int) float32 {
	return float32(C.AMotionEvent_getSize(e.ptr, C.size_t(index)))
}

// The tools a pointer can be.
const (
	ToolUnknown = int(C.AMOTION_EVENT_TOOL_TYPE_UNKNOWN)
	ToolFinger  = int(C.AMOTION_EVENT_TOOL_TYPE_FINGER)
	ToolStylus  = int(C.AMOTION_EVENT_TOOL_TYPE_STYLUS)
	ToolMouse   = int(C.AMOTION_EVENT_TOOL_TYPE_MOUSE)
	ToolEraser  = int(C.AMOTION_EVENT_TOOL_TYPE_ERASER)
)

// ToolType says what is touching: a finger, a stylus, a mouse, or the far
// end of a stylus being used as a rubber.
func (e *InputEvent) ToolType(index int) int {
	return int(C.AMotionEvent_getToolType(e.ptr, C.size_t(index)))
}

// The mouse buttons, as a set of bits.
const (
	ButtonPrimary   int32 = C.AMOTION_EVENT_BUTTON_PRIMARY
	ButtonSecondary int32 = C.AMOTION_EVENT_BUTTON_SECONDARY
	ButtonTertiary  int32 = C.AMOTION_EVENT_BUTTON_TERTIARY
	ButtonBack      int32 = C.AMOTION_EVENT_BUTTON_BACK
	ButtonForward   int32 = C.AMOTION_EVENT_BUTTON_FORWARD
)

// Buttons is which mouse buttons are down.
func (e *InputEvent) Buttons() int32 {
	return int32(C.AMotionEvent_getButtonState(e.ptr))
}

// The axes worth naming: the two scroll wheels, and the sticks, hat and
// triggers of a game controller.
const (
	AxisVScroll  = int(C.AMOTION_EVENT_AXIS_VSCROLL)
	AxisHScroll  = int(C.AMOTION_EVENT_AXIS_HSCROLL)
	AxisX        = int(C.AMOTION_EVENT_AXIS_X)
	AxisY        = int(C.AMOTION_EVENT_AXIS_Y)
	AxisZ        = int(C.AMOTION_EVENT_AXIS_Z)
	AxisRZ       = int(C.AMOTION_EVENT_AXIS_RZ)
	AxisHatX     = int(C.AMOTION_EVENT_AXIS_HAT_X)
	AxisHatY     = int(C.AMOTION_EVENT_AXIS_HAT_Y)
	AxisLTrigger = int(C.AMOTION_EVENT_AXIS_LTRIGGER)
	AxisRTrigger = int(C.AMOTION_EVENT_AXIS_RTRIGGER)
	AxisGas      = int(C.AMOTION_EVENT_AXIS_GAS)
	AxisBrake    = int(C.AMOTION_EVENT_AXIS_BRAKE)
)

// Axis reads one axis of one pointer. Everything a game controller reports
// comes through here rather than through X and Y.
func (e *InputEvent) Axis(axis, index int) float32 {
	return float32(C.AMotionEvent_getAxisValue(e.ptr, C.int32_t(axis), C.size_t(index)))
}

// EventTime is when it happened, in nanoseconds on the same clock as
// SystemClock.uptimeMillis. It is what a fling's speed is measured against.
func (e *InputEvent) EventTime() int64 {
	return int64(C.AMotionEvent_getEventTime(e.ptr))
}
