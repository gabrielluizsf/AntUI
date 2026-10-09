package css

// CSSClasses is the class table a template draws from: one profile of CSS
// classes per widget kind. The template reads [CSSClasses.GetStyleViewport]
// to paint and to lay out, and the table caches the computed styles — keyed
// by the window they were read against — so a screen that draws its widgets
// every frame computes each style once.
//
//	classes := NewTable()
//	classes.Button = "primary"
//	classes.Input = "wide"
//
// A stylesheet then styles exactly those classes:
//
//	.primary { background-color: #3E63DD; }
//	.wide    { width: 80%; }
type CSSClasses struct {
	// One class list per widget kind. Empty strings mean the widget carries
	// no classes, only its tag.
	Body, Flex, Grid, MultiCol, Table, Label, Button, Checkbox, Radio, Slider,
	Input, Select, TextArea, Switch, Progress, DatePicker string

	sheet *Sheet
	specs map[styleKey]Style // computed styles, cleared on SetStyle/SetProperty
	vp    Viewport           // the window those styles were computed for
}
