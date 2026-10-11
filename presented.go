package antui

// Presented reports whether the last [Window.End] put anything on the screen.
// A template can read it to see whether the frame it just drew was a change
// or a repeat: the input state is settled either way, and a repeat is free.
func (win *Window) Presented() bool { return win.presented }
