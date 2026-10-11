package template

import "github.com/gabrielluizsf/antui/template/css"

func ensureGridTracks(tracks *[]css.GridTrack, count int) {
	for len(*tracks) < count {
		*tracks = append(*tracks, css.GridTrack{Kind: css.GridTrackSingle, Size: css.GridTrackSize{Kind: css.GridTrackAuto}})
	}
}
