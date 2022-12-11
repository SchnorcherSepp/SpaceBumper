package ai

import (
	"SpaceBumper/myai/megagrid"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// memory
var lastPosCell *megagrid.Cell
var lastPosCellIter = 0

func MyAI(conn *net.TCPConn, world *megagrid.WorldMap) {

	// main loop
	for {
		world.Mux.Lock()
		//-------------------------------------------------------------------------------

		// my position
		nextStar := world.NextStar

		// goto
		goTo(conn, world, nextStar)

		//-------------------------------------------------------------------------------
		world.Mux.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

//------------- TARGET -----------------------------------------------------------------------------------------------//

func goTo(conn *net.TCPConn, world *megagrid.WorldMap, targetCell *megagrid.Cell) {

	// no more targets
	if targetCell == nil {
		setMove(megagrid.NewVector(0, 0), conn)
		return
	}

	// get current ship and position
	ship := world.PlayerShip
	myPosCell := world.MyPos(world.PlayerName)

	// calc path to target
	if path := world.MGrid.AStar(myPosCell, targetCell); path < 0 {
		setMove(megagrid.NewVector(0, 0), conn)
		return
	}

	// next step/cell in my path
	nextCell := myPosCell.AParent
	if nextCell == nil {
		setMove(megagrid.NewVector(0, 0), conn)
		return
	}

	// max. speed for straight lines
	if nextCell.AParent != nil {
		old := myPosCell
		mid := nextCell
		nxt := nextCell.AParent

		for {
			// calc vector
			v1 := megagrid.NewVector(float64(old.XCol()-mid.XCol()), float64(old.YRow()-mid.YRow()))
			v2 := megagrid.NewVector(float64(mid.XCol()-nxt.XCol()), float64(mid.YRow()-nxt.YRow()))
			// check straight line
			if v1.Angle() == v2.Angle() {
				nextCell = nextCell.AParent
			} else {
				break
			}
			// prepare next round
			if nextCell.AParent != nil {
				old = mid
				mid = nxt
				nxt = nextCell.AParent
			} else {
				break
			}
		}
	}

	// unfreeze
	if lastPosCell == nil || lastPosCell.XCol() != myPosCell.XCol() || lastPosCell.YRow() != myPosCell.YRow() {
		lastPosCell = myPosCell
		lastPosCellIter = 0
	} else {
		lastPosCellIter++
	}
	if lastPosCellIter > 40 {
		v := megagrid.NewVector((rand.Float64()-0.5)*1000, (rand.Float64()-0.5)*1000)
		println("UNFREEZ", v.X, v.Y)
		setMove(v, conn)
		time.Sleep(300 * time.Millisecond)
		return
	}

	// calc move
	x := float64(nextCell.XCol()-myPosCell.XCol()) - ship.Velocity.X*0.55 // important factor!!
	y := float64(nextCell.YRow()-myPosCell.YRow()) - ship.Velocity.Y*0.55
	v := megagrid.NewVector(x*1000, y*1000)
	println("MOVE", v.X, v.Y)
	setMove(v, conn)
}

//----------------- HELPER -------------------------------------------------------------------------------------------//

func setMove(acceleration *megagrid.Vector, conn *net.TCPConn) {
	mov := fmt.Sprintf("%f|%f\n", acceleration.X, acceleration.Y)
	_, _ = conn.Write([]byte(mov))
}
