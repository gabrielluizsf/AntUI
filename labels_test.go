package antui

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func TestMeasuringALabel(t *testing.T) {
	if got, want := LabelWidth("OK", 16), canvas.TextWidth("OK")+32; got != want {
		t.Errorf("a label is %d wide, want %d", got, want)
	}
	if got, want := LabelHeight(6), canvas.FontHeight+12; got != want {
		t.Errorf("a label is %d tall, want %d", got, want)
	}
	if got := LabelWidth("", 16); got != 32 {
		t.Errorf("an empty label is %d wide, want its padding", got)
	}
}

func TestMeasuringARowOfLabels(t *testing.T) {
	labels := []string{"Cancel", "Empty", "A scene", "Something moving"}

	want := 0
	for i, label := range labels {
		if i > 0 {
			want += 8
		}
		want += canvas.TextWidth(label) + 32
	}
	if got := LabelsWidth(labels, 16, 8); got != want {
		t.Errorf("the row is %d wide, want %d", got, want)
	}
	if got := LabelsWidth(nil, 16, 8); got != 0 {
		t.Errorf("no labels are %d wide", got)
	}
	// One label is one label: no gap is counted before the first.
	if got, want := LabelsWidth(labels[:1], 16, 8), LabelWidth(labels[0], 16); got != want {
		t.Errorf("one label in a row is %d wide, want %d", got, want)
	}
}

func TestWrappingARowOfLabels(t *testing.T) {
	labels := []string{"Cancel", "Empty", "A scene", "Something moving"}
	whole := LabelsWidth(labels, 16, 8)

	// Room for the lot is one row.
	if rows := WrapLabels(labels, 16, 8, whole); len(rows) != 1 || rows[0] != 4 {
		t.Errorf("%v rows with room for them all, want one of four", rows)
	}
	// A pixel short of the lot is two rows, and every label is on one of them.
	rows := WrapLabels(labels, 16, 8, whole-1)
	if len(rows) < 2 {
		t.Fatalf("%v rows with no room for them all", rows)
	}
	counted := 0
	for _, count := range rows {
		counted += count
	}
	if counted != len(labels) {
		t.Errorf("%d labels across the rows, want %d", counted, len(labels))
	}
	// No row is wider than it was allowed to be, unless one label is.
	at := 0
	for _, count := range rows {
		row := labels[at : at+count]
		at += count
		if width := LabelsWidth(row, 16, 8); width > whole-1 && count > 1 {
			t.Errorf("a row of %d is %d wide, past %d", count, width, whole-1)
		}
	}

	// A label wider than anything gets a row to itself rather than nothing.
	rows = WrapLabels([]string{"a", "an unreasonably long answer", "b"}, 16, 8, 40)
	if len(rows) != 3 {
		t.Errorf("%v rows, want one each", rows)
	}
	if rows := WrapLabels(nil, 16, 8, 100); rows != nil {
		t.Errorf("%v rows for no labels", rows)
	}
}