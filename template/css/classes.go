package css

import (
	"fmt"
	"os"
)

// Role names: the tag each widget kind carries into the cascade. A selector
// may match any of them, alone or with the classes the template gives them:
//
//	button                { background-color: #3E63DD; }
//	select:focus          { border-color: #3E63DD; }
//	textarea.wide         { width: 80%; }
//
// Roles are the template's tags, and NewTable pre-fills one key per widget
// kind so nothing in a stylesheet goes unheard.
const (
	RoleBody       = "body"
	RoleFlex       = "flex"
	RoleLabel      = "label"
	RoleButton     = "button"
	RoleCheckbox   = "checkbox"
	RoleRadio      = "radio"
	RoleSlider     = "slider"
	RoleInput      = "input"
	RoleSelect     = "select"
	RoleTextArea   = "textarea"
	RoleSwitch     = "switch"
	RoleProgress   = "progress"
	RoleDatePicker = "datepicker"
)

// AllRoles is the full set of widget tags, in a stable order.
var AllRoles = []string{
	RoleBody, RoleFlex, RoleLabel, RoleButton, RoleCheckbox, RoleRadio,
	RoleSlider, RoleInput, RoleSelect, RoleTextArea, RoleSwitch,
	RoleProgress, RoleDatePicker,
}

// CSSClasses is the class table a template draws from: one profile of CSS
// classes per widget kind. The template reads [CSSClasses.GetStyle] to paint
// and to lay out, and the table caches the computed styles so a screen that
// draws its widgets every frame computes each style once.
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
	Body, Flex, Label, Button, Checkbox, Radio, Slider, Input,
	Select, TextArea, Switch, Progress, DatePicker string

	sheet *Sheet
	specs map[styleKey]Style // computed styles, cleared on SetStyle/SetProperty
}

type styleKey struct {
	tag   string
	state State
	width int
}

// NewTable makes an empty class table, with no stylesheet behind it. Load the
// sheet with [CSSClasses.SetStyle] before drawing anything.
func NewTable() *CSSClasses {
	return &CSSClasses{specs: make(map[styleKey]Style)}
}

// SetStyle loads a CSS file into the table, replacing whatever it held, and
// clears the style cache. Compile errors — a malformed selector or an
// unbalanced block — stop the load; properties the engine does not know are
// skipped and reported through the sheet's warnings.
func (c *CSSClasses) SetStyle(cssFile string) error {
	if c.specs == nil {
		c.specs = make(map[styleKey]Style)
	}
	sh, err := ParseFile(cssFile)
	if err != nil {
		return err
	}
	c.sheet = sh
	c.specs = make(map[styleKey]Style)
	return nil
}

// ParseFile reads a CSS file, exactly as [Parse] reads its text.
func ParseFile(cssFile string) (*Sheet, error) {
	data, err := os.ReadFile(cssFile)
	if err != nil {
		return nil, err
	}
	return Parse(string(data))
}

// Sheet is the parsed stylesheet the table draws from.
func (c *CSSClasses) Sheet() *Sheet { return c.sheet }

// Keyframes returns the @keyframes block with the given name, the definition
// the stylesheet gave that an animation references, or nil.
func (c *CSSClasses) Keyframes(name string) *Keyframes {
	if c.sheet == nil {
		return nil
	}
	return c.sheet.Keyframes(name)
}

// Warnings are the sheet's own: every declaration the engine did not
// understand, for the template to show its author.
func (c *CSSClasses) Warnings() []string {
	if c.sheet == nil {
		return nil
	}
	return c.sheet.Warn
}

// Apply tags the table with a raw rule (a selector followed by a declaration
// block) folded into the sheet as if it were another rule later than the
// file, so it wins the cascade. Parse errors are reported and the rule is
// dropped; media queries inside are honoured.
func (c *CSSClasses) Apply(rule string) error {
	if c.sheet == nil {
		return errNoSheet
	}
	m := &parser{src: rule}
	sh := &Sheet{order: c.sheet.order}
	if err := m.parseRule(sh, Media{}); err != nil {
		return err
	}
	if len(sh.rules) == 0 {
		return errEmptyRule
	}
	c.sheet.rules = append(c.sheet.rules, sh.rules...)
	c.sheet.order = sh.order
	c.sheet.Warn = append(c.sheet.Warn, sh.Warn...)
	c.specs = make(map[styleKey]Style)
	return nil
}

// Property reads the winning raw value of one property for a tag and class
// list, computed the way the cascade would compute it. It answers what the
// stylesheet says, outside of the folded [Style] the drawing uses; width is
// the window's, which media queries test against.
func (c *CSSClasses) Property(tag string, classes []string, prop string, width int) (string, bool) {
	if c.sheet == nil {
		return "", false
	}
	return c.sheet.Property(tag, classes, prop, width)
}

// SetProperty overrides a property for every widget of the given tag and
// classes, appended after the sheet — so it wins whatever the file said for
// those selectors — and clears the cache.
func (c *CSSClasses) SetProperty(tag string, classes []string, prop, raw string) {
	if c.sheet == nil {
		return
	}
	c.sheet.SetProperty(tag, classes, prop, raw)
	c.specs = make(map[styleKey]Style)
}

// GetStyle computes the winning style for one widget: body for the window
// itself, or one of the Role tags. Classes are the widget's own, combined
// with the class table's entries for its kind. media queries test against
// width, and state carries the pseudo-classes the widget is in. The result is
// cached, so a template that paints every widget each frame computes each
// style once.
func (c *CSSClasses) GetStyle(tag string, classes []string, state State, width int) Style {
	if c.sheet == nil {
		return Style{}
	}
	all := append(append([]string(nil), c.classesOf(tag)...), classes...)
	key := styleKey{tag: tag, state: state, width: width}
	if st, ok := c.specs[key]; ok {
		return st
	}
	st := c.sheet.Style(tag, all, state, width)
	if c.specs != nil {
		c.specs[key] = st
	}
	return st
}

// classesOf is the class table's entry for a widget kind.
func (c *CSSClasses) classesOf(tag string) []string {
	var s string
	switch tag {
	case RoleBody:
		s = c.Body
	case RoleFlex:
		s = c.Flex
	case RoleLabel:
		s = c.Label
	case RoleButton:
		s = c.Button
	case RoleCheckbox:
		s = c.Checkbox
	case RoleRadio:
		s = c.Radio
	case RoleSlider:
		s = c.Slider
	case RoleInput:
		s = c.Input
	case RoleSelect:
		s = c.Select
	case RoleTextArea:
		s = c.TextArea
	case RoleSwitch:
		s = c.Switch
	case RoleProgress:
		s = c.Progress
	case RoleDatePicker:
		s = c.DatePicker
	}
	return ClassList(s)
}

var (
	errNoSheet   = fmt.Errorf("css: no stylesheet loaded")
	errEmptyRule = fmt.Errorf("css: empty rule")
)
