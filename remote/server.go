package remote

import (
	"SpaceBumper/core"
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
	"sync"
)

type server struct {
	host       string
	port       string
	world      *core.WorldMap
	waitPlayer int

	mux *sync.Mutex
}

// RunServer starts a server and makes the game world available remotely.
func RunServer(host, port string, world *core.WorldMap, waitPlayer int) {

	// Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("RunServer: %v\n", err)
	}

	// Close the listener when the application closes.
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	// Freeze world
	world.Freeze(true) // undo in registerPlayer()

	// server
	ser := &server{
		host:       host,
		port:       port,
		world:      world,
		waitPlayer: waitPlayer,
		mux:        new(sync.Mutex),
	}

	fmt.Println("START SERVER [" + host + ":" + port + "]")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, ser)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, ser *server) {
	ser.mux.Lock()
	defer ser.mux.Unlock()

	// prepare line reader
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// read first line (ended with \n or \r\n)
	line, _ := tp.ReadLine()

	// EXIT
	if line == "EXIT" {
		fmt.Printf("EXIT by %v\n", conn.RemoteAddr())
		os.Exit(0)
	}

	// vars
	var retMsg = ""

	// extract command
	// format:  "{pass}|{name}|{color}\n"
	param := strings.Split(line, "|")
	if len(param) != 3 {
		retMsg = "ERROR: invalid command! use '{pass}|{name}|{color}\\n'"

	} else {
		// extract name and color
		name := param[1]
		color := param[2]
		fmt.Printf("request: name=%s, color=%s\n", name, color)

		// add player
		id, err := ser.world.AddPlayer(name, color, conn)
		if err != nil {
			retMsg = err.Error()
		} else {
			retMsg = fmt.Sprintf("PLAYERID:%d", id)
		}
	}

	// write
	_, err := conn.Write([]byte(retMsg + "\n"))
	if err != nil {
		fmt.Printf("comWrite: %v\n", err)
	}

	// un-freeze
	if ser.waitPlayer <= len(ser.world.Players()) {
		ser.world.Freeze(false)
	}
}
