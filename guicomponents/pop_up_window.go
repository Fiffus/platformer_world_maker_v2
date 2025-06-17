package guicomponents

import (
	"image/color"
	"platformer_world_maker_v2/attributes"

	"github.com/hajimehoshi/ebiten/v2"
)

type PopUpWindow struct {
	rect    attributes.Rect
	confirm Button
	cancel  Button
	entry   Entry
	title   Label
	active  bool
}

func (puw *PopUpWindow) Construct(title, initialEntryValue string) {
	var width, height int = ebiten.Monitor().Size()
	var size attributes.Spatial = attributes.Spatial{X: 500, Y: 300}

	puw.rect = attributes.Rect{
		Position: attributes.Spatial{X: float64(width)/2 - size.X/2, Y: float64(height)/2 - size.Y/2},
		Size:     size,
	}
	puw.confirm.Construct(
		attributes.Spatial{
			X: puw.rect.Right(),
			Y: puw.rect.Bottom(),
		},
		"Confirm",
	)
	puw.cancel.Construct(
		attributes.Spatial{
			X: puw.rect.Left(),
			Y: puw.rect.Bottom(),
		},
		"Cancel",
	)
	puw.entry.Construct(
		attributes.Spatial{
			X: puw.rect.Left(),
			Y: puw.rect.MidLeft().Y,
		},
	)
	puw.title.Construct(
		attributes.Spatial{
			X: puw.rect.Left(),
			Y: puw.rect.Top(),
		},
		title,
	)
	puw.active = false
}

func (puw *PopUpWindow) Confirmed() bool {
	return puw.confirm.Pressed()
}

func (puw *PopUpWindow) Cancelled() bool {
	return puw.cancel.Pressed()
}

func (puw *PopUpWindow) Value() string {
	return puw.entry.text
}

func (puw *PopUpWindow) Activate() {
	puw.active = true
}

func (puw *PopUpWindow) Active() bool {
	return puw.active
}

func (puw *PopUpWindow) Update() {
	if !puw.active {
		return
	}

	puw.confirm.HighLight()
	puw.cancel.HighLight()
	puw.entry.Update()

	if puw.confirm.Pressed() {
		puw.active = false
	}
	if puw.cancel.Pressed() {
		puw.active = false
	}
}

func (puw *PopUpWindow) Draw(surface *ebiten.Image) {
	if !puw.active {
		return
	}

	puw.rect.Draw(surface, color.RGBA{100, 100, 120, 170}, attributes.Spatial{X: 0, Y: 0})
	puw.confirm.Draw(surface)
	puw.cancel.Draw(surface)
	puw.entry.Draw(surface)
	puw.title.Draw(surface)
}
