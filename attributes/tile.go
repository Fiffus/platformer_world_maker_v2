package attributes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	rect      Rect
	imageName string
	image     *ebiten.Image
}

func (t *Tile) Construct(position Vector, baseTileSize float64, textureName string, image *ebiten.Image) {
	t.rect = Rect{
		Position: position,
		Size:     Vector{X: baseTileSize, Y: baseTileSize},
	}
	t.imageName = textureName
	t.image = image
}

func (t *Tile) Rect() Rect {
	return t.rect
}

func (t *Tile) ImageName() string {
	return t.imageName
}

func (t *Tile) Image() *ebiten.Image {
	return t.image
}

func (t *Tile) SetImage(newImageName string, newImage *ebiten.Image) {
	t.imageName = newImageName
	t.image = newImage
}

func (t *Tile) Draw(surface *ebiten.Image, opacity float64, offset Vector) {
	options := &ebiten.DrawImageOptions{}
	options.ColorScale.Scale(1, 1, 1, float32(opacity))
	options.GeoM.Scale(SCALE, SCALE)
	var topMargin float64 = float64(int(float64(t.image.Bounds().Max.Y)*SCALE) % int(t.rect.Size.Y))
	if float64(t.image.Bounds().Max.Y)*SCALE < t.rect.Size.Y {
		topMargin = -(t.rect.Size.Y - topMargin)
	}
	if topMargin != 0 {
		topMargin -= SCALE
	}
	var leftMargin float64 = t.rect.Center().X - ((float64(t.image.Bounds().Max.X) * SCALE) / 2)
	options.GeoM.Translate(
		leftMargin-offset.X,
		t.rect.Top()-offset.Y-topMargin,
	)
	surface.DrawImage(t.image, options)
}
