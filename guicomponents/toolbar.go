package guicomponents

import (
	"image/color"
	"platformer_world_maker_v2/attributes"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Toolbar struct {
	field                     *ebiten.Image
	rect                      attributes.Rect
	tools                     attributes.Tools
	scrollbar                 attributes.Rect
	scrollbarBackground       attributes.Rect
	scrollbarColor            color.RGBA
	scrollbarColorDefault     color.RGBA
	scrollbarColorHighlighted color.RGBA
	selectedNew               bool
	scrolling                 bool
	selectedImg               *ebiten.Image
	selectedName              string
	boundsY                   float64
}

func (tb *Toolbar) Construct(textures map[string]*ebiten.Image) {
	tb.loadTools(textures)
	tb.selectedNew = false
	var screenWidth, screenHeight int = ebiten.Monitor().Size()
	var width float64 = float64(screenWidth) * 0.217
	var height float64 = float64(screenHeight) * 0.5
	tb.rect = attributes.Rect{
		Position: attributes.Spatial{X: 0, Y: 0},
		Size:     attributes.Spatial{X: width, Y: height},
	}
	tb.scrollbarBackground = attributes.Rect{
		Position: attributes.Spatial{X: width - 35 + 10, Y: 0},
		Size:     attributes.Spatial{X: 35, Y: height},
	}
	tb.field = ebiten.NewImage(int(width), int(height))
	tb.scrollbarColorDefault = color.RGBA{210, 210, 210, 255}
	tb.scrollbarColorHighlighted = color.RGBA{255, 255, 255, 255}
	tb.scrollbarColor = tb.scrollbarColorDefault
	var lastRowImgHeight int = 0
	for i := range 4 {
		img, _ := tb.tools[len(tb.tools)-1][i].GetCurrentImage()
		if img != nil {
			if lastRowImgHeight < img.Bounds().Max.Y {
				lastRowImgHeight = img.Bounds().Max.Y
			}
		}
	}
	tb.boundsY = float64(len(tb.tools)-1)*66 + 10 + float64(len(tb.tools)-1)*20 + 100
	var scrollbarHeight float64 = height / (tb.boundsY / height)
	if scrollbarHeight < 0 {
		scrollbarHeight = height
	}
	tb.scrollbar = attributes.Rect{
		Position: attributes.Spatial{X: width - 35 + 10, Y: 0},
		Size:     attributes.Spatial{X: 35, Y: scrollbarHeight},
	}
	tb.scrolling = false
}

func (tb *Toolbar) loadTools(textures map[string]*ebiten.Image) {
	var extend int = 1
	if len(textures)%4 == 0 {
		extend = 0
	}
	tb.tools = make(attributes.Tools, len(textures)/4+extend)
	tb.tools[0] = make([]attributes.Tool, 4)
	var row int = 0
	var col int = 0
	var sortedTextures []string = make([]string, len(textures))
	var j int = 0
	for name := range textures {
		sortedTextures[j] = name
		j++
	}
	sort.Strings(sortedTextures)
	for i, imgName := range sortedTextures {
		if i%4 == 0 && i != 0 {
			row += 1
			tb.tools[row] = make([]attributes.Tool, 4)
		}
		col = i - row*4
		var tool attributes.Tool = attributes.Tool{}
		tool.Construct(
			attributes.Rect{
				Position: attributes.Spatial{X: float64(col*3)*22 + 10 + float64(col*2)*20, Y: float64(row)*66 + 10 + float64(row)*20},
				Size:     attributes.Spatial{X: 60, Y: 60},
			},
			textures[imgName],
			imgName,
		)
		tb.tools[row][col] = tool
	}
}

func (tb *Toolbar) SelectedImage() *ebiten.Image {
	return tb.selectedImg
}
func (tb *Toolbar) SelectedName() string {
	return tb.selectedName
}
func (tb *Toolbar) SelectedNew() bool {
	return tb.selectedNew
}

func (tb *Toolbar) Update() {
	tb.selectedNew = false
	tb.updateTools()
	tb.updateScrollBar()
}

func (tb *Toolbar) updateTools() {
	for i, row := range tb.tools {
		for j, tool := range row {
			img, name := tool.GetCurrentImage()
			if img == nil {
				continue
			}
			if tb.rect.CollideRect(tool.Rect()) {
				if tool.Pressed() {
					tb.selectedNew = true
					tb.selectedImg = img
					tb.selectedName = name
				}
			}
			tb.tools[i][j].Offset(tb.scrollbar.Top() * tb.boundsY / tb.rect.Size.Y)
		}
	}
}

func (tb *Toolbar) updateScrollBar() {
	var x, y int = ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if tb.scrollbar.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)}) {
			tb.scrolling = true
		}
	}
	if tb.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)}) {
		_, yScroll := ebiten.Wheel()
		yScroll *= 11
		if !(tb.scrollbar.Top()-float64(yScroll) > 0) {
			tb.scrollbar.Position.Y = 0
		} else if !(tb.scrollbar.Bottom()-float64(yScroll) < tb.rect.Bottom()) {
			tb.scrollbar.Position.Y = tb.scrollbarBackground.Bottom() - tb.scrollbar.Size.Y
		} else {
			if tb.scrollbar.Top()-float64(yScroll) > 0 && tb.scrollbar.Bottom()-float64(yScroll) < tb.rect.Bottom() {
				tb.scrollbar.Position.Y -= float64(yScroll)
			}
		}
	}
	if tb.scrollbarBackground.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)}) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			tb.scrollbar.Position.Y = float64(y) - tb.scrollbar.Size.Y/2
			if tb.scrollbar.Top() < 0 {
				tb.scrollbar.Position.Y = 0
			}
			if tb.scrollbar.Bottom() > tb.scrollbarBackground.Bottom() {
				tb.scrollbar.Position.Y = tb.scrollbarBackground.Bottom() - tb.scrollbar.Size.Y
			}
		}
	}
	if tb.scrolling {
		tb.scrollbarColor = tb.scrollbarColorHighlighted
		if float64(y)-tb.scrollbar.Size.Y/2 > 0 && float64(y)+tb.scrollbar.Size.Y/2 < tb.rect.Bottom() {
			tb.scrollbar.Position.Y = float64(y) - tb.scrollbar.Size.Y/2
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			tb.scrolling = false
		}
		return
	}
	tb.scrollbarColor = tb.scrollbarColorDefault
}
func (tb *Toolbar) Draw(surface *ebiten.Image) {
	tb.rect.Draw(tb.field, color.RGBA{20, 20, 23, 255}, attributes.Spatial{X: 0, Y: 0})
	for _, row := range tb.tools {
		for _, tool := range row {
			img, _ := tool.GetCurrentImage()
			if img != nil {
				if tb.rect.CollideRect(tool.Rect()) {
					tool.Draw(tb.field)
				}
			}
		}
	}
	tb.scrollbarBackground.Draw(tb.field, color.RGBA{58, 55, 94, 255}, attributes.Spatial{X: 0, Y: 0})
	tb.scrollbar.Draw(tb.field, tb.scrollbarColor, attributes.Spatial{X: 0, Y: 0})
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(1, 1)
	options.GeoM.Translate(10, 10)
	surface.DrawImage(tb.field, options)
}
