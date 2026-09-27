package linux

import (
	"testing"
	"time"
)

// TestIdleWaitSitsStillInsteadOfSpinning is the property the whole blocking
// pump rests on: with nothing to do the driver waits, and it waits about as
// long as it takes for the frame the program asked for. A driver that returned
// zero here would go round the loop as fast as the machine can manage, and a
// program waiting for a 30 Hz frame would burn a core drawing nothing.
func TestIdleWaitSitsStillInsteadOfSpinning(t *testing.T) {
	x := &Driver{}

	// A program that has not said when its next frame is due: wait the
	// backstop, not nothing.
	if got := x.idleWait(); got != maxIdleWait {
		t.Errorf("with no deadline set the driver waits %v, want %v", got, maxIdleWait)
	}

	// A program pacing at 60 Hz: wait until its frame, capped at the backstop.
	for _, wait := range []time.Duration{0, 4 * time.Millisecond, 16 * time.Millisecond, time.Second} {
		x.SetFrameDeadline(time.Now().Add(wait))
		got := x.idleWait()
		if got > maxIdleWait {
			t.Errorf("a frame due in %v: the driver waits %v, never more than the %v backstop",
				wait, got, maxIdleWait)
		}
		if got < 0 {
			t.Errorf("a frame due in %v: the driver waits %v, never backwards", wait, got)
		}
		if want := min(wait, maxIdleWait); got > want || want-got > 2*time.Millisecond {
			t.Errorf("a frame due in %v: the driver waits %v, want it to wake about when the frame is", wait, got)
		}
	}

	// A frame that is already late is not waited on at all: a window that
	// cannot keep up should run flat out, not idle between frames it is late for.
	x.SetFrameDeadline(time.Now().Add(-time.Second))
	if got := x.idleWait(); got != 0 {
		t.Errorf("a frame that is already due: the driver waits %v, want no wait at all", got)
	}
}

// TestSetFrameDeadlineZeroMeansNoPacing checks the other half of the deadline:
// clearing it puts the driver back on its own short wait, which is what a
// program that asked for no frame rate gets.
func TestSetFrameDeadlineZeroMeansNoPacing(t *testing.T) {
	x := &Driver{}
	x.SetFrameDeadline(time.Now())
	if got := x.idleWait(); got > maxIdleWait {
		t.Fatalf("a driver with a deadline waits %v, over the backstop", got)
	}
	x.SetFrameDeadline(time.Time{})
	if got := x.idleWait(); got != maxIdleWait {
		t.Errorf("after clearing the deadline the driver waits %v, want the %v backstop", got, maxIdleWait)
	}
}
