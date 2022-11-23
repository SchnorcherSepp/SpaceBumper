package resources

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
		Bg:     loadGameImg("game/bg.jpg"),
		Block:  loadGameImg("game/block.png"),
		Blue:   loadGameImg("game/blue.png"),
		Boost:  loadGameImg("game/boost.png"),
		Error:  loadGameImg("game/error.png"),
		Green:  loadGameImg("game/green.png"),
		Logo:   loadGameImg("game/logo.png"),
		Orange: loadGameImg("game/orange.png"),
		Red:    loadGameImg("game/red.png"),
		Slow:   loadGameImg("game/slow.png"),
		Spawn:  loadGameImg("game/spawn.png"),
		Star:   loadGameImg("game/star.png"),
		Anti:   loadGameImg("game/anti.png"),
		Tile:   loadGameImg("game/tile.png"),
	}
}

//go:embed game
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
