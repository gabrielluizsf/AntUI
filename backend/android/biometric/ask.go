//go:build android

package biometric

import (
	"fmt"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
)

// Ask puts the platform's prompt on screen and waits.
//
// It returns nil when the user proved who they are, and one of the errors
// above otherwise. A finger that was not recognised is **not** a failure —
// the prompt stays up and the user tries again — so this only returns when
// the matter is settled one way or the other.
//
// cancel is the label on the way out, and cannot be empty: the platform
// insists that the user always have one.
//
// **Not from the frame loop**, for the same reason as every other thing that
// waits for a person.
func Ask(title, subtitle, cancel string) error {
	if cancel == "" {
		cancel = "Cancel"
	}
	can, err := Can()
	if err != nil {
		return err
	}
	if can != Available {
		return fmt.Errorf("biometric: %s", can)
	}

	done := make(chan error, 1)
	var once sync.Once
	finish := func(err error) { once.Do(func() { done <- err }) }

	token, release := app.Listen(askListener(finish))
	defer release()

	if e := app.RunOnUISync(func() { err = authenticate(token, title, subtitle, cancel) }); e != nil {
		return e
	}
	if err != nil {
		return plainly(err)
	}
	return <-done
}
