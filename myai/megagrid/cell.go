package megagrid

import (
	"bytes"
	"fmt"
)

//--------  Constants  -----------------------------------------------------------------------------------------------//

// all supported cell types
const (
	None    = ' ' // None cell: ships on this cell fell into the void
	Blocked = '#' // Blocked cell: ships crash at this cell and rebound
	Boost   = 'b' // Boost cell: boosts velocity on contact
	Slow    = 's' // Slow cell: slow velocity on contact
	Tile    = '.' // Tile cell: is normal flat ground that can be driven on
	Star    = 'x' // Star cell: increases the score on contact and then disappears (Tile)
	Anti    = 'a' // Anti cell: reduce the score on contact and then disappears (Tile)
	Spawn   = 'o' // Spawn cell: point where players are randomly placed
)

// CellTypes all supported cell types as slice
var CellTypes = []byte{None, Blocked, Boost, Slow, Tile, Star, Anti, Spawn}

// CellSize is the dimension of the square cell
const CellSize = 40.0 // cell image 40x40

//--------  Struct  --------------------------------------------------------------------------------------------------//

// Cell always has a defined type (see CellTypes).
// It's positioned on a chessboard (Grid or MegaGrid) with rows and columns (int).
type Cell struct {
	cType byte // cell types
	xCol  int  // column of the grid
	yRow  int  // row of the grid
	cost  int  // set by UpdateCost()

	// set by UpdateNeighbours()
	neighbours     []*Cell
	neighboursCost []int

	// set by AStar()
	AG      int   // distance from start node
	AH      int   // (heuristic) distance from end node
	AF      int   // G cost + H cost
	AOpen   bool  // is in open list
	AClosed bool  // is in close list
	AParent *Cell // parent of this cell
}

// NewCell create a new cell.
// If the type is invalid, then the default (None) is used.
func NewCell(cType byte, xCol, yRow int) *Cell {
	c := &Cell{
		cType: cType,
		xCol:  xCol,
		yRow:  yRow,
	}
	c.SetType(cType)
	return c
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

// Type return the CellTypes.
func (c *Cell) Type() byte {
	return c.cType
}

// XCol return the column on the Grid.
func (c *Cell) XCol() int {
	return c.xCol
}

// YRow return the row on the Grid.
func (c *Cell) YRow() int {
	return c.yRow
}

// Cost return the move cost.
func (c *Cell) Cost() int {
	return c.cost
}

// Neighbours return all valid neighbour cells of this cell
func (c *Cell) Neighbours() []*Cell {
	return c.neighbours
}

// NeighboursCost is identical with Neighbours(), but return the cost to move to this neighbour
func (c *Cell) NeighboursCost() []int {
	return c.neighboursCost
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

// SetType change the CellTypes.
// If the type is invalid, then the default (None) is used.
func (c *Cell) SetType(cType byte) {
	// check cell type
	if bytes.IndexByte(CellTypes, cType) < 0 {
		fmt.Printf("err: invalid cell type '%s' at X %d and Y %d\n", string(cType), c.xCol, c.yRow)
		cType = None
	}
	// set type
	c.cType = cType
}
