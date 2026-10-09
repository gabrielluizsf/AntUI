package css

import (
	"fmt"
)

// classesOf is the class table's entry for a widget kind.
func (c *CSSClasses) classesOf(tag string) []string {
	var s string
	switch tag {
	case RoleBody:
		s = c.Body
	case RoleFlex:
		s = c.Flex
	case RoleGrid:
		s = c.Grid
	case RoleMultiCol:
		s = c.MultiCol
	case RoleTable:
		s = c.Table
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
