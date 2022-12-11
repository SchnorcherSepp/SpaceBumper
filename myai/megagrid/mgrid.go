package megagrid

import "math"

const ScalingFactor = 5.0

type MegaGrid [][]*Cell

func NewMegaGrid(g Grid) MegaGrid {
	xWidth, yHeight := g.Dimensions()

	// make mega grid
	mg := make(MegaGrid, xWidth*ScalingFactor)
	for xCol := 0; xCol < xWidth*ScalingFactor; xCol++ {
		mg[xCol] = make([]*Cell, yHeight*ScalingFactor)
	}

	// set data
	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			for ix := 0; ix < ScalingFactor; ix++ {
				for iy := 0; iy < ScalingFactor; iy++ {
					newX := ScalingFactor*xCol + ix
					newY := ScalingFactor*yRow + iy
					mg[newX][newY] = NewCell(g[xCol][yRow].Type(), newX, newY)
				}
			}
		}
	}

	// update costs & neighbours
	mg.UpdateCosts()
	mg.UpdateNeighbours()

	// return
	//mg.Print()
	return mg
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

func (mg MegaGrid) Dimensions() (xWidth, yHeight int) {
	return dimensions(mg)
}

func (mg MegaGrid) ToSlice() []*Cell {
	return toSlice(mg)
}

func (mg MegaGrid) OutOfBound(xCol, yRow int) bool {
	return outOfBound(mg, xCol, yRow)
}

func (mg MegaGrid) CellByCoordinates(x, y float64) *Cell {
	xCol := int(math.Floor(x / (CellSize / ScalingFactor)))
	yRow := int(math.Floor(y / (CellSize / ScalingFactor)))

	if mg.OutOfBound(xCol, yRow) {
		return NewCell(None, xCol, yRow)
	} else {
		return mg[xCol][yRow]
	}
}

func (mg MegaGrid) Stars() []*Cell {
	list := make([]*Cell, 0, 256)
	for _, c := range mg.ToSlice() {
		if c.Type() == Star && (c.xCol-2)%ScalingFactor == 0 && (c.yRow-2)%ScalingFactor == 0 {
			list = append(list, c)
		}
	}
	return list
}

func (mg MegaGrid) Print() {
	printGrid(mg)
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

func (mg MegaGrid) UpdateCosts() {
	xWidth, yHeight := mg.Dimensions()

	// extend block
	blockCells := make([]*Cell, 0, xWidth*yHeight)
	avoidCells := make([]*Cell, 0, xWidth*yHeight)

	// set data
	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			cell := mg[xCol][yRow]
			cost := 10 // basis cost is 10

			switch cell.Type() {
			case None: // None cell: ships on this cell fell into the void
				cost *= -1 // BLOCK
				blockCells = append(blockCells, cell)
				break
			case Blocked: // Blocked cell: ships crash at this cell and rebound
				cost *= -1 // BLOCK
				blockCells = append(blockCells, cell)
				break
			case Anti: // Anti cell: reduce the score on contact and then disappears (Tile)
				cost *= 5 // bad
				avoidCells = append(avoidCells, cell)
				break
			case Slow: // Slow cell: slow velocity on contact
				cost *= 5 // bad
				avoidCells = append(avoidCells, cell)
				break
			case Boost: // Boost cell: boosts velocity on contact
				cost *= 7 // unpredictable
				avoidCells = append(avoidCells, cell)
				break
			case Tile: // Tile cell: is normal flat ground that can be driven on
				cost *= 1 // normal
				break
			case Star: // Star cell: increases the score on contact and then disappears (Tile)
				cost *= 1 // normal
				break
			case Spawn: // Spawn cell: point where players are randomly placed
				cost *= 1 // normal
				break
			}

			cell.cost = cost
		}
	}

	// process extend block
	for _, c := range blockCells {
		for ix := -2; ix < 3; ix++ {
			for iy := -2; iy < 3; iy++ {
				x := c.xCol + ix
				y := c.yRow + iy
				if !mg.OutOfBound(x, y) && mg[x][y].cost >= 0 {
					mg[x][y].cost = 9999
				}
			}
		}
	}

	// process extend speed, slow, anti
	for _, c := range avoidCells {
		for ix := -2; ix < 3; ix++ {
			for iy := -2; iy < 3; iy++ {
				x := c.xCol + ix
				y := c.yRow + iy
				if !mg.OutOfBound(x, y) && mg[x][y].cost >= 0 && mg[x][y].cost < 30 {
					mg[x][y].cost = 30
				}
			}
		}
	}
}

func (mg MegaGrid) UpdateNeighbours() {
	for _, c := range mg.ToSlice() {
		cellNeighbours(c, mg)
	}
}

func (mg MegaGrid) AStar(start *Cell, end *Cell) int {
	if path := astar(end, start, mg); !path {
		return -1
	} else {
		return start.AG
	}
}

//--------  Helper  --------------------------------------------------------------------------------------------------//

func cellNeighbours(currentCell *Cell, mg MegaGrid) {

	neighbours := make([]*Cell, 0, 8)
	neighboursCost := make([]int, 0, 8)

	// cell is invalid?
	if currentCell.Cost() < 0 {
		return // EXIT
	}

	// helper
	type Co struct {
		X int
		Y int

		Cell   *Cell
		factor float64

		// diagonal also need left and right free
		NecessaryLeft  *Co
		NecessaryRight *Co
	}

	// straight
	up := &Co{X: currentCell.xCol, Y: currentCell.yRow - 1, factor: 1}
	down := &Co{X: currentCell.xCol, Y: currentCell.yRow + 1, factor: 1}
	left := &Co{X: currentCell.xCol - 1, Y: currentCell.yRow, factor: 1}
	right := &Co{X: currentCell.xCol + 1, Y: currentCell.yRow, factor: 1}

	// diagonal
	// need left and right free
	leftUp := &Co{X: currentCell.xCol - 1, Y: currentCell.yRow - 1, NecessaryLeft: left, NecessaryRight: up, factor: 1.4}
	rightUp := &Co{X: currentCell.xCol + 1, Y: currentCell.yRow - 1, NecessaryLeft: up, NecessaryRight: right, factor: 1.4}
	leftDown := &Co{X: currentCell.xCol - 1, Y: currentCell.yRow + 1, NecessaryLeft: down, NecessaryRight: left, factor: 1.4}
	rightDown := &Co{X: currentCell.xCol + 1, Y: currentCell.yRow + 1, NecessaryLeft: right, NecessaryRight: down, factor: 1.4}

	// allNeighbour
	for _, n := range []*Co{up, down, left, right, leftUp, rightUp, leftDown, rightDown} {

		// out of bound
		if mg.OutOfBound(n.X, n.Y) {
			continue
		}
		n.Cell = mg[n.X][n.Y]

		// blocked
		if n.Cell.Cost() < 0 {
			n.Cell = nil // remove invalid cell
			continue
		}

		// diagonal: left and right blocked?
		if n.NecessaryLeft != nil && (n.NecessaryLeft.Cell == nil || n.NecessaryRight.Cell == nil || n.NecessaryLeft.Cell.cost > 9000 || n.NecessaryRight.Cell.cost > 9000) {
			n.Cell = nil // remove invalid cell
			continue
		}

		// cell is ok
		neighbours = append(neighbours, n.Cell)
		neighboursCost = append(neighboursCost, int(n.factor*float64(currentCell.cost+n.Cell.cost)))
	}

	// update cell
	currentCell.neighbours = neighbours
	currentCell.neighboursCost = neighboursCost
}
