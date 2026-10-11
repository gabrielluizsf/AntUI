package template

func (b *gridBatch) gridAreas() map[string]gridArea {
	areas := make(map[string]gridArea)
	for row, cells := range b.st.GridTemplateAreas {
		for col, name := range cells {
			if name == "." || name == "..." {
				continue
			}
			area, exists := areas[name]
			if !exists {
				areas[name] = gridArea{col: col, row: row, colSpan: 1, rowSpan: 1}
				continue
			}
			area.colSpan = max(area.col, col) - min(area.col, col) + 1
			area.rowSpan = max(area.row, row) - min(area.row, row) + 1
			if area.col != col || area.row != row {
				area.hasOther = true
				area.otherCol, area.otherRow = col, row
			}
			areas[name] = area
		}
	}
	return areas
}
