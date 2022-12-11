package gui

import (
	"SpaceBumper/core"
	"SpaceBumper/gui/resources"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"sort"
)

const GameSpeed = 60

// interface check: ebiten.Game
var _ ebiten.Game = (*Game)(nil)

// Game is the GUI
type Game struct {
	screenWidth  int
	screenHeight int
	world        *core.WorldMap
	callUpdate   bool
}

// RunGame starts a GUI window and displays the specified world.
// The callUpdate option activates the core.WorldMap Update() call with 60 Ticks per second.
// Do not activate this option if the update is done externally.
//
// This call is blocking.
func RunGame(title string, world *core.WorldMap, callUpdate bool) error {

	// config game
	game := &Game{
		screenWidth:  world.XWidth() * core.CellSize,  // cell image 40x40
		screenHeight: world.YHeight() * core.CellSize, // cell image 40x40
		world:        world,
		callUpdate:   callUpdate,
	}

	// config window
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowIcon([]image.Image{resources.Games.Logo})
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(GameSpeed) // default: 60 ticks per second

	// run (BLOCKING)
	return ebiten.RunGame(game)
}

//--------------------------------------------------------------------------------------------------------------------//

// Layout accepts a native outside size in device-independent pixels and returns the game's logical screen
// size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On browsers, the outside is a body
// element. On mobiles, the outside is the view's size.
//
// Even though the outside size and the screen size differ, the rendering scale is automatically adjusted to
// fit with the outside.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first frame.
//
// If Layout returns non-positive numbers, the caller can panic.
//
// You can return a fixed screen size if you don't care, or you can also return a calculated screen size
// adjusted with the given outside size.
func (g *Game) Layout(_, _ int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// Update updates a game by one tick. The given argument represents a screen image.
//
// Update updates only the game logic and Draw draws the screen.
//
// In the first frame, it is ensured that Update is called at least once before Draw. You can use Update
// to initialize the game state.
//
// After the first frame, Update might not be called or might be called once
// or more for one frame. The frequency is determined by the current TPS (tick-per-second).
func (g *Game) Update() error {

	// player control
	id := 0
	ship, err := g.world.Player(id)
	if err == nil && ship.IsAlive() && ship.Remote() == nil {
		// keys
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// cursor position
			x, y := ebiten.CursorPosition()
			// ship position
			pos := ship.Position()
			// calc acceleration vector
			cmd := core.NewVector((float64(x)-pos.X())/100, (float64(y)-pos.Y())/100)
			// set move command
			ship.Move(cmd)
		} else {
			// reset move command
			ship.Move(new(core.Vector))
		}
	}

	// call world update
	if g.callUpdate {
		g.world.Update()
	}

	// return
	return nil
}

// Draw draws the game screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted as the game screen.
func (g *Game) Draw(screen *ebiten.Image) {

	// DRAW: background image
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Scale(float64(g.screenWidth)/2600.0, float64(g.screenHeight)/1839.0) // bgImage is 2600px * 1839px
	op.Filter = ebiten.FilterLinear                                              // Specify linear filter.
	screen.DrawImage(resources.Games.Bg, op)

	// DRAW: cell images (Map)
	xWidth := g.world.XWidth()
	yHeight := g.world.YHeight()
	for xCol := 0; xCol < xWidth; xCol++ {
		for yRow := 0; yRow < yHeight; yRow++ {
			cell := g.world.Cell(xCol, yRow)

			// prepare image
			op := new(ebiten.DrawImageOptions)
			op.GeoM.Translate(float64(xCol*core.CellSize), float64(yRow*core.CellSize)) // cell image 40x40
			op.Filter = ebiten.FilterLinear                                             // Specify linear filter.

			// draw cell
			switch cell.Type() {
			case core.Blocked:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Block, op)
			case core.Boost:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Boost, op)
			case core.Slow:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Slow, op)
			case core.Tile:
				screen.DrawImage(resources.Games.Tile, op)
			case core.Star:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Star, op)
			case core.Anti:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Anti, op)
			case core.Spawn:
				screen.DrawImage(resources.Games.Tile, op)
				screen.DrawImage(resources.Games.Spawn, op)
			case core.None:
				// draw nothing (= background image)
			default:
				screen.DrawImage(resources.Games.Error, op) // ERROR
			}
		}
	}

	// DRAW: ships
	for _, s := range g.world.Players() {
		op := new(ebiten.DrawImageOptions)

		// get ship image
		var sImg *ebiten.Image
		switch s.Color() {
		case "red":
			sImg = resources.Games.Red
		case "blue":
			sImg = resources.Games.Blue
		case "green":
			sImg = resources.Games.Green
		case "orange":
			sImg = resources.Games.Orange
		default:
			sImg = resources.Games.Error
		}

		// Move the image's center to the screen's upper-left corner.
		// This is a preparation for rotating. When geometry matrices are applied,
		// the origin point is the upper-left corner.
		w, h := sImg.Size()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

		// Rotate the image. As a result, the anchor point of this rotate is
		// the center of the image.
		//   90° × π/180 =1,571 rad
		angle := s.Angle() + 4.71239 // add 270° to align the ship image
		op.GeoM.Rotate(angle)

		// Move the image to the final position.
		pos := s.Position()
		op.GeoM.Translate(pos.X(), pos.Y())

		// draw ship
		op.Filter = ebiten.FilterLinear // Specify linear filter.
		screen.DrawImage(sImg, op)

		// TEXT: debug messages
		iteration, endtime, maxUpdateTime := g.world.Stats()
		msg := fmt.Sprintf("\n  round=%d/%d, maxUpdateTime=%v\n", iteration, endtime, maxUpdateTime)
		for i, p := range sortPlayer(g.world.Players()) {
			if p != nil {
				if p.IsAlive() {
					msg += fmt.Sprintf("  %d. %s: %d (Speed %.2f -> %.2f)\n", i+1, p.Name(), p.Score(), p.Acceleration().Length(), p.Velocity().Length())
				} else {
					msg += fmt.Sprintf("  %d. %s: dead\n", i+1, p.Name())
				}
			}
		}
		ebitenutil.DebugPrint(screen, msg)

		// TEXT: Name
		name := fmt.Sprintf("%s", s.Name())
		namePosX := pos.X() - (6 / 2 * float64(len(name)))
		namePosY := pos.Y() + 5 + core.CellRadius
		ebitenutil.DebugPrintAt(screen, name, int(namePosX), int(namePosY))

		// TEXT: Score
		score := fmt.Sprintf("%d", s.Score())
		scorePosX := pos.X() - (6 / 2 * float64(len(score)))
		scorePosY := pos.Y() - 8
		ebitenutil.DebugPrintAt(screen, score, int(scorePosX), int(scorePosY))
	}
}

//--------------------------------------------------------------------------------------------------------------------//

func sortPlayer(in []*core.Ship) []*core.Ship {
	// clone
	out := make([]*core.Ship, 0, len(in))
	out = append(out, in...)

	// sort
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Score() > out[j].Score()
	})

	// return
	return out
}
