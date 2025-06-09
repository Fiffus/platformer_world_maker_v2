package guicomponents

import (
	"platformer_world_maker_v2/attributes"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DimensionChanger struct {
	entry      Entry
	valid      bool
	changeable bool
}

func (dc *DimensionChanger) Construct(position attributes.Spatial) {
	dc.entry.Construct(position)
	dc.entry.text = "230x160"
	dc.valid = false
	dc.changeable = false
}

func (dc *DimensionChanger) Valid() bool {
	return dc.valid
}

func (dc *DimensionChanger) Value() string {
	return dc.entry.text
}

func (dc *DimensionChanger) Changeable() bool {
	return dc.changeable
}

func (dc *DimensionChanger) ChangeUsed() {
	dc.changeable = false
}

func (dc *DimensionChanger) validateInput() {
	if dc.entry.entered {
		var dimensions []string = strings.Split(dc.entry.text, "x")
		if dc.entry.text == "" {
			dc.valid = false
			return
		}
		if len(dimensions) == 2 {
			if _, err := strconv.Atoi(dimensions[0]); err != nil {
				dc.valid = false
				return
			}
			if _, err := strconv.Atoi(dimensions[1]); err != nil {
				dc.valid = false
				return
			}
			dc.valid = true
		}
	}
}

func (dc *DimensionChanger) Update() {
	var x, y int = ebiten.CursorPosition()
	dc.validateInput()
	if dc.entry.entered {
		if dc.valid && ((inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !dc.entry.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)})) || ebiten.IsKeyPressed(ebiten.KeyEnter)) {
			dc.changeable = true
		}
	}
	dc.entry.enter()
	dc.entry.editText()
}

func (dc *DimensionChanger) Draw(surface *ebiten.Image) {
	dc.entry.Draw(surface)
}
