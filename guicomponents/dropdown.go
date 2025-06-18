package guicomponents

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type DropDown struct {
	rect    attributes.Rect
	addRect attributes.Rect

	clr    attributes.Color
	addClr attributes.Color

	text     string
	fontFace *text.GoTextFace

	rows      []DropDownRow
	activeRow int

	opened          bool
	justAddedNewRow bool
}

func (dd *DropDown) Construct(position attributes.Vector) {
	dd.rect = attributes.Rect{
		Position: position,
		Size: attributes.Vector{
			X: 140,
			Y: 60,
		},
	}
	dd.addRect = attributes.Rect{
		Position: attributes.Vector{X: position.X + dd.rect.Size.X, Y: position.Y},
		Size:     attributes.Vector{X: 60, Y: 60},
	}
	dd.text = "Layer 0"
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	dd.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   40,
	}
	dd.clr = attributes.Color{
		Current:   color.RGBA{58, 55, 94, 255},
		Normal:    color.RGBA{58, 55, 94, 255},
		Highlight: color.RGBA{77, 73, 122, 255},
	}
	dd.addClr = dd.clr
	var firstRow DropDownRow = DropDownRow{}
	firstRow.Construct(attributes.Vector{X: position.X, Y: position.Y + dd.rect.Size.Y}, "Layer 0")
	firstRow.Activate()
	dd.rows = []DropDownRow{firstRow}
	dd.activeRow = 0
	dd.opened = false
	dd.justAddedNewRow = false
}

func (dd *DropDown) LoadLayers(numLayers int) {
	dd.rows = make([]DropDownRow, numLayers)
	for i := range numLayers {
		dd.rows[i].Construct(attributes.Vector{X: dd.rect.Position.X, Y: dd.rect.Position.Y + dd.rect.Size.Y*float64(i+1)}, fmt.Sprintf("Layer %d", i))
	}
	dd.rows[0].Activate()
}

func (dd *DropDown) highLightAdd() {
	var x, y int = ebiten.CursorPosition()

	if dd.addRect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		dd.addClr.Current = dd.addClr.Highlight
		return
	}
	dd.addClr.Current = dd.addClr.Normal
}

func (dd *DropDown) highLight() {
	var x, y int = ebiten.CursorPosition()

	if dd.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		dd.clr.Current = dd.clr.Highlight
		return
	}
	dd.clr.Current = dd.clr.Normal
}

func (dd *DropDown) addNewRow() {
	var x, y int = ebiten.CursorPosition()

	if dd.addRect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			var newRow DropDownRow = DropDownRow{}
			newRow.Construct(attributes.Vector{X: dd.rect.Position.X, Y: dd.rect.Position.Y + float64((len(dd.rows)+1)*60)}, fmt.Sprintf("Layer %d", len(dd.rows)))
			dd.rows = append(dd.rows, newRow)
			dd.justAddedNewRow = true
		}
	}
}

func (dd *DropDown) ActiveRow() int {
	return dd.activeRow
}

func (dd *DropDown) Active() bool {
	return dd.opened
}

func (dd *DropDown) JustAddedNewRow() bool {
	return dd.justAddedNewRow
}

func (dd *DropDown) toggle() {
	var x, y int = ebiten.CursorPosition()

	if !dd.opened {
		if dd.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				dd.opened = true
				return
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if dd.addRect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
			return
		}
		for i := range dd.rows {
			if dd.rows[i].rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
				return
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		dd.opened = false
	}
}

func (dd *DropDown) Update() {
	dd.justAddedNewRow = false
	dd.toggle()
	dd.highLight()
	dd.highLightAdd()
	dd.addNewRow()

	if !dd.opened {
		return
	}
	for i := range dd.rows {
		dd.rows[i].Update()

		if dd.rows[i].PressedLeft() {
			dd.rows[dd.activeRow].Dectivate()
			dd.rows[i].Activate()
			dd.activeRow = i
			dd.text = dd.rows[i].text
		}
	}
}

func (dd *DropDown) Draw(surface *ebiten.Image) {
	dd.rect.Draw(surface, dd.clr.Current, attributes.Vector{X: 0, Y: 0})
	dd.addRect.Draw(surface, dd.addClr.Current, attributes.Vector{X: 0, Y: 0})

	var textWidth, textHeight = text.Measure(dd.text, dd.fontFace, dd.fontFace.Size+10)

	options := text.DrawOptions{}
	options.GeoM.Translate(dd.rect.Left()+dd.rect.Size.X/2-float64(len(dd.text))/2-textWidth/2, dd.rect.Top()+dd.rect.Size.Y/2-textHeight/2)
	options.ColorScale.Scale(1, 1, 1, 1)
	text.Draw(surface, dd.text, dd.fontFace, &options)

	var addWidth, addHeight = text.Measure("+", dd.fontFace, dd.fontFace.Size+10)

	addOptions := text.DrawOptions{}
	addOptions.GeoM.Translate(dd.addRect.Left()+dd.addRect.Size.X/2-1.0/2.0-addWidth/2, dd.addRect.Top()+dd.addRect.Size.Y/2-addHeight/2)
	text.Draw(surface, "+", dd.fontFace, &addOptions)

	if !dd.opened {
		return
	}
	for i := range dd.rows {
		dd.rows[i].Draw(surface)
	}
}
