package ai

import (
	"SpaceBumper/myai/megagrid"
	"SpaceBumper/myai/mgui"
	"net"
	"time"
)

const NAME = "Glawischnig"
const COLOR = "green"

func RunAI(srvAddr, srvPort string) {

	// connect to server
	tcpAddr, err := net.ResolveTCPAddr("tcp", srvAddr+":"+srvPort)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	// start game
	_, _ = conn.Write([]byte("pass|" + NAME + "|" + COLOR + "\n"))

	//---------------------------------------------------------------

	// create dummy world
	world, err := megagrid.NewWorldMap(make([]byte, 0), NAME)
	if err != nil {
		panic(err)
	}

	// start world updater
	go world.UpdateStreamHandler(conn)

	// wait for world initialisation
	for {
		world.Mux.Lock()
		if world.Iteration > 0 && len(world.Grid) > 0 && world.PlayerShip != nil {
			world.Mux.Unlock()
			break
		}
		world.Mux.Unlock()
		time.Sleep(500 * time.Microsecond)
	}

	//---------------------------------------------------------------

	// start AI
	println("START AI")
	go MyAI(conn, world)

	//---------------------------------------------------------------

	// run GridView (BLOCKING)
	err = mgui.RunGridView("Grid View", world)
	if err != nil {
		panic(err)
	}
}
