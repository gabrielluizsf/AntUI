package apk

import (
	"fmt"
	"strings"
)

// validPackage checks the application id against Android's rules, which are
// stricter than they look: at least two parts, each starting with a letter,
// and no Java keyword anywhere — the id becomes a package name in generated
// code, so "com.example.class" is refused by the platform, not by us.
func validPackage(p string) error {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return fmt.Errorf("apk: package %q needs at least two parts, like \"com.example.app\"", p)
	}
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("apk: package %q has an empty part", p)
		}
		if !isLetter(part[0]) && part[0] != '_' {
			return fmt.Errorf("apk: package part %q must start with a letter", part)
		}
		for i := range len(part) {
			ch := part[i]
			if !isLetter(ch) && !isDigit(ch) && ch != '_' {
				return fmt.Errorf("apk: package part %q may only hold letters, digits and _", part)
			}
		}
		if javaKeywords[part] {
			return fmt.Errorf("apk: package part %q is a Java keyword, which Android refuses", part)
		}
	}
	return nil
}
