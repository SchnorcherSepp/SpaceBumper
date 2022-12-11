package megagrid

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"
)

type Grid [][]*Cell

func NewGrid(txt []byte) (Grid, error) {

	// remove magic bytes
	txt = bytes.ReplaceAll(txt, []byte{0xef, 0xbb, 0xbf}, []byte{})

	// split lines
	s := strings.ReplaceAll(string(txt), "\r", "") // remove '\r'
	s = strings.ReplaceAll(s, "|", "")             // remove '|'
	lines := strings.Split(s, "\n")                // split lines ('\n')

	// parse data
	var xWidth int
	var yHeight = len(lines)
	var grid Grid

	for yRow, l := range lines {
		// action for first line
		if yRow == 0 {
			xWidth = len(l)           // get width
			grid = make(Grid, xWidth) // init Grid
		}
		// CHECK: each line must have the same width
		if xWidth != len(l) {
			return nil, errors.New(fmt.Sprintf("invalid width: %d is not %d in line %d", len(l), xWidth, yRow+1))
		}
		// more grit init
		for xCol := 0; xCol < xWidth; xCol++ {
			if grid[xCol] == nil {
				grid[xCol] = make([]*Cell, yHeight)
			}
			// add cells to Gr
			grid[xCol][yRow] = NewCell(l[xCol], xCol, yRow)
		}
	}

	// return
	//grid.Print()
	return grid, nil
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

func (g Grid) Dimensions() (xWidth, yHeight int) {
	return dimensions(g)
}

func (g Grid) ToSlice() []*Cell {
	return toSlice(g)
}

func (g Grid) CellByCoordinates(x, y float64) *Cell {
	xCol := int(math.Floor(x / CellSize))
	yRow := int(math.Floor(y / CellSize))

	if g.OutOfBound(xCol, yRow) {
		return NewCell(None, xCol, yRow)
	} else {
		return g[xCol][yRow]
	}
}

func (g Grid) OutOfBound(xCol, yRow int) bool {
	return outOfBound(g, xCol, yRow)
}

func (g Grid) Stars() []*Cell {
	list := make([]*Cell, 0, 256)
	for _, c := range g.ToSlice() {
		if c.Type() == Star {
			list = append(list, c)
		}
	}
	return list
}

func (g Grid) Print() {
	printGrid(g)
}

//--------  Helper  --------------------------------------------------------------------------------------------------//

func dimensions(g [][]*Cell) (xWidth, yHeight int) {
	xWidth = len(g)
	if xWidth < 1 {
		yHeight = 0
	} else {
		yHeight = len(g[0])
	}
	return
}

func toSlice(g [][]*Cell) []*Cell {
	xWidth, yHeight := dimensions(g)
	list := make([]*Cell, xWidth*yHeight)

	i := 0
	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			list[i] = g[xCol][yRow]
			i++
		}
	}
	return list
}

func outOfBound(g [][]*Cell, xCol, yRow int) bool {
	xWidth, yHeight := dimensions(g)
	return xWidth <= xCol || xCol < 0 || yHeight <= yRow || yRow < 0
}

func printGrid(g [][]*Cell) {
	xWidth, yHeight := dimensions(g)

	// top border
	fmt.Print("+")
	for xCol := 0; xCol < xWidth; xCol++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	// grid
	for yRow := 0; yRow < yHeight; yRow++ {
		fmt.Print("|") // left border
		for xCol := 0; xCol < xWidth; xCol++ {
			fmt.Print(string(g[xCol][yRow].Type())) // cell type
		}
		fmt.Print("|\n") // right border
	}

	// bottom border
	fmt.Print("+")
	for xCol := 0; xCol < xWidth; xCol++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")
}
