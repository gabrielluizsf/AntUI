package css

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
	RoleGrid       = "grid"
	RoleMultiCol   = "multicolumn"
	RoleTable      = "table"
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
	RoleBody, RoleFlex, RoleGrid, RoleMultiCol, RoleTable, RoleLabel,
	RoleButton, RoleCheckbox, RoleRadio, RoleSlider, RoleInput, RoleSelect,
	RoleTextArea, RoleSwitch, RoleProgress, RoleDatePicker,
}
