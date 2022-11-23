package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WorldMap represents the current world status
// with the grid and all players.
type WorldMap struct {
	freeze        bool
	iteration     uint64
	endtime       uint64
	maxUpdateTime time.Duration

	xWidth  int       // grid size (width)
	yHeight int       // grid size (height)
	grid    [][]*Cell // grid (map)
	spawns  []*Cell   // spawn cell list

	players []*Ship // all players (alive and dead)
}

// LoadWorldMap loads the map file '{mapName}.txt' from the local directory or the 'maps' subdirectory.
// The Order is './*' then './maps/*' then '../maps/*' and then '../../maps/*'.
//
// The map is defined with characters in a text file. Each character is a cell (see CellTypes).
// The chars are the columns, the lines are the rows of the grid.
// All lines must contain the same number of characters.
// Each line must end with the character '|'.
func LoadWorldMap(mapName string, endtime uint64) (*WorldMap, error) {

	// search map file
	if !strings.HasSuffix(strings.ToLower(mapName), ".txt") {
		mapName += ".txt"
	}
	paths := []string{mapName, "maps/" + mapName, "../maps/" + mapName, "../../maps/" + mapName}
	for _, p := range paths {
		_, e := os.Stat(p)
		if e == nil {
			mapName = p
			break
		}
	}

	// read file
	b, err := os.ReadFile(mapName)
	if err != nil {
		return nil, err
	}

	// remove utf8 magic bytes
	b = bytes.ReplaceAll(b, []byte{0xef, 0xbb, 0xbf}, []byte{})

	// split lines
	s := strings.ReplaceAll(string(b), "\r", "") // remove '\r'
	s = strings.ReplaceAll(s, "|", "")           // remove '|'
	lines := strings.Split(s, "\n")              // split lines ('\n')

	// parse data
	var xWidth int
	var yHeight = len(lines)
	var grid [][]*Cell
	var spawns = make([]*Cell, 0)

	for yRow, l := range lines {
		// action for first line
		if yRow == 0 {
			xWidth = len(l)                // init width
			grid = make([][]*Cell, xWidth) // init grid
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
			// add cells to grid
			c := NewCell(l[xCol], xCol, yRow)
			grid[xCol][yRow] = c
			// save spawn positions
			if c.Type() == Spawn {
				spawns = append(spawns, c)
			}
		}
	}

	// end time
	if endtime <= 0 {
		endtime = math.MaxUint64
	}

	// build map
	wm := &WorldMap{
		freeze:        false,
		iteration:     0,
		endtime:       endtime,
		maxUpdateTime: 0,
		xWidth:        xWidth,
		yHeight:       yHeight,
		grid:          grid,
		spawns:        spawns,
		players:       make([]*Ship, 0, len(spawns)),
	}

	// return
	wm.Print()
	return wm, nil
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

// Stats returns
//
//	iteration the current iteration
//  endtime is the max. iteration
//	maxUpdateTime the longest running time of the Update() function.
func (m *WorldMap) Stats() (iteration, endtime uint64, maxUpdateTime time.Duration) {
	return m.iteration, m.endtime, m.maxUpdateTime
}

// XWidth returns the grid width
func (m *WorldMap) XWidth() int {
	return m.xWidth
}

// YHeight returns the grid height
func (m *WorldMap) YHeight() int {
	return m.yHeight
}

// Grid returns the grid (map)
func (m *WorldMap) Grid() [][]*Cell {
	return m.grid
}

// Spawns returns all spawn cells
func (m *WorldMap) Spawns() []*Cell {
	return m.spawns
}

// Players returns all players
func (m *WorldMap) Players() []*Ship {
	return m.players
}

// MaxPlayers returns the maximum supported players of this map (is the spawner count).
func (m *WorldMap) MaxPlayers() int {
	return len(m.spawns)
}

// Player returns the requested player.
// Throws an error if the player was not found.
func (m *WorldMap) Player(playerID int) (*Ship, error) {
	if playerID < 0 || playerID >= len(m.players) {
		return nil, fmt.Errorf("invalid player id")
	}
	return m.players[playerID], nil
}

// Cell returns the requested cell.
// If accessed outside the grid, the default value (None) is returned.
func (m *WorldMap) Cell(xCol, yRow int) *Cell {
	// out of bound
	if len(m.grid) <= xCol || xCol < 0 {
		return NewCell(None, xCol, yRow)
	}
	if len(m.grid[xCol]) <= yRow || yRow < 0 {
		return NewCell(None, xCol, yRow)
	}

	// return cell
	return m.grid[xCol][yRow]
}

// CellByVector returns the cell pointed to by the specified coordinates.
func (m *WorldMap) CellByVector(v *Vector) *Cell {
	xCol := int(math.Floor(v.X() / CellSize))
	yRow := int(math.Floor(v.Y() / CellSize))
	return m.Cell(xCol, yRow)
}

// TouchingCells returns all cells touched by a ship.
// The ship is slightly smaller than a cell.
func (m *WorldMap) TouchingCells(v *Vector) []*Cell {
	r := CellRadius - CellRadius*0.1 - 1

	// get all possible cells
	all := make([]*Cell, 0, 8)
	all = append(all, m.CellByVector(NewVector(v.X()+r, v.Y()-r)))
	all = append(all, m.CellByVector(NewVector(v.X()-r, v.Y()+r)))
	all = append(all, m.CellByVector(NewVector(v.X()+r, v.Y()+r)))
	all = append(all, m.CellByVector(NewVector(v.X()-r, v.Y()-r)))
	all = append(all, m.CellByVector(NewVector(v.X()+r, v.Y())))
	all = append(all, m.CellByVector(NewVector(v.X()-r, v.Y())))
	all = append(all, m.CellByVector(NewVector(v.X(), v.Y()+r)))
	all = append(all, m.CellByVector(NewVector(v.X(), v.Y()-r)))

	// return (distinct)
	return removeDuplicate(all)
}

// FreeSpawn returns a random, free spawn point.
// A spawn is free when no player is touching it.
func (m *WorldMap) FreeSpawn() *Vector {
	// find spawn without other player
	free := make([]*Cell, 0, len(m.spawns)+1)

	for _, spawn := range m.spawns {
		isFree := true
		for _, player := range m.players {
			for _, cell := range m.TouchingCells(player.Position()) {
				if cell.Center().X() == spawn.Center().X() && cell.Center().Y() == spawn.Center().Y() {
					isFree = false
				}
			}
		}
		if isFree {
			free = append(free, spawn)
		}
	}

	// shuffle spawns
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })

	// add fallback spawn
	free = append(free, m.spawns[0]) // error fallback

	// return
	return free[0].Center().Clone()
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

// Freeze disable the world update.
func (m *WorldMap) Freeze(f bool) {
	m.freeze = f
}

// Update updates the world (move, score, velocity, ...).
// Call this several times per second in the background. (default 60/s)
func (m *WorldMap) Update() {

	// Freeze
	if m.freeze || m.iteration > m.endtime {
		return // no updates
	}

	// start timer
	start := time.Now()
	//--------------------------------------

	// do your thing
	for _, ship := range m.players {
		ship.Update()
	}

	// build protocol
	pOut := make([]byte, 0, 2000)
	pOut = append(pOut, []byte(ProtocolStatus(m))...)
	if m.iteration%2 == 0 {
		pOut = append(pOut, []byte(ProtocolPlayer(m))...)
	}
	if m.iteration%5 == 0 {
		pOut = append(pOut, []byte(ProtocolMap(m))...)
	}

	// send protocol to remote
	var wg sync.WaitGroup
	wg.Add(len(m.players))
	{ // go routines
		for _, p := range m.players {
			go func(p *Ship) {
				defer wg.Done()
				if p.remoteRW != nil {
					// write status
					if _, err := p.remoteRW.Write(pOut); err != nil {
						fmt.Printf("ERR remote player %s: %v\n", p.name, err)
						p.remoteRW = nil // disable remote
					}
				}
			}(p)
		}
	}
	wg.Wait() // WAITING

	// iteration
	m.iteration++

	//--------------------------------------
	// maxUpdateTime
	duration := time.Since(start)
	if m.maxUpdateTime.Microseconds() < duration.Microseconds() {
		m.maxUpdateTime = duration
		if m.maxUpdateTime > 16*time.Millisecond {
			// 16ms is fast enough for 60 updates per second
			fmt.Println("WARNING:", "maxUpdateTime", duration)
		}
	}
}

// AddPlayer registers and spawns a new player in the world.
// A unique name must be set.
// A valid color must be set (red, blue, green or orange).
// For remote param see Ship.Remote().
// There is a maximum number of players (see MaxPlayers).
func (m *WorldMap) AddPlayer(name, color string, remote io.ReadWriter) (playerID int, err error) {
	// check max player
	if len(m.players) >= m.MaxPlayers() {
		return -1, errors.New("maximum number of players reached")
	}

	// check name
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\r", "")
	if len(name) <= 0 || len(name) > 20 {
		return -1, errors.New("player name must be between 1 and 20 characters long")
	}

	// check double names
	for _, p := range m.players {
		if p.Name() == name {
			return -1, errors.New("player name already taken")
		}
	}

	// check color
	if color != "red" && color != "blue" && color != "green" && color != "orange" {
		return -1, errors.New("player color must be red, blue, green or orange")
	}

	// add / spawn
	playerID = len(m.players)
	ship := NewShip(m, playerID, name, color, remote)
	m.players = append(m.players, ship)

	// set spawn position
	ship.Spawn()

	// start move command listener
	if remote != nil {
		go func(r io.ReadWriter, p *Ship) {
			for {
				// prepare line reader
				tp := textproto.NewReader(bufio.NewReader(r))
				// read first line (ended with \n or \r\n)
				line, _ := tp.ReadLine()
				// parse param
				param := strings.Split(line, "|")
				if len(param) == 2 {
					x, errX := strconv.ParseFloat(param[0], 64)
					y, errY := strconv.ParseFloat(param[1], 64)
					if errX == nil && errY == nil {
						p.Move(NewVector(x, y))
					}
				}
			}
		}(remote, ship)
	}

	// return
	return playerID, nil
}

//--------  Helper  --------------------------------------------------------------------------------------------------//

// Print outputs the map on the console.
func (m *WorldMap) Print() {
	// top border
	fmt.Print("+")
	for xCol := 0; xCol < m.xWidth; xCol++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	// map
	for yRow := 0; yRow < m.yHeight; yRow++ {
		fmt.Print("|") // left border
		for xCol := 0; xCol < m.xWidth; xCol++ {
			fmt.Print(string(m.grid[xCol][yRow].Type())) // cell type
		}
		fmt.Print("|\n") // right border
	}

	// bottom border
	fmt.Print("+")
	for xCol := 0; xCol < m.xWidth; xCol++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")
}

// removeDuplicate removes cell duplicates from the list.
func removeDuplicate(in []*Cell) []*Cell {
	allKeys := make(map[string]bool)
	out := make([]*Cell, 0, len(in))
	for _, c := range in {
		key := fmt.Sprintf("%x|%d|%d", c.Type(), c.XCol(), c.YRow())
		if _, ok := allKeys[key]; !ok {
			allKeys[key] = true
			out = append(out, c)
		}
	}
	return out
}
