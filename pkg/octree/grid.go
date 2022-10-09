package octree

const CellSizeFactor = 5.0

type Grid struct {
	Width          int
	Height         int
	AABB           *AABB
	Depth          int
	SquaredSpacing float64
	Items          map[int]*GridCell
}

type GridCell struct {
	Neighbours []*GridCell
	points     []*Vector3
	sparseGrid *Grid
}

type GridIndex struct {
	I, J, K int
}

func NewGridIndex(i, j, k int) *GridIndex {
	return &GridIndex{I: i, J: j, K: k}
}

func (gi *GridIndex) Compare(b *GridIndex) bool {
	if gi.I < b.I {
		return true
	} else if gi.I == b.I && gi.J < b.J {
		return true
	} else if gi.I == b.I && gi.J == b.J && gi.K < b.K {
		return true
	}

	return false
}
func (c *GridCell) IsDistant(p *Point, squaredSpacing float64) bool {
	for _, point := range c.points {
		if p.Pos.SquaredDistanceTo(point) < squaredSpacing {
			return false
		}
	}

	return true
}

func NewGridCell(grid *Grid, index *GridIndex) *GridCell {
	cell := &GridCell{
		sparseGrid: grid,
	}

	for i := Max(index.I-1, 0); i < Min(grid.Width-1, index.I+1); i++ {
		for j := Max(index.J-1, 0); j < Min(grid.Height-1, index.J+1); j++ {
			for k := Max(index.K-1, 0); k < Min(grid.Depth-1, index.K+1); k++ {
				key := (k << 40) | (j << 20) | i
				if neighbour, ok := grid.Items[key]; ok {
					if &neighbour != &cell {
						cell.Neighbours = append(cell.Neighbours, neighbour)
						neighbour.Neighbours = append(cell.Neighbours, cell)
					}
				}
			}
		}
	}

	return cell
}

func (c *GridCell) Add(p *Point) {
	c.points = append(c.points, p.Pos)
}

func NewGrid(aabb *AABB, spacing float64) *Grid {
	return &Grid{
		AABB:           aabb,
		Width:          int(aabb.Size.X / (spacing * CellSizeFactor)),
		Height:         int(aabb.Size.Y / (spacing * CellSizeFactor)),
		Depth:          int(aabb.Size.Z / (spacing * CellSizeFactor)),
		SquaredSpacing: spacing * spacing,
		Items:          make(map[int]*GridCell),
	}
}

func (g *Grid) IsDistant(point *Point, cell *GridCell) bool {
	if !cell.IsDistant(point, g.SquaredSpacing) {
		return false
	}

	for _, neighbour := range cell.Neighbours {
		if !neighbour.IsDistant(point, g.SquaredSpacing) {
			return false
		}
	}

	return true
}

func (g *Grid) Add(p *Point) bool {
	nx := (int)(float64(g.Width) * (p.Pos.X - g.AABB.Min.X) / g.AABB.Size.X)
	ny := (int)(float64(g.Height) * (p.Pos.Y - g.AABB.Min.Y) / g.AABB.Size.Y)
	nz := (int)(float64(g.Depth) * (p.Pos.Z - g.AABB.Min.Z) / g.AABB.Size.Z)

	i := Min(nx, g.Width-1)
	j := Min(ny, g.Height-1)
	k := Min(nz, g.Depth-1)

	index := NewGridIndex(i, j, k)

	key := (k << 40) | (j << 20) | i

	if _, ok := g.Items[key]; !ok {
		g.Items[key] = NewGridCell(g, index)
	}

	if g.IsDistant(p, g.Items[key]) {
		g.Items[key].Add(p)
		return true
	} else {
		return false
	}
}
