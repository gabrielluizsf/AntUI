package css

import "testing"

func gridRule(t *testing.T, role, rule string) Style {
	t.Helper()
	sh, err := Parse(role + " {" + rule + "}")
	if err != nil {
		t.Fatal(err)
	}
	return sh.Style(role, nil, StateNone, 800)
}

func TestGridTemplateParsesTracksAndRepeat(t *testing.T) {
	tracks, ok := parseGridTemplate("100px 1fr repeat(2, minmax(80px, 2fr))", Units{})
	if !ok || len(tracks) != 4 {
		t.Fatalf("tracks = %d, %v want 4,true", len(tracks), ok)
	}
	first := tracks[0].Track
	if first.Kind != GridTrackSingle || first.Size.Kind != GridTrackLength || first.Size.Length.Resolve(Units{}) != 100 {
		t.Errorf("first track = %+v", first)
	}
	if tracks[1].Track.Size.Kind != GridTrackFlexible || tracks[1].Track.Size.Fr != 1 {
		t.Errorf("second track = %+v", tracks[1].Track)
	}
	for _, segment := range tracks[2:] {
		track := segment.Track
		if track.Kind != GridTrackMinMax || track.Min.Kind != GridTrackLength || track.Min.Length.Resolve(Units{}) != 80 || track.Max.Kind != GridTrackFlexible || track.Max.Fr != 2 {
			t.Errorf("repeated track = %+v", track)
		}
	}
}

func TestGridTemplateParsesLineNamesCaseSensitively(t *testing.T) {
	tracks, ok := parseGridTemplate("[Full-Start] 120px [content-start] 1fr [content-start]", Units{})
	if !ok || len(tracks) != 2 {
		t.Fatalf("tracks = %d, %v want 2,true", len(tracks), ok)
	}
	axis := tracks.Resolve(400, 0, 2)
	if len(axis.Tracks) != 2 || len(axis.Names) != 3 {
		t.Fatalf("axis = %d tracks, %d lines", len(axis.Tracks), len(axis.Names))
	}
	if line, ok := axis.LineIndex("Full-Start", false); !ok || line != 0 {
		t.Errorf("Full-Start = %d, %v want 0,true", line, ok)
	}
	if line, ok := axis.LineIndex("content-start", false); !ok || line != 1 {
		t.Errorf("content-start = %d, %v want 1,true", line, ok)
	}
	if line, ok := axis.LineIndex("content-start", true); !ok || line != 2 {
		t.Errorf("the last content-start = %d, %v want 2,true", line, ok)
	}
	if _, ok := axis.LineIndex("full-start", false); ok {
		t.Error("line names are case-sensitive")
	}
	if _, ok := axis.LineIndex("full-end", true); ok {
		t.Error("a name no line carries must not be found")
	}
}

func TestGridTemplateNamesRideRepeat(t *testing.T) {
	tracks, ok := parseGridTemplate("repeat(2, [edge-start] 1fr [edge-end])", Units{})
	if !ok || len(tracks) != 2 {
		t.Fatalf("tracks = %d, %v want 2,true", len(tracks), ok)
	}
	axis := tracks.Resolve(0, 0, 1)
	if line, ok := axis.LineIndex("edge-start", false); !ok || line != 0 {
		t.Errorf("edge-start = %d, %v want 0,true", line, ok)
	}
	if line, ok := axis.LineIndex("edge-end", true); !ok || line != 2 {
		t.Errorf("edge-end = %d, %v want 2,true", line, ok)
	}
}

func TestGridTemplateKeepsAutoRepeatForTheLayout(t *testing.T) {
	for _, raw := range []string{
		"repeat(auto-fill, minmax(200px, 1fr))",
		"repeat(auto-fit, minmax(200px, 1fr))",
		"100px repeat(auto-fill, 80px) 1fr",
		"[start] repeat(auto-fill, minmax(10em, 1fr)) [end]",
	} {
		tracks, ok := parseGridTemplate(raw, Units{})
		if !ok {
			t.Fatalf("%q must parse", raw)
		}
		var autos int
		for _, segment := range tracks {
			if segment.Auto {
				autos++
			}
		}
		if autos != 1 {
			t.Errorf("%q = %d auto segments, want 1", raw, autos)
		}
	}
	// auto-fit asks for the same minimum as auto-fill, so the count the layout
	// works out differs only in the repetitions nothing lands in.
	fill, _ := parseGridTemplate("repeat(auto-fill, minmax(200px, 1fr))", Units{})
	fit, _ := parseGridTemplate("repeat(auto-fit, minmax(200px, 1fr))", Units{})
	if len(fill.Resolve(650, 10, 2).Tracks) != 3 || len(fit.Resolve(650, 10, 2).Tracks) != 2 {
		t.Errorf("650px with two items = %d fill, %d fit, want 3 and 2",
			len(fill.Resolve(650, 10, 2).Tracks), len(fit.Resolve(650, 10, 2).Tracks))
	}
	if len(fill.Resolve(650, 10, 6).Tracks) != 3 || len(fit.Resolve(650, 10, 6).Tracks) != 3 {
		t.Errorf("650px with six items = %d fill, %d fit, want 3 and 3",
			len(fill.Resolve(650, 10, 6).Tracks), len(fit.Resolve(650, 10, 6).Tracks))
	}
	if len(fill.Resolve(0, 0, 3).Tracks) != 1 || len(fill.Resolve(-40, 0, 3).Tracks) != 1 {
		t.Error("an auto-repeat never falls below one track")
	}
}

func TestGridTemplateRejectsInvalidTracks(t *testing.T) {
	for _, raw := range []string{
		"", "subgrid", "repeat(0, 1fr)", "repeat(auto-fit, 1fr)", "repeat(auto-fill, 1fr)",
		"repeat(auto-fit, auto)", "repeat(auto-fit, 200px 1fr)", "repeat(auto-fit, [a] 1fr [b])",
		"repeat(auto-fit)", "repeat(2)", "repeat(2,)", "2,,1fr", "minmax(1fr, 1fr)",
		"minmax(fit-content(10px), 1fr)", "minmax(100px,,1fr)", "-1fr", "0fr", "1fr 2fr garbage",
		"[", "[]", "[auto]", "1fr [a", "repeat(2, 1fr", "fit-content(-1px)", "fit-content(1fr)",
		"repeat(1001, 1fr)",
	} {
		if _, ok := parseGridTemplate(raw, Units{}); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestGridTemplateParsesContentKeywords(t *testing.T) {
	tracks, ok := parseGridTemplate("min-content max-content fit-content(120px) minmax(min-content, max-content)", Units{})
	if !ok || len(tracks) != 4 {
		t.Fatalf("tracks = %d, %v want 4,true", len(tracks), ok)
	}
	want := []GridTrackSizeKind{GridTrackMinContent, GridTrackMaxContent, GridTrackFitContent}
	for i, kind := range want {
		if tracks[i].Track.Size.Kind != kind {
			t.Errorf("track %d = %+v, want kind %d", i, tracks[i].Track.Size, kind)
		}
	}
	if tracks[2].Track.Size.Length.Resolve(Units{}) != 120 {
		t.Errorf("fit-content limit = %v, want 120", tracks[2].Track.Size.Length)
	}
	pair := tracks[3].Track
	if pair.Kind != GridTrackMinMax || pair.Min.Kind != GridTrackMinContent || pair.Max.Kind != GridTrackMaxContent {
		t.Errorf("minmax of content keywords = %+v", pair)
	}
}

func TestGridTemplateAreasParse(t *testing.T) {
	areas, ok := parseGridTemplateAreas(`"header header" "sidebar main"`)
	if !ok || len(areas) != 2 || len(areas[0]) != 2 || len(areas[1]) != 2 {
		t.Fatalf("areas = %#v, %v", areas, ok)
	}
	if areas[0][0] != "header" || areas[1][1] != "main" {
		t.Errorf("areas = %#v", areas)
	}
	caseAreas, ok := parseGridTemplateAreas(`"Foo foo"`)
	if !ok || caseAreas[0][0] != "Foo" || caseAreas[0][1] != "foo" {
		t.Errorf("case-sensitive areas = %#v, %v", caseAreas, ok)
	}
	single, ok := parseGridTemplateAreas(`'header header' 'side main'`)
	if !ok || len(single) != 2 || single[0][0] != "header" || single[1][1] != "main" {
		t.Errorf("single-quoted areas = %#v, %v", single, ok)
	}
	null, ok := parseGridTemplateAreas(`"a ..."`)
	if !ok || len(null) != 1 || len(null[0]) != 2 || null[0][1] != "..." {
		t.Errorf("null areas = %#v, %v", null, ok)
	}
	for _, raw := range []string{`"one" "two three"`, `header main`, `"unclosed`, `"a," "b"`, `"a ." ". a"`} {
		if _, ok := parseGridTemplateAreas(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestGridAutoFlowParsesDirectionAndDense(t *testing.T) {
	for _, tc := range []struct {
		raw   string
		flow  uint8
		dense bool
	}{
		{"row", GridAutoFlowRow, false},
		{"column", GridAutoFlowColumn, false},
		{"row dense", GridAutoFlowRow, true},
		{"dense column", GridAutoFlowColumn, true},
	} {
		flow, dense, ok := parseGridAutoFlow(tc.raw)
		if !ok || flow != tc.flow || dense != tc.dense {
			t.Errorf("%q = %d,%v,%v", tc.raw, flow, dense, ok)
		}
	}
	for _, raw := range []string{"", "row column", "dense dense", "sideways"} {
		if _, _, ok := parseGridAutoFlow(raw); ok {
			t.Errorf("%q must not parse", raw)
		}
	}
}

func TestGridPlacementParsesLinesSpansAndArea(t *testing.T) {
	placement, ok := parseGridPlacement("2 / span 2")
	if !ok || placement.Start.Kind != GridLineIndex || placement.Start.Index != 2 || placement.End.Kind != GridLineSpan || placement.End.Index != 2 {
		t.Errorf("placement = %+v, %v", placement, ok)
	}
	column, row, ok := parseGridArea("sidebar")
	if !ok || column.Start.Kind != GridLineName || column.Start.Name != "sidebar" || column.End.Name != "sidebar" || row.Start.Name != "sidebar" || row.End.Name != "sidebar" {
		t.Errorf("area = %+v, %+v, %v", column, row, ok)
	}
	column, row, ok = parseGridArea("Header")
	if !ok || column.Start.Name != "Header" || row.Start.Name != "Header" {
		t.Errorf("case-sensitive area name = %+v, %+v, %v", column, row, ok)
	}
	column, row, ok = parseGridArea("sidebar-start / main-start")
	if !ok || column.Start.Name != "main-start" || row.Start.Name != "sidebar-start" ||
		column.End.Kind != GridLineAuto || row.End.Kind != GridLineAuto {
		t.Errorf("named line placement = %+v, %+v, %v", column, row, ok)
	}
	column, row, ok = parseGridArea("sidebar-start / main-start / footer-end / content-end")
	if !ok || column.Start.Name != "main-start" || column.End.Name != "content-end" ||
		row.Start.Name != "sidebar-start" || row.End.Name != "footer-end" {
		t.Errorf("four named lines = %+v, %+v, %v", column, row, ok)
	}
	if _, _, ok := parseGridArea("a / b / c / d / e"); ok {
		t.Error("grid area with five lines must not parse")
	}
	column, row, ok = parseGridArea("1 / 2 / 3 / 4")
	if !ok || row.Start.Index != 1 || row.End.Index != 3 || column.Start.Index != 2 || column.End.Index != 4 {
		t.Errorf("numeric area = %+v, %+v, %v", column, row, ok)
	}
	column, row, ok = parseGridArea("1")
	if !ok || row.Start.Index != 1 || column.Start.Kind != GridLineAuto {
		t.Errorf("single line area = %+v, %+v, %v", column, row, ok)
	}
	span, ok := parseGridPlacement("span 2")
	if !ok || span.Start.Kind != GridLineSpan || span.Start.Index != 2 || span.End.Kind != GridLineAuto {
		t.Errorf("span placement = %+v, %v", span, ok)
	}
	for _, tc := range []struct {
		raw  string
		kind GridLineKind
		idx  int
		span int
	}{
		{"span", GridLineSpan, 1, 0},
		{"span 3", GridLineSpan, 3, 0},
		{"2 span", GridLineIndex, 2, 1},
		{"2 span 3", GridLineIndex, 2, 3},
		{"-2 span 3", GridLineIndex, -2, 3},
		{"-1", GridLineIndex, -1, 0},
	} {
		line, ok := parseGridLine(tc.raw)
		if !ok || line.Kind != tc.kind || line.Index != tc.idx || line.Span != tc.span {
			t.Errorf("line %q = %+v, %v", tc.raw, line, ok)
		}
	}
	for _, tc := range []struct {
		raw      string
		name     string
		backward bool
		span     int
	}{
		{"sidebar-start", "sidebar-start", false, 0},
		{"Sidebar-Start", "Sidebar-Start", false, 0},
		{"-sidebar-start", "sidebar-start", true, 0},
		{"sidebar-start span 2", "sidebar-start", false, 2},
		{"-sidebar-start span 2", "sidebar-start", true, 2},
	} {
		line, ok := parseGridLine(tc.raw)
		if !ok || line.Kind != GridLineName || line.Name != tc.name || line.Backward != tc.backward || line.Span != tc.span {
			t.Errorf("line %q = %+v, %v", tc.raw, line, ok)
		}
	}
	for _, raw := range []string{"span\t2", "span\n2"} {
		span, ok := parseGridPlacement(raw)
		if !ok || span.Start.Kind != GridLineSpan || span.Start.Index != 2 {
			t.Errorf("span %q = %+v, %v", raw, span, ok)
		}
	}
	if _, ok := parseGridLine("0"); ok {
		t.Error("grid line zero must not parse")
	}
	for _, raw := range []string{
		"2 /", "/ 2", "2 // 3", "span span", "1 span 2 3", "1 2", "auto auto",
		"span 0", "a b", "[a]", "0 span 2", "2 span 0",
	} {
		if _, ok := parseGridPlacement(raw); ok {
			t.Errorf("placement %q must not parse", raw)
		}
	}
}

func TestGridPlacementLonghandsAcceptLineNames(t *testing.T) {
	for _, raw := range []string{
		"grid-column: content-start / content-end;",
		"grid-column: -content-start;",
		"grid-row: header-start span 2;",
		"grid-column-start: sidebar;",
		"grid-row-end: -footer;",
		"grid-area: header-start / main-start;",
	} {
		st := gridRule(t, "card", raw)
		placed := st.GridColumn.Start.Kind != GridLineAuto || st.GridColumn.End.Kind != GridLineAuto ||
			st.GridRow.Start.Kind != GridLineAuto || st.GridRow.End.Kind != GridLineAuto
		if !placed {
			t.Errorf("%s did not reach the style: %+v", raw, st.GridColumn)
		}
	}
	st := gridRule(t, "card", "grid-column: -content-start;")
	if st.GridColumn.Start.Name != "content-start" || !st.GridColumn.Start.Backward {
		t.Errorf("backward longhand = %+v", st.GridColumn)
	}
	st = gridRule(t, "card", "grid-row: header-start span 2;")
	if st.GridRow.Start.Name != "header-start" || st.GridRow.Start.Span != 2 {
		t.Errorf("span longhand = %+v", st.GridRow)
	}
}

func TestGridPropertiesReachStyle(t *testing.T) {
	st := gridRule(t, "grid", `
		display: grid;
		grid-template-columns: 120px 1fr;
		grid-template-rows: 40px 2fr;
		grid-template-areas: "header main" "side main";
		grid-auto-flow: column dense;
		justify-content: center;
		align-content: space-between;
		justify-items: end;
		align-items: center;
		gap: 6px 10px;
	`)
	if st.Display != DisplayGrid || len(st.GridTemplateColumns) != 2 || len(st.GridTemplateRows) != 2 || len(st.GridTemplateAreas) != 2 {
		t.Errorf("grid style = %+v", st)
	}
	if st.GridAutoFlow != GridAutoFlowColumn || !st.GridAutoFlowDense || st.JustifyContent != JustifyCenter || st.AlignContent != ContentSpaceBetween || st.JustifyItems != AlignFlexEnd || st.AlignItems != AlignCenter {
		t.Errorf("grid alignment = %+v", st)
	}
	if st.RowGap.Resolve(Units{}) != 6 || st.ColumnGap.Resolve(Units{}) != 10 {
		t.Errorf("gap = %v/%v", st.RowGap, st.ColumnGap)
	}
}

func TestGridItemPropertiesReachStyle(t *testing.T) {
	st := gridRule(t, "button", `
		grid-column: 2 / span 2;
		grid-row: 1;
		justify-self: center;
		align-self: end;
	`)
	if st.GridColumn.Start.Index != 2 || st.GridColumn.End.Kind != GridLineSpan || st.GridColumn.End.Index != 2 || st.GridRow.Start.Index != 1 || st.GridRow.End.Kind != GridLineAuto || st.JustifySelf != AlignCenter || st.AlignSelf != AlignFlexEnd {
		t.Errorf("item grid style = %+v", st)
	}
	st = gridRule(t, "label", "grid-area: header;")
	if st.GridColumn.Start.Name != "header" || st.GridColumn.End.Name != "header" || st.GridRow.Start.Name != "header" || st.GridRow.End.Name != "header" {
		t.Errorf("grid area = %+v", st)
	}
}

func TestGridInitialValues(t *testing.T) {
	st := gridRule(t, "grid", `
		grid-template-columns: initial;
		grid-template-rows: initial;
		grid-template-areas: initial;
		grid-auto-flow: initial;
		justify-items: initial;
		justify-self: initial;
		grid-column: initial;
		grid-row: initial;
		grid-area: initial;
	`)
	if st.GridTemplateColumns != nil || st.GridTemplateRows != nil || st.GridTemplateAreas != nil || st.GridAutoFlow != GridAutoFlowRow || st.GridAutoFlowDense || st.JustifyItems != AlignStretch || st.JustifySelf != AlignAuto {
		t.Errorf("initial grid = %+v", st)
	}
	if st.GridColumn.Start.Kind != GridLineAuto || st.GridColumn.End.Kind != GridLineAuto || st.GridRow.Start.Kind != GridLineAuto || st.GridRow.End.Kind != GridLineAuto {
		t.Errorf("initial placement = %+v", st)
	}
}
