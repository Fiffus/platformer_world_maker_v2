package canvas

import (
	"image/color"
	"platformer_world_maker_v2/attributes"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cursor struct {
	rect  attributes.Rect
	color color.RGBA
}

func (c *Cursor) Construct(rect attributes.Rect) {
	c.rect = rect
	c.color = color.RGBA{58, 55, 94, 80}
}

func (c *Cursor) move() {
	var x, y int = ebiten.CursorPosition()
	c.rect.Position = attributes.Spatial{X: float64(x) - c.rect.Size.X/2, Y: float64(y) - c.rect.Size.Y/2}
}

func (c *Cursor) changeSize() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) && ebiten.IsKeyPressed(ebiten.KeyControlLeft) {
		c.rect.Size.X += 3
		c.rect.Size.Y += 3
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) && c.rect.Size.X-3 > 33 {
		c.rect.Size.X -= 3
		c.rect.Size.Y -= 3
	}
}

func (c *Cursor) Rect() attributes.Rect {
	return c.rect
}

func (c *Cursor) Update() {
	c.changeSize()
	c.move()
}

func (c *Cursor) Draw(surface *ebiten.Image) {
	c.rect.Draw(surface, c.color, attributes.Spatial{X: 0, Y: 0})
}
