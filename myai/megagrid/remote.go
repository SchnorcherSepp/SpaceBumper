package megagrid

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// UpdateStreamHandler receives status updates continuously and updates the global variables
func (w *WorldMap) UpdateStreamHandler(in io.Reader) {
	const sStat = "START STATUS"
	const eStat = "END STATUS"
	const sPly = "START PLAYER"
	const ePly = "END PLAYER"
	const sMap = "START MAP"
	const eMap = "END MAP"

	// blocks
	text := new(strings.Builder)

	// read lines
	tp := textproto.NewReader(bufio.NewReader(in))
	for {
		line, err := tp.ReadLine()
		if err != nil {
			fmt.Printf("ERR: updateStreamHandler: %v\n", err)
			continue
		}

		// start new block (reset old block)
		if strings.HasPrefix(line, sStat) || strings.HasPrefix(line, sPly) || strings.HasPrefix(line, sMap) {
			text.Reset()
			continue
		}

		// block end found -> process STATUS
		if strings.HasPrefix(line, eStat) {
			w.parseStat(text.String())
			continue
		}
		// block end found -> process PLAYER
		if strings.HasPrefix(line, ePly) {
			w.parsePly(text.String())
			continue
		}
		// block end found -> process MAP
		if strings.HasPrefix(line, eMap) {
			w.parseMap(text.String())
			continue
		}

		// process block
		text.WriteString(line)
		text.WriteString("\n")
	}
}

// parse STATUS block
func (w *WorldMap) parseStat(text string) {
	w.Mux.Lock() // LOCK
	defer w.Mux.Unlock()

	for _, line := range strings.Split(text, "\n") {
		args := strings.Split(line, ":")
		if len(args) == 2 {
			if args[0] == "Iteration" {
				i, _ := strconv.ParseUint(strings.TrimSpace(args[1]), 10, 64)
				w.Iteration = i
			} else if args[0] == "Endtime" {
				i, _ := strconv.ParseUint(strings.TrimSpace(args[1]), 10, 64)
				w.EndTime = i
			} else if args[0] == "MaxUpdateTime" {
				d, _ := time.ParseDuration(strings.TrimSpace(args[1]))
				w.MaxUpdateTime = d
			}
		}
	}
}

// parse PLAYER block
func (w *WorldMap) parsePly(text string) {
	w.Mux.Lock() // LOCK
	defer w.Mux.Unlock()

	currentPlayerID := -1
	for _, line := range strings.Split(text, "\n") {
		// split elements
		els := strings.Split(line, "|")
		if len(els) == 10 {
			for _, el := range els {
				// split args
				args := strings.Split(el, ":")
				if len(args) == 2 {
					// process args
					if args[0] == "PlayerID" {
						currentPlayerID, _ = strconv.Atoi(strings.TrimSpace(args[1]))
						if len(w.Players) <= currentPlayerID {
							w.Players = append(w.Players, NewShip(w, currentPlayerID, "tmp", "red"))
						}
					} else if args[0] == "Name" {
						w.Players[currentPlayerID].Name = strings.TrimSpace(args[1])
					} else if args[0] == "Color" {
						w.Players[currentPlayerID].Color = strings.TrimSpace(args[1])
					} else if args[0] == "Position" {
						v := strings.Split(strings.TrimSpace(args[1]), ",")
						v1, _ := strconv.ParseFloat(v[0], 64)
						v2, _ := strconv.ParseFloat(v[1], 64)
						w.Players[currentPlayerID].Position = NewVector(v1, v2)
					} else if args[0] == "Velocity" {
						v := strings.Split(strings.TrimSpace(args[1]), ",")
						v1, _ := strconv.ParseFloat(v[0], 64)
						v2, _ := strconv.ParseFloat(v[1], 64)
						w.Players[currentPlayerID].Velocity = NewVector(v1, v2)
					} else if args[0] == "Acceleration" {
						v := strings.Split(strings.TrimSpace(args[1]), ",")
						v1, _ := strconv.ParseFloat(v[0], 64)
						v2, _ := strconv.ParseFloat(v[1], 64)
						w.Players[currentPlayerID].Acceleration = NewVector(v1, v2)
					} else if args[0] == "Score" {
						s, _ := strconv.Atoi(strings.TrimSpace(args[1]))
						w.Players[currentPlayerID].Score = s
					}
				}
			}
		}
	}
}

// parse MAP block
func (w *WorldMap) parseMap(text string) {
	w.Mux.Lock() // LOCK
	defer w.Mux.Unlock()

	// reset grid
	grid := make([][]*Cell, 0)

	// changes only
	needUpdate := false

	// update grid
	lines := strings.Split(text, "\n")
	for yH := 0; yH < len(lines); yH++ {
		line := lines[yH]
		if len(line) > 0 {
			for xW := 0; xW < len(line); xW++ {
				// extend grid
				if len(grid) <= xW {
					grid = append(grid, make([]*Cell, 0))
				}
				if len(grid[xW]) <= yH {
					grid[xW] = append(grid[xW], nil)
				}
				// add data
				typ := line[xW]
				grid[xW][yH] = NewCell(typ, xW, yH)
				if w.Grid.OutOfBound(xW, yH) || w.Grid[xW][yH].cType != typ {
					needUpdate = true
				}
			}
		}
	}

	// set world
	if needUpdate {
		w.Grid = grid
		w.MGrid = NewMegaGrid(grid)
		w.UpdateStars(w.PlayerName)
		println("set new grid and update neighbours and costs")
	}
}
