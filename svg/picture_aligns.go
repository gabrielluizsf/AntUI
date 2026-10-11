package svg

// pictureAligns is every alignment the attribute spells, lowercased because the
// values are matched case-insensitively the same way the names are.
var pictureAligns = map[string]PictureAlign{
	"xminymin": AlignXMinYMin, "xmidymin": AlignXMidYMin, "xmaxymin": AlignXMaxYMin,
	"xminymid": AlignXMinYMid, "xmidymid": AlignXMidYMid, "xmaxymid": AlignXMaxYMid,
	"xminymax": AlignXMinYMax, "xmidymax": AlignXMidYMax, "xmaxymax": AlignXMaxYMax,
}
