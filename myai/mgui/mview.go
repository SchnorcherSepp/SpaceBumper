package mgui

import (
	"SpaceBumper/myai/megagrid"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"sort"
)

// interface check: ebiten.Game
var _ ebiten.Game = (*GridView)(nil)

type GridView struct {
	screenWidth  int
	screenHeight int
	world        *megagrid.WorldMap
	callUpdate   bool
}

func RunGridView(title string, world *megagrid.WorldMap) error {
	xWidth, yHeight := world.Grid.Dimensions()

	gridView := &GridView{
		screenWidth:  xWidth * megagrid.CellSize,  // cell image 40x40
		screenHeight: yHeight * megagrid.CellSize, // cell image 40x40
		world:        world,
	}

	// config window
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon([]image.Image{Games.Logo})
	ebiten.SetWindowSize(gridView.screenWidth, gridView.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(30) // default: 60 ticks per second

	// run (BLOCKING)
	return ebiten.RunGame(gridView)
}

//--------------------------------------------------------------------------------------------------------------------//

func (g *GridView) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *GridView) Update() error {
	return nil
}

func (g *GridView) Draw(screen *ebiten.Image) {
	g.world.Mux.Lock()
	defer g.world.Mux.Unlock()

	// DRAW: background image
	drawBackground(screen, g.screenWidth, g.screenHeight)

	// DRAW: cell images (Map)
	dCost := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	// scaling
	grid := g.world.MGrid
	xWidth, yHeight := grid.Dimensions()
	drawGrid(screen, grid, !dCost, xWidth, yHeight)
	// draw Neighbours & path
	drawNeighbours(screen, grid)
	drawPath(screen, g.world)

	// draw player (ships)
	drawPlayer(screen, g.world)

	// draw debug (status)
	drawDebug(screen, g.world)
}

func drawDebug(screen *ebiten.Image, w *megagrid.WorldMap) {
	msg := fmt.Sprintf("\n  round=%d/%d, maxUpdateTime=%v\n", w.Iteration, w.EndTime, w.MaxUpdateTime)
	for i, p := range sortRPlayer(w.Players) {
		if p != nil {
			if p.Score > 0 {
				msg += fmt.Sprintf("  %d. %s: %d (Speed %.2f -> %.2f)\n", i+1, p.Name, p.Score, p.Acceleration.Length(), p.Velocity.Length())
			} else {
				msg += fmt.Sprintf("  %d. %s: dead\n", i+1, p.Name)
			}
		}
	}
	ebitenutil.DebugPrint(screen, msg)
}

func drawPlayer(screen *ebiten.Image, w *megagrid.WorldMap) {
	// DRAW: ships
	for _, s := range w.Players {
		op := new(ebiten.DrawImageOptions)

		// get ship image
		var sImg *ebiten.Image
		switch s.Color {
		case "red":
			sImg = Games.Red
		case "blue":
			sImg = Games.Blue
		case "green":
			sImg = Games.Green
		case "orange":
			sImg = Games.Orange
		default:
			sImg = Games.Error
		}

		// Move the image's center to the screen's upper-left corner.
		// This is a preparation for rotating. When geometry matrices are applied,
		// the origin point is the upper-left corner.
		w, h := sImg.Size()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

		// Rotate the image. As a result, the anchor point of this rotate is
		// the center of the image.
		//   90° × π/180 =1,571 rad
		angle := s.Velocity.Angle() + 4.71239 // add 270° to align the ship image
		op.GeoM.Rotate(angle)

		// Move the image to the final position.
		pos := s.Position
		op.GeoM.Translate(pos.X, pos.Y)

		// draw ship
		op.Filter = ebiten.FilterLinear // Specify linear filter.
		screen.DrawImage(sImg, op)

		// TEXT: Name
		name := fmt.Sprintf("%s (%d)", s.Name, s.Score)
		namePosX := pos.X - (6 / 2 * float64(len(name)))
		namePosY := pos.Y + 5 + megagrid.CellSize/2
		ebitenutil.DebugPrintAt(screen, name, int(namePosX), int(namePosY))
	}
}

func drawBackground(screen *ebiten.Image, screenWidth, screenHeight int) {
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Scale(float64(screenWidth)/2600.0, float64(screenHeight)/1839.0) // bgImage is 2600px * 1839px
	op.Filter = ebiten.FilterLinear                                          // Specify linear filter.
	screen.DrawImage(Games.Bg, op)
}

func drawNeighbours(screen *ebiten.Image, grid megagrid.MegaGrid) {
	x, y := ebiten.CursorPosition()
	cell := grid.CellByCoordinates(float64(x), float64(y))

	if cell.Neighbours() != nil {

		// Neighbours
		for i, c := range cell.Neighbours() {
			op := new(ebiten.DrawImageOptions)
			var sImg = Games.Error2

			// draw
			op.Filter = ebiten.FilterLinear // Specify linear filter.
			posX := float64(c.XCol() * 8)
			posY := float64(c.YRow() * 8)
			op.GeoM.Translate(posX, posY)
			screen.DrawImage(sImg, op)

			// TEXT: cost
			score := fmt.Sprintf("%d", cell.NeighboursCost()[i])
			scorePosX := posX - (6 / 2 * float64(len(score))) + (posX-float64(x))*4 + 20
			scorePosY := posY - 8 + (posY-float64(y))*4 + 20
			ebitenutil.DebugPrintAt(screen, score, int(scorePosX), int(scorePosY))
		}
	}

}

func drawPath(screen *ebiten.Image, world *megagrid.WorldMap) {

	// player position
	myPos := world.MyPos(world.PlayerName)
	op := new(ebiten.DrawImageOptions)
	var sImg = Games.Error
	op.Filter = ebiten.FilterLinear // Specify linear filter.
	posX := float64(myPos.XCol()) * 8
	posY := float64(myPos.YRow()) * 8
	op.GeoM.Translate(posX, posY)
	screen.DrawImage(sImg, op)

	// find first cell
	start := myPos

	// draw path
	for i := 0; i < 1000; i++ {
		if start == nil || start.AG <= 0 {
			break
		}

		op := new(ebiten.DrawImageOptions)
		var sImg = Games.Error
		op.Filter = ebiten.FilterLinear // Specify linear filter.
		posX := float64(start.XCol()) * 8
		posY := float64(start.YRow()) * 8
		op.GeoM.Translate(posX, posY)
		screen.DrawImage(sImg, op)

		// nxt
		start = start.AParent
	}
}

func drawGrid(screen *ebiten.Image, grid megagrid.MegaGrid, cost bool, xWidth, yHeight int) {

	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			cell := grid[xCol][yRow]

			// prepare image
			op := new(ebiten.DrawImageOptions)
			op.GeoM.Translate(float64(xCol*8), float64(yRow*8)) // cell image 40x40
			op.Filter = ebiten.FilterLinear                     // Specify linear filter.

			// color
			if cost {
				lvl := cell.Cost()
				if lvl < 0 {
					op.ColorM.Scale(1, 0, 0, 255)
				} else if lvl == 10 {
					op.ColorM.Scale(0, 1, 0, 1)
				} else if lvl > 9000 {
					op.ColorM.Scale(1, 0, 0, 1)
				} else if lvl > 10 {
					op.ColorM.Scale(0, 0, 1, float64(lvl)/50)
				}
			}

			// draw cell
			switch cell.Type() {
			case megagrid.Blocked:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Block, op)
			case megagrid.Boost:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Boost, op)
			case megagrid.Slow:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Slow, op)
			case megagrid.Tile:
				screen.DrawImage(Games.Tile, op)
			case megagrid.Star:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Star, op)
			case megagrid.Anti:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Anti, op)
			case megagrid.Spawn:
				screen.DrawImage(Games.Tile, op)
				screen.DrawImage(Games.Spawn, op)
			case megagrid.None:
				// draw nothing (= background image)
				if cost {
					screen.DrawImage(Games.Tile, op)
				}
			default:
				screen.DrawImage(Games.Error, op) // ERROR
			}
		}
	}
}

//--------------------------------------------------------------------------------------------------------------------//

func sortRPlayer(in []*megagrid.Ship) []*megagrid.Ship {
	// clone
	out := make([]*megagrid.Ship, 0, len(in))
	out = append(out, in...)

	// sort
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Score > out[j].Score
	})

	// return
	return out
}
