package guicomponents

import (
	"bytes"
	"log"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	text     string
	fontFace *text.GoTextFace
	position attributes.Vector
}

func (l *Label) Construct(position attributes.Vector, labelText string) {
	l.text = labelText
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	l.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   40,
	}
	l.position = position
}

func (l *Label) SetText(newText string) {
	l.text = newText
}

func (l *Label) Draw(surface *ebiten.Image) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(l.position.X, l.position.Y)
	options.ColorScale.Scale(1, 1, 1, 1)
	text.Draw(surface, l.text, l.fontFace, options)
}
