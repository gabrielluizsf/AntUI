package template

import "github.com/gabrielluizsf/antui/template/css"

func prependGridTracks(tracks *[]css.GridTrack, count int) {
	if count <= 0 {
		return
	}
	prefix := make([]css.GridTrack, count)
	for i := range prefix {
		prefix[i] = css.GridTrack{Kind: css.GridTrackSingle, Size: css.GridTrackSize{Kind: css.GridTrackAuto}}
	}
	*tracks = append(prefix, *tracks...)
}
