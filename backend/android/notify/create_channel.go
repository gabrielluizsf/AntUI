//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// CreateChannel makes a channel, or does nothing if it already exists.
// Calling it every launch is the intended use: it is cheap and it is how a
// name or a description gets updated.
//
// Below API 26 there are no channels and this does nothing successfully.
func CreateChannel(c Channel) error {
	if c.ID == "" {
		return fmt.Errorf("notify: a channel needs an id")
	}
	a := app.Current()
	if a == nil {
		return app.ErrNoUIThread
	}
	if a.Config().SDK < 26 {
		return nil
	}
	if c.Importance == 0 {
		c.Importance = Default
	}
	return jni.Do(func(e *jni.Env) error {
		return makeChannel(e, c)
	})
}
