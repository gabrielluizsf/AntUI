//go:build android

package antui

// familySet is as much of fonts.xml as this needs.
type familySet struct {
	Families []struct {
		Name  string `xml:"name,attr"`
		Fonts []struct {
			Weight string `xml:"weight,attr"`
			Style  string `xml:"style,attr"`
			File   string `xml:",chardata"`
		} `xml:"font"`
	} `xml:"family"`
}
