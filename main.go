package main

import (
	"SpaceBumper/core"
	"SpaceBumper/gui"
	"SpaceBumper/remote"
	"flag"
	"os"
	"strconv"
	"time"
)

const VERSION = "1.0b"

func main() {

	// world settings
	mapName := flag.String("map", "map1", "the name of the player map")
	endtime := flag.Uint64("endtime", 10800, "maximum ticks until the game ends")

	// remote server settings
	remotePly := flag.Bool("remote", false, "starts the server for remote play")
	srvAddr := flag.String("addr", "localhost", "server ip; needs remote=true")
	srvPort := flag.String("port", "3333", "server port; needs remote=true")
	player := flag.String("player", "2", "how many players to wait for; needs remote=true")

	// local player settings
	localPly := flag.Bool("local", true, "enable local game with mouse; needs headless=false")
	localName := flag.String("name", "Local Player", "your local player name; needs local=true")
	localColor := flag.String("color", "blue", "your local player color; needs local=true")

	// gui settings
	headless := flag.Bool("headless", false, "enable or disable GUI")

	// parse flags
	flag.Parse()

	// print defaults
	if len(os.Args) <= 1 {
		println("SpaceBumper", VERSION)
		println("---------------")
		flag.PrintDefaults()
		println()
		os.Exit(0)
	}

	//------------------------------------------------------------------------------

	// create world
	world, err := core.LoadWorldMap(*mapName, *endtime)
	if err != nil {
		panic(err)
	}

	// start server
	if *remotePly {
		waitPlayer, err := strconv.Atoi(*player)
		if err != nil {
			panic(err)
		}
		go remote.RunServer(*srvAddr, *srvPort, world, waitPlayer)
	}

	// add local player
	if *localPly {
		_, err = world.AddPlayer(*localName, *localColor, nil)
		if err != nil {
			panic(err)
		}
	}

	// run GUI (blocking)
	if *headless {
		for {
			world.Update()
			time.Sleep(16 * time.Millisecond) // ~ 60 tick/sec
		}
	} else {
		if err := gui.RunGame("Space Bumper", world, true); err != nil {
			panic(err)
		}
	}
}
