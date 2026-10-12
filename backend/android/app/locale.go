//go:build android

package app

// Locale is the language and country joined the way a tag is written:
// "pt-BR", or just "pt" when there is no country.
func (c Config) Locale() string {
	switch {
	case c.Language == "":
		return ""
	case c.Country == "":
		return c.Language
	}
	return c.Language + "-" + c.Country
}
