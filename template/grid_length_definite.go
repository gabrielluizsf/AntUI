package template

import "github.com/gabrielluizsf/antui/template/css"

func gridLengthDefinite(length css.Length) bool {
	return !length.Auto() && !length.None()
}
