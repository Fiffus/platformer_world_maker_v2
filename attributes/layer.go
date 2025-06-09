package attributes

type Layer [][]Tile

func (l *Layer) Construct(rows, cols int, baseTileSize float64) {
	*l = make(Layer, rows)
	for row := range *l {
		(*l)[row] = make([]Tile, cols)
		for col := range (*l)[row] {
			(*l)[row][col].Construct(
				Spatial{X: float64(col) * baseTileSize, Y: float64(row) * baseTileSize},
				baseTileSize,
				"air",
				nil,
			)
		}
	}
}
