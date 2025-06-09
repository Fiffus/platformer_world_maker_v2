package attributes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Tools [][]Tool

type Tool struct {
	image            *ebiten.Image
	imageName        string
	rect             Rect
	originalPosition Spatial
}

func (t *Tool) Construct(rect Rect, texture *ebiten.Image, imageName string) {
	t.image = texture
	t.imageName = imageName
	t.rect = rect
	t.originalPosition = t.rect.Position
}

func (t *Tool) calculateScale() float64 {
	return float64(t.rect.Size.X / float64(t.image.Bounds().Max.X))
}

func (t *Tool) Pressed() bool {
	var x, y int = ebiten.CursorPosition()
	if t.rect.CollidePoint(Spatial{X: float64(x), Y: float64(y)}) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	return false
}

func (t *Tool) Offset(y float64) {
	t.rect.Position.Y = t.originalPosition.Y - y
}

func (t *Tool) Rect() Rect {
	return t.rect
}

func (t *Tool) GetCurrentImage() (*ebiten.Image, string) {
	return t.image, t.imageName
}

func (t *Tool) Draw(surface *ebiten.Image) {
	t.rect.Draw(surface, color.RGBA{58, 55, 94, 100}, Spatial{X: 0, Y: 0})
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(t.calculateScale(), t.calculateScale())
	options.GeoM.Translate(t.rect.Left(), t.rect.Top())
	surface.DrawImage(t.image, options)
}
