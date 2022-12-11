package megagrid

import (
	"math"
	"sort"
)

func astar(start *Cell, end *Cell, grid MegaGrid) bool {
	if start == nil || end == nil || grid == nil {
		return false
	}

	// reset all cells
	for _, c := range grid.ToSlice() {
		c.AG = 0           // distance from start node
		c.AH = 0           // (heuristic) distance from end node
		c.AF = c.AG + c.AH // G cost + H cost
		c.AClosed = false
		c.AOpen = false
	}

	// define lists
	open := make([]*Cell, 0, 1024)   // to be evaluated
	closed := make([]*Cell, 0, 1024) // already evaluated

	// add start cell
	open = append(open, start)
	start.AOpen = true

	// find a way
	for {
		// there is no way :(
		if len(open) == 0 {
			// error
			return false
		}

		// cell from open with the lowest f_cost
		sortList(open)
		current := open[0]

		// remove cell from open
		open = open[1:]
		current.AOpen = false

		// and add to close
		closed = append(closed, current)
		current.AClosed = true

		// is current the target cell?
		if current.xCol == end.xCol && current.yRow == end.yRow {
			// success
			return true
		}

		// foreach neighbour of the current sell
		for i, n := range current.neighbours {

			// ignore cells in close
			if n.AClosed || n.cType == Blocked || n.cType == None {
				continue // skip this neighbour
			}

			// calculating G, H and F
			nCost := current.neighboursCost[i] // costs from current to this neighbour
			dx := float64(n.xCol - end.xCol)
			dy := float64(n.yRow - end.yRow)
			d := math.Sqrt(dx*dx + dy*dy)

			g := nCost + current.AG // distance from start node
			h := int(d) * 10        // (heuristic) distance from end node
			f := g + h              // G cost + H cost

			// neighbour cost
			if g < n.AG || !n.AOpen {
				n.AParent = current
				n.AG = g
				n.AH = h
				n.AF = f
				if !n.AOpen {
					open = append(open, n)
					n.AOpen = true
				}
			}
		}
	}

}

func sortList(list []*Cell) {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].AF < list[j].AF
	})
}
