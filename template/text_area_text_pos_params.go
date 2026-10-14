package template

// TextAreaTextPosParams holds parameters for hit-testing a multi-line text area.
type TextAreaTextPosParams struct {
	TextAreaParams
	MouseX, MouseY int
}