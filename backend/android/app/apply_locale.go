//go:build android

package app

/*
#include <android/configuration.h>
*/
import "C"

// applyLocale reads the language and country.
//
// They are written into two chars each and are not terminated, so they
// cannot be read as C strings. An unset one comes back as two zero bytes.
func applyLocale(cfg *C.AConfiguration, out *Config) {
	var lang, country [2]C.char
	C.AConfiguration_getLanguage(cfg, &lang[0])
	C.AConfiguration_getCountry(cfg, &country[0])
	out.Language = twoLetters(lang)
	out.Country = twoLetters(country)
}
