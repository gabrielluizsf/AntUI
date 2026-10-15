package template

import "github.com/gabrielluizsf/antui"

// Context is one window running one template through a stack of screens. Screens
// push one another with Go, come back with Back, and the Context keeps the
// record — the stack IS the map home.
//
//	win, _ := antui.OpenWith(antui.Options{Title: "Shop", Width: 480, Height: 320})
//
//	app := template.NewContext(win, template.Cyberpunk())
//	app.Run(Home{})
//
// while Home's Draw decides when to leave:
//
//	func (Home) Draw(c *Context) {
//	    c.Label(0, 0, "Home")
//	    if c.Button(canvas.Area{X: c.Center(160), Y: 120, Width: 160, Height: 40}, "To Shop").Is(event.Button, event.Click) {
//	        c.Go(Shop{})
//	    }
//	}
type Context struct {
	win   *antui.Window
	tpl   Template
	stack []Screen
}
