package apk

// Feature is a uses-feature entry.
type Feature struct {
	Name     string `json:"name"`
	Required bool   `json:"required,omitempty"`
}
