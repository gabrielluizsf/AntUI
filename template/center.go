package template

// Center is the left edge of a widget of the given width, centered
// horizontally — where a button goes to sit in the middle of the window.
func (c *Context) Center(w int) int { return (c.win.Width() - w) / 2 }
