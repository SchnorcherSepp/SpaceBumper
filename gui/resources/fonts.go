package resources

import (
	"embed"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	_ "image/jpeg" // needed for ebitenutil.NewImageFromReader()
	_ "image/png"  // needed for ebitenutil.NewImageFromReader()
	"log"
)

var Fonts *FontResources

type FontResources struct {
	Face         font.Face
	TitleFace    font.Face
	BigTitleFace font.Face
	ToolTipFace  font.Face
}

func init() {
	const fontFaceRegular = "fonts/NotoSans-Regular.ttf"
	const fontFaceBold = "fonts/NotoSans-Bold.ttf"

	Fonts = &FontResources{
		Face:         loadFont(fontFaceRegular, 20),
		TitleFace:    loadFont(fontFaceBold, 24),
		BigTitleFace: loadFont(fontFaceBold, 28),
		ToolTipFace:  loadFont(fontFaceRegular, 15),
	}
}

//go:embed fonts
var fFS embed.FS

func loadFont(name string, size float64) font.Face {
	// read bytes
	fontData, err := fFS.ReadFile(name)
	if err != nil {
		log.Fatalf("err: loadFont: %v\n", err)
	}

	// parse font
	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalf("err: loadFont: %v\n", err)
	}

	// return face
	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
