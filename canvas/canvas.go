package canvas

import (
	"image/color"
	"math"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/loader"

	"github.com/hajimehoshi/ebiten/v2"
)

type Canvas struct {
	canvas       *ebiten.Image
	position     attributes.Spatial
	camera       Camera
	Layers       []attributes.Layer
	activeLayer  int
	baseTileSize float64
}

func (c *Canvas) Contruct(images map[string]*ebiten.Image) {
	var screenWidth, screenHeight int = ebiten.Monitor().Size()
	var width float64 = float64(screenWidth) * 0.75

	c.canvas = ebiten.NewImage(int(width), screenHeight)
	c.position = attributes.Spatial{X: float64(screenWidth) - width, Y: 0}
	c.camera = Camera{}
	c.camera.Construct(attributes.Spatial{X: width, Y: float64(screenHeight)})
	c.baseTileSize = loader.CalcBaseSize(images)

	c.Layers = []attributes.Layer{}
	c.Layers = append(c.Layers, attributes.Layer{})
	c.Layers[0].Construct(160, 230, c.baseTileSize)
}

func (c *Canvas) ChangeDimensions(dimensions [2]int) {
	for layerIndex := range c.Layers {
		var temp attributes.Layer = make(attributes.Layer, dimensions[1])

		for row := range dimensions[1] {
			temp[row] = make([]attributes.Tile, dimensions[0])
			for col := range dimensions[0] {
				// keep old
				if row < len(c.Layers[0]) && col < len(c.Layers[0][0]) {
					temp[row][col] = c.Layers[layerIndex][row][col]
					continue
				}
				// add new tiles when enlarging canvas
				temp[row][col] = attributes.Tile{}
				temp[row][col].Construct(
					attributes.Spatial{X: float64(col) * c.baseTileSize, Y: float64(row) * c.baseTileSize},
					c.baseTileSize,
					"air",
					nil,
				)
			}
		}
		c.Layers[layerIndex] = make(attributes.Layer, dimensions[1])
		for row := range dimensions[1] {
			c.Layers[layerIndex][row] = make([]attributes.Tile, dimensions[0])
			for col := range dimensions[0] {
				c.Layers[layerIndex][row][col] = temp[row][col]
			}
		}
	}
	c.CheckBoundsAfterDimensionChange()
}

func (c *Canvas) CheckBoundsAfterDimensionChange() {
	c.camera.CheckBoundsAfterDimensionChange(
		attributes.Spatial{
			X: float64(len(c.Layers[0][0])) * c.baseTileSize * attributes.SCALE,
			Y: float64(len(c.Layers[0])) * c.baseTileSize * attributes.SCALE,
		},
	)
}

func (c *Canvas) SetActiveLayer(newActive int) {
	c.activeLayer = newActive
}

func (c *Canvas) ActiveLayer() int {
	return c.activeLayer
}

func (c *Canvas) Update(currentImageName string, currentImage *ebiten.Image, cursor attributes.Rect) {
	c.camera.Move(attributes.Spatial{X: float64(len(c.Layers[0][0])) * c.baseTileSize * attributes.SCALE, Y: float64(len(c.Layers[0])) * c.baseTileSize * attributes.SCALE})
	var screenWidth, screenHeight int = ebiten.Monitor().Size()

	cursor.Position.X += c.camera.Offset().X - float64(screenWidth)*0.25
	cursor.Position.Y += c.camera.Offset().Y

	var startRow int = int(math.Round(c.camera.offset.Y/(attributes.SCALE*c.baseTileSize))) - 2
	var startCol int = int(math.Round(c.camera.offset.X/(attributes.SCALE*c.baseTileSize))) - 2
	var endRow int = int(math.Round((c.camera.offset.Y+float64(screenHeight))/(attributes.SCALE*c.baseTileSize))) + 2
	var endCol int = int(math.Round((c.camera.offset.X+float64(screenWidth))/(attributes.SCALE*c.baseTileSize))) + 2

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row > -1 && row < len(c.Layers[0]) && col > -1 && col < len(c.Layers[0][0]) {
				if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
					if c.Layers[c.activeLayer][row][col].Rect().CollidePoint(cursor.Center()) && cursor.Right()-c.camera.Offset().X+float64(screenWidth)*0.25 > float64(screenWidth)*0.25 {
						if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
							c.Layers[c.activeLayer][row][col].SetImage(currentImageName, currentImage)
						}
						if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
							c.Layers[c.activeLayer][row][col].SetImage("air", nil)
						}
					}
					continue
				}
				if cursor.CollideRect(c.Layers[c.activeLayer][row][col].Rect()) && cursor.Right()-c.camera.Offset().X+float64(screenWidth)*0.25 > float64(screenWidth)*0.25 {
					if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
						c.Layers[c.activeLayer][row][col].SetImage(currentImageName, currentImage)
					}
					if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
						c.Layers[c.activeLayer][row][col].SetImage("air", nil)
					}
				}
			}
		}
	}
}

func (c *Canvas) AddLayer() {
	var newLayer attributes.Layer = attributes.Layer{}
	newLayer.Construct(len(c.Layers[0]), len(c.Layers[0][0]), c.baseTileSize)
	c.Layers = append(c.Layers, newLayer)
}

func (c *Canvas) drawLayers() {
	var screenWidth, screenHeight int = ebiten.Monitor().Size()

	var startRow int = int(math.Round(c.camera.offset.Y/(attributes.SCALE*c.baseTileSize))) - 2
	var startCol int = int(math.Round(c.camera.offset.X/(attributes.SCALE*c.baseTileSize))) - 2
	var endRow int = int(math.Round((c.camera.offset.Y+float64(screenHeight))/(attributes.SCALE*c.baseTileSize))) + 2
	var endCol int = int(math.Round((c.camera.offset.X+float64(screenWidth))/(attributes.SCALE*c.baseTileSize))) + 2

	for layerIndex := len(c.Layers) - 1; layerIndex >= 0; layerIndex-- {
		if layerIndex == c.activeLayer {
			continue
		}
		for row := startRow; row <= endRow; row++ {
			for col := startCol; col <= endCol; col++ {
				if row > -1 && row < len(c.Layers[0]) && col > -1 && col < len(c.Layers[0][0]) {
					if c.Layers[layerIndex][row][col].Image() != nil && c.camera.Rect().CollideRect(c.Layers[layerIndex][row][col].Rect()) {
						var opacity float64 = 1 / (float64(layerIndex) + 0.2) // non-active layers are going to have at least by 0.2 lower opacity than the active layer
						c.Layers[layerIndex][row][col].Draw(c.canvas, opacity, c.camera.offset)
					}
				}
			}
		}
	}
	// draw active layer last
	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row > -1 && row < len(c.Layers[0]) && col > -1 && col < len(c.Layers[0][0]) {
				if c.Layers[c.activeLayer][row][col].Image() != nil && c.camera.Rect().CollideRect(c.Layers[c.activeLayer][row][col].Rect()) {
					c.Layers[c.activeLayer][row][col].Draw(c.canvas, 1, c.camera.offset)
				}
			}
		}
	}
}

func (c *Canvas) Draw(surface *ebiten.Image) {
	c.canvas.Fill(color.RGBA{15, 15, 15, 255})

	c.drawLayers()

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(c.position.X, c.position.Y)
	surface.DrawImage(c.canvas, options)
}
