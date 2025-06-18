package guicomponents

import (
	"bytes"
	"image/color"
	"log"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type DropDownRow struct {
	rect     attributes.Rect
	clr      attributes.Color
	text     string
	fontFace *text.GoTextFace
	active   bool
}

func (ddr *DropDownRow) Construct(position attributes.Vector, rowText string) {
	ddr.rect = attributes.Rect{
		Position: position,
		Size: attributes.Vector{
			X: 140,
			Y: 60,
		},
	}
	ddr.text = rowText
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	ddr.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   40,
	}
	ddr.clr = attributes.Color{
		Current:   color.RGBA{58, 55, 94, 255},
		Normal:    color.RGBA{58, 55, 94, 255},
		Highlight: color.RGBA{77, 73, 122, 255},
	}
	ddr.active = false
}

func (ddr *DropDownRow) PressedLeft() bool {
	var x, y int = ebiten.CursorPosition()

	if ddr.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	return false
}

func (ddr *DropDownRow) PressedRight() bool {
	var x, y int = ebiten.CursorPosition()

	if ddr.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			return true
		}
	}
	return false
}

func (ddr *DropDownRow) highlight() {
	var x, y int = ebiten.CursorPosition()

	if ddr.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		ddr.clr.Current = ddr.clr.Highlight
	}

	if !ddr.active {
		ddr.clr.Current = ddr.clr.Normal
	}
}

func (ddr *DropDownRow) Activate() {
	ddr.active = true
}

func (ddr *DropDownRow) Dectivate() {
	ddr.active = false
}

func (ddr *DropDownRow) Update() {
	ddr.highlight()
}

func (ddr *DropDownRow) Draw(surface *ebiten.Image) {
	ddr.rect.Draw(surface, ddr.clr.Current, attributes.Vector{X: 0, Y: 0})

	var textWidth, textHeight = text.Measure(ddr.text, ddr.fontFace, ddr.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(ddr.rect.Left()+ddr.rect.Size.X/2-float64(len(ddr.text))/2-textWidth/2, ddr.rect.Top()+ddr.rect.Size.Y/2-textHeight/2)
	options.ColorScale.Scale(1, 1, 1, 1)
	text.Draw(surface, ddr.text, ddr.fontFace, options)
}
