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

type Entry struct {
	rect     attributes.Rect
	clr      attributes.Color
	text     string
	fontFace *text.GoTextFace
	entered  bool
}

func (e *Entry) Construct(position attributes.Vector) {
	e.rect = attributes.Rect{
		Position: position,
		Size: attributes.Vector{
			X: 200,
			Y: 60,
		},
	}
	e.clr = attributes.Color{
		Current:   color.RGBA{58, 55, 94, 255},
		Normal:    color.RGBA{58, 55, 94, 255},
		Highlight: color.RGBA{77, 73, 122, 255},
	}
	e.text = ""
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	e.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   40,
	}
	e.entered = false
}

func (e *Entry) Rect() attributes.Rect {
	return e.rect
}

func (e *Entry) IsActive() bool {
	return e.entered
}

func (e *Entry) enter() {
	var x, y int = ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !e.entered {
		if e.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
			e.entered = true
			e.clr.Current = e.clr.Highlight
			return
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && e.entered && !e.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		e.entered = false
		e.clr.Current = e.clr.Normal
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && e.entered {
		e.entered = false
		e.clr.Current = e.clr.Normal
	}
}

func (e *Entry) editText() {
	if e.entered {
		if len(ebiten.InputChars()) > 0 {
			e.text = fmt.Sprintf("%v%v", e.text, string(ebiten.InputChars()[len(ebiten.InputChars())-1]))
		}
		if len(e.text)-1 >= 0 && inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			e.text = e.text[:len(e.text)-1]
		}
	}
}

func (e *Entry) Value() string {
	return e.text
}

func (e *Entry) Update() {
	e.enter()
	e.editText()
}

func (e *Entry) Draw(surface *ebiten.Image) {
	e.rect.Draw(surface, e.clr.Current, attributes.Vector{X: 0, Y: 0})

	var textWidth, textHeight = text.Measure(e.text, e.fontFace, e.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(e.rect.Left()+e.rect.Size.X/2-float64(len(e.text))/2-textWidth/2, e.rect.Top()+e.rect.Size.Y/2-textHeight/2)
	options.ColorScale.Scale(1, 1, 1, 1)
	text.Draw(surface, e.text, e.fontFace, options)
}
