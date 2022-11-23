package resources

import (
	"embed"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"image/color"
	_ "image/jpeg" // needed for ebitenutil.NewImageFromReader()
	_ "image/png"  // needed for ebitenutil.NewImageFromReader()
	"log"
	"strconv"
)

var Uis *UiResources

type UiResources struct {
	Bg             *image.NineSlice
	SeparatorColor color.Color
}

var Texts *TextResources

type TextResources struct {
	IdleColor     color.Color
	DisabledColor color.Color
	Face          font.Face
	TitleFace     font.Face
	BigTitleFace  font.Face
	SmallFace     font.Face
}

var Buttons *ButtonResources

type ButtonResources struct {
	Image   *widget.ButtonImage
	Text    *widget.ButtonTextColor
	Face    font.Face
	Padding widget.Insets
}

var Checkboxes *CheckboxResources

type CheckboxResources struct {
	Image   *widget.ButtonImage
	Graphic *widget.CheckboxGraphicImage
	Spacing int
}

var Labels *LabelResources

type LabelResources struct {
	Text *widget.LabelColor
	Face font.Face
}

var ComboButtons *ComboButtonResources

type ComboButtonResources struct {
	Image   *widget.ButtonImage
	Text    *widget.ButtonTextColor
	Face    font.Face
	Graphic *widget.ButtonImageImage
	Padding widget.Insets
}

var Lists *ListResources

type ListResources struct {
	Image        *widget.ScrollContainerImage
	Track        *widget.SliderTrackImage
	TrackPadding widget.Insets
	Handle       *widget.ButtonImage
	HandleSize   int
	Face         font.Face
	Entry        *widget.ListEntryColor
	EntryPadding widget.Insets
}

var Sliders *SliderResources

type SliderResources struct {
	TrackImage *widget.SliderTrackImage
	Handle     *widget.ButtonImage
	HandleSize int
}

var Panels *PanelResources

type PanelResources struct {
	Image   *image.NineSlice
	Padding widget.Insets
}

var TabBooks *TabBookResources

type TabBookResources struct {
	IdleButton     *widget.ButtonImage
	SelectedButton *widget.ButtonImage
	ButtonFace     font.Face
	ButtonText     *widget.ButtonTextColor
	ButtonPadding  widget.Insets
}

var Headers *HeaderResources

type HeaderResources struct {
	Bg      *image.NineSlice
	Padding widget.Insets
	Face    font.Face
	Color   color.Color
}

var TextInputs *TextInputResources

type TextInputResources struct {
	Image   *widget.TextInputImage
	Padding widget.Insets
	Face    font.Face
	Color   *widget.TextInputColor
}

var ToolTips *ToolTipResources

type ToolTipResources struct {
	Bg      *image.NineSlice
	Padding widget.Insets
	Face    font.Face
	Color   color.Color
}

//--------------------------------------------------------------------------------------------------------------------//

func init() {

	const (
		bgColor = "131a22"

		textIdleColor     = "dff4ff"
		textDisabledColor = "5a7a91"

		labelIdleColor     = textIdleColor
		labelDisabledColor = textDisabledColor

		buttonIdleColor     = textIdleColor
		buttonDisabledColor = labelDisabledColor

		listSelectedBg         = "4b687a"
		listDisabledSelectedBg = "2a3944"

		headerColor = textIdleColor

		textInputCaretColor         = "e7c34b"
		textInputDisabledCaretColor = "766326"

		toolTipColor = bgColor

		separatorColor = listDisabledSelectedBg
	)

	// -----------------------------

	Uis = &UiResources{
		Bg:             image.NewNineSliceColor(hexToColor(bgColor)),
		SeparatorColor: hexToColor(separatorColor),
	}

	// -----------------------------

	Texts = &TextResources{
		IdleColor:     hexToColor(textIdleColor),
		DisabledColor: hexToColor(textDisabledColor),
		Face:          Fonts.Face,
		TitleFace:     Fonts.TitleFace,
		BigTitleFace:  Fonts.BigTitleFace,
		SmallFace:     Fonts.ToolTipFace,
	}

	// -----------------------------

	Buttons = &ButtonResources{
		Image: &widget.ButtonImage{
			Idle:     loadNineSlice("gui/button-idle.png", 12, 0),
			Hover:    loadNineSlice("gui/button-hover.png", 12, 0),
			Pressed:  loadNineSlice("gui/button-pressed.png", 12, 0),
			Disabled: loadNineSlice("gui/button-disabled.png", 12, 0),
		},
		Text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},
		Face: Fonts.Face,
		Padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}

	// -----------------------------

	Checkboxes = &CheckboxResources{
		Image: &widget.ButtonImage{
			Idle:     loadNineSlice("gui/checkbox-idle.png", 20, 0),
			Hover:    loadNineSlice("gui/checkbox-hover.png", 20, 0),
			Pressed:  loadNineSlice("gui/checkbox-hover.png", 20, 0),
			Disabled: loadNineSlice("gui/checkbox-disabled.png", 20, 0),
		},
		Graphic: &widget.CheckboxGraphicImage{
			Checked:   loadButtonImg("gui/checkbox-checked-idle.png", "gui/checkbox-checked-disabled.png"),
			Unchecked: loadButtonImg("gui/checkbox-unchecked-idle.png", "gui/checkbox-unchecked-disabled.png"),
			Greyed:    loadButtonImg("gui/checkbox-greyed-idle.png", "gui/checkbox-greyed-disabled.png"),
		},
		Spacing: 10,
	}

	// -----------------------------

	Labels = &LabelResources{
		Text: &widget.LabelColor{
			Idle:     hexToColor(labelIdleColor),
			Disabled: hexToColor(labelDisabledColor),
		},

		Face: Fonts.Face,
	}

	// -----------------------------

	ComboButtons = &ComboButtonResources{
		Image: &widget.ButtonImage{
			Idle:     loadNineSlice("gui/combo-button-idle.png", 12, 0),
			Hover:    loadNineSlice("gui/combo-button-hover.png", 12, 0),
			Pressed:  loadNineSlice("gui/combo-button-pressed.png", 12, 0),
			Disabled: loadNineSlice("gui/combo-button-disabled.png", 12, 0),
		},
		Text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},
		Face:    Fonts.Face,
		Graphic: loadButtonImg("gui/arrow-down-idle.png", "gui/arrow-down-disabled.png"),
		Padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}
	// -----------------------------

	Panels = &PanelResources{
		Image: loadNineSlice("gui/panel-idle.png", 10, 10),
		Padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}

	// -----------------------------

	TabBooks = &TabBookResources{
		SelectedButton: &widget.ButtonImage{
			Idle:     loadNineSlice("gui/button-selected-idle.png", 12, 0),
			Hover:    loadNineSlice("gui/button-selected-hover.png", 12, 0),
			Pressed:  loadNineSlice("gui/button-selected-pressed.png", 12, 0),
			Disabled: loadNineSlice("gui/button-selected-disabled.png", 12, 0),
		},
		IdleButton: &widget.ButtonImage{
			Idle:     loadNineSlice("gui/button-idle.png", 12, 0),
			Hover:    loadNineSlice("gui/button-hover.png", 12, 0),
			Pressed:  loadNineSlice("gui/button-pressed.png", 12, 0),
			Disabled: loadNineSlice("gui/button-disabled.png", 12, 0),
		},
		ButtonFace: Fonts.Face,
		ButtonText: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},
		ButtonPadding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}

	// -----------------------------

	Headers = &HeaderResources{
		Bg: loadNineSlice("gui/header.png", 446, 9),
		Padding: widget.Insets{
			Left:   25,
			Right:  25,
			Top:    4,
			Bottom: 4,
		},
		Face:  Fonts.BigTitleFace,
		Color: hexToColor(headerColor),
	}

	// -----------------------------

	ToolTips = &ToolTipResources{
		Bg: image.NewNineSlice(loadGuiImg("gui/tool-tip.png"), [3]int{19, 6, 13}, [3]int{19, 5, 13}),
		Padding: widget.Insets{
			Left:   15,
			Right:  15,
			Top:    10,
			Bottom: 10,
		},
		Face:  Fonts.ToolTipFace,
		Color: hexToColor(toolTipColor),
	}

	// -----------------------------

	TextInputs = &TextInputResources{
		Image: &widget.TextInputImage{
			Idle:     image.NewNineSlice(loadGuiImg("gui/text-input-idle.png"), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(loadGuiImg("gui/text-input-disabled.png"), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		},
		Padding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		},
		Face: Fonts.Face,
		Color: &widget.TextInputColor{
			Idle:          hexToColor(textIdleColor),
			Disabled:      hexToColor(textDisabledColor),
			Caret:         hexToColor(textInputCaretColor),
			DisabledCaret: hexToColor(textInputDisabledCaretColor),
		},
	}

	// -----------------------------

	Lists = &ListResources{
		Image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(loadGuiImg("gui/list-idle.png"), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(loadGuiImg("gui/list-disabled.png"), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(loadGuiImg("gui/list-mask.png"), [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},
		Track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(loadGuiImg("gui/list-track-idle.png"), [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(loadGuiImg("gui/list-track-idle.png"), [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(loadGuiImg("gui/list-track-disabled.png"), [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},
		TrackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},
		Handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-idle.png"), 0, 5),
			Hover:    image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-hover.png"), 0, 5),
			Pressed:  image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-hover.png"), 0, 5),
			Disabled: image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-idle.png"), 0, 5),
		},
		HandleSize: 5,
		Face:       Fonts.Face,
		Entry: &widget.ListEntryColor{
			Unselected:         hexToColor(textIdleColor),
			DisabledUnselected: hexToColor(textDisabledColor),

			Selected:         hexToColor(textIdleColor),
			DisabledSelected: hexToColor(textDisabledColor),

			SelectedBackground:         hexToColor(listSelectedBg),
			DisabledSelectedBackground: hexToColor(listDisabledSelectedBg),
		},
		EntryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}

	// -----------------------------

	Sliders = &SliderResources{
		TrackImage: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(loadGuiImg("gui/slider-track-idle.png"), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(loadGuiImg("gui/slider-track-idle.png"), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(loadGuiImg("gui/slider-track-disabled.png"), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		},
		Handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-idle.png"), 0, 5),
			Hover:    image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-hover.png"), 0, 5),
			Pressed:  image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-hover.png"), 0, 5),
			Disabled: image.NewNineSliceSimple(loadGuiImg("gui/slider-handle-disabled.png"), 0, 5),
		},
		HandleSize: 6,
	}
}

//--------------------------------------------------------------------------------------------------------------------//

//go:embed gui
var guFS embed.FS

func loadGuiImg(name string) *ebiten.Image {
	// open reader
	r, err := guFS.Open(name)
	if err != nil {
		log.Fatalf("err: loadGuiImg: %v\n", err)
	}
	// get image
	eim, _, err := ebitenutil.NewImageFromReader(r)
	if err != nil {
		log.Fatalf("err: loadGuiImg: %v\n", err)
	}
	// return
	return eim
}

func loadButtonImg(idle string, disabled string) *widget.ButtonImageImage {
	idleImage := loadGuiImg(idle)

	var disabledImage *ebiten.Image
	if disabled != "" {
		disabledImage = loadGuiImg(disabled)
	}

	return &widget.ButtonImageImage{
		Idle:     idleImage,
		Disabled: disabledImage,
	}
}

func loadNineSlice(path string, centerWidth int, centerHeight int) *image.NineSlice {
	img := loadGuiImg(path)
	w, h := img.Size()

	return image.NewNineSlice(img,
		[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
		[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight},
	)
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		log.Fatalf("err: hexToColor: %v\n", err)
	}

	return color.RGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
