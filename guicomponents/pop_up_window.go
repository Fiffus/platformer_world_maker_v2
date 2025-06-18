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
	var size attributes.Vector = attributes.Vector{X: 500, Y: 300}

	puw.rect = attributes.Rect{
		Position: attributes.Vector{X: float64(width)/2 - size.X/2, Y: float64(height)/2 - size.Y/2},
		Size:     size,
	}
	puw.confirm.Construct(
		attributes.Vector{
			X: puw.rect.Right(),
			Y: puw.rect.Bottom(),
		},
		"Confirm",
	)
	puw.confirm.rect.Position.X -= puw.confirm.rect.Size.X + 10
	puw.confirm.rect.Position.Y -= puw.confirm.rect.Size.Y + 10
	puw.cancel.Construct(
		attributes.Vector{
			X: puw.rect.Left(),
			Y: puw.rect.Bottom(),
		},
		"Cancel",
	)
	puw.cancel.rect.Position.X += 10
	puw.cancel.rect.Position.Y -= puw.cancel.rect.Size.Y + 10
	puw.entry.Construct(
		attributes.Vector{
			X: puw.rect.Left(),
			Y: puw.rect.MidLeft().Y,
		},
	)
	puw.entry.rect.Position.X = puw.rect.Center().X - puw.entry.rect.Size.X/2
	puw.entry.rect.Position.Y = puw.rect.Center().Y - puw.entry.rect.Size.Y/2
	puw.title.Construct(
		attributes.Vector{
			X: puw.rect.Left() + 10,
			Y: puw.rect.Top() + 10,
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

	puw.rect.Draw(surface, color.RGBA{100, 100, 120, 170}, attributes.Vector{X: 0, Y: 0})
	puw.confirm.Draw(surface)
	puw.cancel.Draw(surface)
	puw.entry.Draw(surface)
	puw.title.Draw(surface)
}
