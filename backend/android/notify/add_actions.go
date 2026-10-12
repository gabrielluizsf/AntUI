//go:build android

package notify

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// addActions puts the buttons on a notification being built.
func addActions(e *jni.Env, builder jni.Object, n Notification) error {
	if len(n.Actions) == 0 {
		return nil
	}
	actionMu.Lock()
	token := actionTok
	actionMu.Unlock()
	if token == 0 {
		return fmt.Errorf("notify: this notification has actions and none has " +
			"been registered with OnAction, so a tap would go nowhere")
	}

	for i, a := range n.Actions {
		if err := addAction(e, builder, n.ID, int32(i), a, token); err != nil {
			return err
		}
	}
	return nil
}
