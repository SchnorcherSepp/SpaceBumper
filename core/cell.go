package core

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

// CellRadius is the cell radius
const CellRadius = CellSize / 2

//--------  Struct  --------------------------------------------------------------------------------------------------//

// Cell always has a defined type (see CellTypes).
// It's positioned on a chessboard (grid) with rows and columns (int).
// At the same time, the cell has dimensions (CellSize) in a coordinate system (float).
type Cell struct {
	cType       byte    // cell types
	xCol        int     // column of the grid
	yRow        int     // row of the grid
	topLeft     *Vector // cell border (coordinate top left)
	bottomRight *Vector // cell border (coordinate bottom right)
	center      *Vector // cell position (coordinates of the center)
}

// NewCell create a new cell.
// If the type is invalid, then the default (None) is used.
func NewCell(cType byte, xCol, yRow int) *Cell {
	// check cell type
	if bytes.IndexByte(CellTypes, cType) < 0 {
		fmt.Printf("err: invalid cell type '%s' at xCol %d and yRow %d\n", string(cType), xCol, yRow)
		cType = None
	}

	// calc cell border
	topLeft := NewVector(float64(xCol)*CellSize, float64(yRow)*CellSize)
	bottomRight := NewVector(float64(xCol)*CellSize+CellSize, float64(yRow)*CellSize+CellSize)

	// calc position (cell center)
	center := NewVector(topLeft.X()+CellRadius, topLeft.Y()+CellRadius)

	// return
	return &Cell{
		cType:       cType,
		xCol:        xCol,
		yRow:        yRow,
		topLeft:     topLeft,
		bottomRight: bottomRight,
		center:      center,
	}
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

// Type return the CellTypes.
// The value is immutable.
func (c *Cell) Type() byte {
	return c.cType
}

// XCol return the column on the grid.
// The value is immutable.
func (c *Cell) XCol() int {
	return c.xCol
}

// YRow return the row on the grid.
// The value is immutable.
func (c *Cell) YRow() int {
	return c.yRow
}

// TopLeft return the upper left coordinate of this cell.
// The value is immutable (vector clone).
func (c *Cell) TopLeft() *Vector {
	return c.topLeft.Clone()
}

// BottomRight return the lower right coordinate of this cell.
// The value is immutable (vector clone).
func (c *Cell) BottomRight() *Vector {
	return c.bottomRight.Clone()
}

// Center return the coordinate of this cell (center).
// The value is immutable (vector clone).
func (c *Cell) Center() *Vector {
	return c.center.Clone()
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

// SetType change the CellTypes.
func (c *Cell) SetType(t byte) {
	c.cType = t
}
