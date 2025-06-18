package attributes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Rect struct {
	Position Vector
	Size     Vector
}

func (r Rect) Center() Vector {
	return Vector{X: r.Size.X/2 + r.Position.X, Y: r.Size.Y/2 + r.Position.Y}
}

func (r Rect) TopRight() Vector {
	return Vector{X: r.Size.X + r.Position.X, Y: r.Position.Y}
}

func (r Rect) BottomLeft() Vector {
	return Vector{X: r.Position.X, Y: r.Size.Y + r.Position.Y}
}

func (r Rect) BottomRight() Vector {
	return Vector{X: r.Size.X + r.Position.X, Y: r.Size.Y + r.Position.Y}
}

func (r Rect) Top() float64 {
	return r.Position.Y
}

func (r Rect) Left() float64 {
	return r.Position.X
}

func (r Rect) Bottom() float64 {
	return r.BottomRight().Y
}

func (r Rect) Right() float64 {
	return r.BottomRight().X
}

func (r Rect) MidTop() Vector {
	return Vector{X: r.Left() + r.Size.X/2, Y: r.Top()}
}

func (r Rect) MidLeft() Vector {
	return Vector{X: r.Left(), Y: r.Top() + r.Size.Y/2}
}

func (r Rect) MidBottom() Vector {
	return Vector{X: r.Left() + r.Size.X/2, Y: r.Bottom()}
}

func (r Rect) MidRight() Vector {
	return Vector{X: r.Right(), Y: r.Top() + r.Size.Y/2}
}

func (r Rect) CollideRect(collisionRect Rect) bool {
	if r.Area() > collisionRect.Area() {
		return r.CollidePoint(collisionRect.Position) || r.CollidePoint(collisionRect.TopRight()) || r.CollidePoint(collisionRect.BottomLeft()) || r.CollidePoint(collisionRect.BottomRight())
	}
	return collisionRect.CollidePoint(r.Position) || collisionRect.CollidePoint(r.TopRight()) || collisionRect.CollidePoint(r.BottomLeft()) || collisionRect.CollidePoint(r.BottomRight())
}

func (r Rect) CollidePoint(point Vector) bool {
	if r.Position.X <= point.X && r.BottomRight().X >= point.X {
		if r.Position.Y <= point.Y && r.BottomRight().Y >= point.Y {
			return true
		}
	}
	return false
}

func (r Rect) Area() float64 {
	return r.Size.X * r.Size.Y
}

func (r Rect) Draw(surface *ebiten.Image, clr color.RGBA, offset Vector) {
	vector.DrawFilledRect(surface, float32(r.Position.X-offset.X), float32(r.Position.Y-offset.Y), float32(r.Size.X), float32(r.Size.Y), clr, false)
}
