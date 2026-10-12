//go:build android

package app

// Info is what the activity knows about itself.
//
// A nil App answers with a zero Info rather than crashing. [Current] returns
// nil while there is no activity — before one exists, and for the moment
// after one is destroyed, which is what a rotation does — and a goroutine
// doing work in the background will call this in that window. Taking the
// whole process down for it, in a library whose job is to make a phone
// bearable, is not the right answer; an empty path that fails the next
// operation with a reason is.
func (a *App) Info() Info {
	if a == nil {
		return Info{}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.info
}
