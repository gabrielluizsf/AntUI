package sdk

import (
	"strconv"
)

// String is the platform as a person would say it: "android-36".
func (p Platform) String() string {
	if p.Name != "" {
		return p.Name
	}
	return "android-" + strconv.Itoa(p.API)
}
