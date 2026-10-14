package template

import (
	"github.com/gabrielluizsf/antui"
)

// WidgetGeometry is the common geometry for all widget painting methods.
type WidgetGeometry struct {
	Win   *antui.Window
	State State
	X, Y  int
	W, H  int
}

// LabelParams holds parameters for painting a label.
type LabelParams struct {
	WidgetGeometry
	Text string
}

// ButtonParams holds parameters for painting a button.
type ButtonParams struct {
	WidgetGeometry
	Label string
}

// CheckboxParams holds parameters for painting a checkbox.
type CheckboxParams struct {
	WidgetGeometry
	Label string
	On    bool
}

// RadioParams holds parameters for painting a radio option.
type RadioParams struct {
	WidgetGeometry
	Label string
	On    bool
}

// SliderParams holds parameters for painting a slider.
type SliderParams struct {
	WidgetGeometry
	Value float32
}

// InputParams holds parameters for painting a single-line text input.
type InputParams struct {
	WidgetGeometry
	Text   string
	Cursor int
	Caret  bool
}

// InputTextPosParams holds parameters for hit-testing a single-line text input.
type InputTextPosParams struct {
	InputParams
	MouseX, MouseY int
}

// TextAreaParams holds parameters for painting a multi-line text area.
type TextAreaParams struct {
	WidgetGeometry
	Text   string
	Cursor int
	Caret  bool
}

// TextAreaTextPosParams holds parameters for hit-testing a multi-line text area.
type TextAreaTextPosParams struct {
	TextAreaParams
	MouseX, MouseY int
}

// SelectParams holds parameters for painting a select dropdown.
type SelectParams struct {
	WidgetGeometry
	Value string
	Open  bool
}

// SelectOptionParams holds parameters for painting a select option.
type SelectOptionParams struct {
	WidgetGeometry
	Label    string
	Selected bool
}

// SwitchParams holds parameters for painting a switch toggle.
type SwitchParams struct {
	WidgetGeometry
	Label string
	On    bool
}

// ProgressParams holds parameters for painting a progress bar.
type ProgressParams struct {
	WidgetGeometry
	Value float32
}

// DatePickerParams holds parameters for painting a date picker calendar.
type DatePickerParams struct {
	WidgetGeometry
	Year, Month, FirstWD, Days, Selected, Today, Hover int
}

// DatePickerBoxParams holds parameters for painting a closed date picker.
type DatePickerBoxParams struct {
	WidgetGeometry
	Value string
	Open  bool
}