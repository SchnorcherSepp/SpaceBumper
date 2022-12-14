package megagrid

import (
	"sync"
	"time"
)

type WorldMap struct {
	Mux *sync.Mutex

	Iteration     uint64
	EndTime       uint64
	MaxUpdateTime time.Duration

	Grid  Grid
	MGrid MegaGrid

	PlayerName string
	PlayerShip *Ship

	Players  []*Ship // all players (alive and dead)
	NextStar *Cell   // on MegaGrid
}

func NewWorldMap(txt []byte, playerName string) (*WorldMap, error) {

	// build grid
	g, err := NewGrid(txt)
	if err != nil {
		return nil, err
	}

	// build mega grid
	mg := NewMegaGrid(g)

	// return world
	world := &WorldMap{
		Mux:        new(sync.Mutex),
		Grid:       g,
		MGrid:      mg,
		PlayerName: playerName,
	}
	return world, nil
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

func (w *WorldMap) MyPos(playerName string) *Cell {
	for _, p := range w.Players {
		if p.Name == playerName {
			w.PlayerShip = p
			return w.MGrid.CellByCoordinates(p.Position.X, p.Position.Y)
		}
	}
	return nil
}

func (w *WorldMap) UpdateStars(playerName string) {
	var bestD = 999999999999999
	var bestS *Cell

	myPos := w.MyPos(playerName)

	// find best star
	for _, star := range w.MGrid.Stars() {
		d := w.MGrid.AStar(myPos, star)
		sts := w.starCluster(star)
		if sts > 50 {
			println("cluster", sts, star.xCol, star.yRow)
			d -= sts * 50
		}
		if d < bestD {
			bestD = d
			bestS = star
		}
	}

	// set best star and calc path
	w.MGrid.AStar(myPos, bestS)
	w.NextStar = bestS
}

func (w *WorldMap) starCluster(star *Cell) int {
	x := star.XCol()
	y := star.YRow()

	sts := 0

	for tmpX := x - 12; tmpX < x+12; tmpX++ {
		for tmpY := y - 12; tmpY < y+12; tmpY++ {
			if !w.MGrid.OutOfBound(tmpX, tmpY) {
				if w.MGrid[tmpX][tmpY].cType == Star {
					sts++
				}
			}
		}
	}

	return sts
}
