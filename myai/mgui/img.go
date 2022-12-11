package mgui

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg" // needed for ebitenutil.NewImageFromReader()
	_ "image/png"  // needed for ebitenutil.NewImageFromReader()
	"log"
)

var Games *GameResources

type GameResources struct {
	Bg     *ebiten.Image
	Block  *ebiten.Image
	Blue   *ebiten.Image
	Boost  *ebiten.Image
	Error  *ebiten.Image
	Error2 *ebiten.Image
	Green  *ebiten.Image
	Logo   *ebiten.Image
	Orange *ebiten.Image
	Red    *ebiten.Image
	Slow   *ebiten.Image
	Spawn  *ebiten.Image
	Star   *ebiten.Image
	Anti   *ebiten.Image
	Tile   *ebiten.Image
}

func init() {
	Games = &GameResources{
		Bg:     loadGameImg("mres/bg.jpg"),
		Block:  loadGameImg("mres/block.png"),
		Blue:   loadGameImg("mres/blue.png"),
		Boost:  loadGameImg("mres/boost.png"),
		Error:  loadGameImg("mres/error.png"),
		Error2: loadGameImg("mres/error2.png"),
		Green:  loadGameImg("mres/green.png"),
		Logo:   loadGameImg("mres/logo.png"),
		Orange: loadGameImg("mres/orange.png"),
		Red:    loadGameImg("mres/red.png"),
		Slow:   loadGameImg("mres/slow.png"),
		Spawn:  loadGameImg("mres/spawn.png"),
		Star:   loadGameImg("mres/star.png"),
		Anti:   loadGameImg("mres/anti.png"),
		Tile:   loadGameImg("mres/tile.png"),
	}
}

//go:embed mres
var gFS embed.FS

func loadGameImg(name string) *ebiten.Image {
	// open reader
	r, err := gFS.Open(name)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// get image
	eim, _, err := ebitenutil.NewImageFromReader(r)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// return
	return eim
}
