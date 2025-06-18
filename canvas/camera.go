package canvas

import (
	"platformer_world_maker_v2/attributes"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	offset attributes.Vector
	rect   attributes.Rect
	speed  float64
}

func (c *Camera) Construct(size attributes.Vector) {
	c.rect = attributes.Rect{
		Position: attributes.Vector{X: 0, Y: 0},
		Size:     size,
	}
	c.offset = c.rect.Position
	c.speed = 30
}

func (c *Camera) Move(bounds attributes.Vector) {
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && c.rect.Top()-c.speed >= 0 {
		c.rect.Position.Y -= c.speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)) && c.rect.Left()-c.speed >= 0 {
		c.rect.Position.X -= c.speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && c.rect.Bottom()+c.speed <= bounds.Y {
		c.rect.Position.Y += c.speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)) && c.rect.Right()+c.speed <= bounds.X {
		c.rect.Position.X += c.speed
	}
	c.scroll(bounds)
	c.offset = c.rect.Position
}

func (c *Camera) scroll(bounds attributes.Vector) {
	_, scroll := ebiten.Wheel()
	scroll *= 100

	mx, _ := ebiten.CursorPosition()
	lowestXbounds, _ := ebiten.Monitor().Size()
	lowestXbounds /= 4

	if mx < lowestXbounds {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && !(c.rect.Left()-float64(scroll) > 0) {
		c.rect.Position.X = 0
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && !(c.rect.Right()-float64(scroll) < bounds.X) {
		c.rect.Position.X = bounds.X - c.rect.Size.X
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && c.rect.Left()-float64(scroll) > 0 && c.rect.Right()-float64(scroll) < bounds.X {
		c.rect.Position.X -= float64(scroll)
		return
	}
	if !(c.rect.Top()-float64(scroll) > 0) {
		c.rect.Position.Y = 0
		return
	}
	if !(c.rect.Bottom()-float64(scroll) < bounds.Y) {
		c.rect.Position.Y = bounds.Y - c.rect.Size.Y
		return
	}
	if c.rect.Top()-float64(scroll) > 0 && c.rect.Bottom()-float64(scroll) < bounds.Y {
		c.rect.Position.Y -= float64(scroll)
	}
}

func (c *Camera) CheckBoundsAfterDimensionChange(bounds attributes.Vector) {
	if c.rect.Bottom() >= bounds.Y {
		c.rect.Position.Y -= c.rect.Bottom() - bounds.Y
	}
	if c.rect.Right() >= bounds.X {
		c.rect.Position.X -= c.rect.Right() - bounds.X
	}
}

func (c *Camera) Rect() attributes.Rect {
	return c.rect
}

func (c *Camera) Offset() attributes.Vector {
	return c.offset
}
