//go:build android

package app

// RunOnUISync runs f on the UI thread and waits for it to finish.
//
// Called from the UI thread it simply calls f, which is what keeps it from
// waiting for itself. From anywhere else it blocks — so it must not be
// called from anything the UI thread is itself waiting on, which in this
// library means it must not be called from inside an activity callback.
func RunOnUISync(f func()) error {
	if OnUIThread() {
		f()
		return nil
	}
	done := make(chan struct{})
	if err := RunOnUI(func() {
		defer close(done)
		f()
	}); err != nil {
		return err
	}
	<-done
	return nil
}
