package guicomponents

import (
	"platformer_world_maker_v2/attributes"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DimensionChanger struct {
	entryX      Entry
	entryY      Entry
	validX      bool
	changeableX bool
	validY      bool
	changeableY bool
}

func (dc *DimensionChanger) Construct(position attributes.Spatial) {
	dc.entryX.Construct(position)
	dc.entryY.Construct(attributes.Spatial{X: position.X, Y: position.Y + dc.entryX.rect.Size.Y})
	dc.entryX.text = "230"
	dc.entryY.text = "160"
	dc.validX = true
	dc.changeableX = false
	dc.validY = true
	dc.changeableY = false
}

func (dc *DimensionChanger) ValidX() bool {
	return dc.validX
}

func (dc *DimensionChanger) ValidY() bool {
	return dc.validY
}

func (dc *DimensionChanger) ValueX() string {
	return dc.entryX.text
}

func (dc *DimensionChanger) ValueY() string {
	return dc.entryY.text
}

func (dc *DimensionChanger) ChangeableX() bool {
	return dc.changeableX
}

func (dc *DimensionChanger) ChangeableY() bool {
	return dc.changeableY
}

func (dc *DimensionChanger) ChangeUsed() {
	dc.changeableX = false
	dc.changeableY = false
}

func (dc *DimensionChanger) validateInputX() {
	if _, err := strconv.Atoi(dc.entryX.text); err != nil {
		dc.validX = false
		return
	}
	dc.validX = true
}

func (dc *DimensionChanger) validateInputY() {
	if _, err := strconv.Atoi(dc.entryY.text); err != nil {
		dc.validY = false
		return
	}
	dc.validY = true
}

func (dc *DimensionChanger) Update() {
	var x, y int = ebiten.CursorPosition()
	dc.validateInputX()
	dc.validateInputY()
	if dc.entryX.entered {
		if dc.validX && ((inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !dc.entryX.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)})) || ebiten.IsKeyPressed(ebiten.KeyEnter)) {
			dc.changeableX = true
		}
	}
	dc.entryX.enter()
	dc.entryX.editText()

	if dc.entryY.entered {
		if dc.validY && ((inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !dc.entryY.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)})) || ebiten.IsKeyPressed(ebiten.KeyEnter)) {
			dc.changeableY = true
		}
	}
	dc.entryY.enter()
	dc.entryY.editText()
}

func (dc *DimensionChanger) Draw(surface *ebiten.Image) {
	dc.entryX.Draw(surface)
	dc.entryY.Draw(surface)
}
