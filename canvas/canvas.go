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
	layers       []attributes.Layer
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

	c.layers = []attributes.Layer{}
	c.layers = append(c.layers, attributes.Layer{})
	c.layers[0].Construct(160, 230, c.baseTileSize)
}

func (c *Canvas) ChangeDimensions(dimensions [2]int) {
	for layerIndex := range c.layers {
		var temp attributes.Layer = make(attributes.Layer, dimensions[1])

		for row := range dimensions[1] {
			temp[row] = make([]attributes.Tile, dimensions[0])
			for col := range dimensions[0] {
				// keep old
				if row < dimensions[1] && col < dimensions[0] {
					temp[row][col] = c.layers[layerIndex][row][col]
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
		c.layers[layerIndex] = make(attributes.Layer, dimensions[1])
		for row := range dimensions[1] {
			c.layers[layerIndex] = make(attributes.Layer, dimensions[0])
			for col := range dimensions[0] {
				c.layers[layerIndex][row][col] = temp[row][col]
			}
		}
	}
	c.CheckBoundsAfterDimensionChange()
}

func (c *Canvas) CheckBoundsAfterDimensionChange() {
	c.camera.CheckBoundsAfterDimensionChange(
		attributes.Spatial{
			X: float64(len(c.layers[0][0])) * c.baseTileSize * attributes.SCALE,
			Y: float64(len(c.layers[0])) * c.baseTileSize * attributes.SCALE,
		},
	)
}

func (c *Canvas) SetActiveLayer(newActive int) {
	c.activeLayer = newActive
}

func (c *Canvas) Update(currentImageName string, currentImage *ebiten.Image, cursor attributes.Rect) {
	c.camera.Move(attributes.Spatial{X: float64(len(c.layers[0][0])) * c.baseTileSize * attributes.SCALE, Y: float64(len(c.layers[0])) * c.baseTileSize * attributes.SCALE})
	var screenWidth, screenHeight int = ebiten.Monitor().Size()

	cursor.Position.X += c.camera.Offset().X - float64(screenWidth)*0.25
	cursor.Position.Y += c.camera.Offset().Y

	var startRow int = int(math.Round(c.camera.offset.Y/(attributes.SCALE*c.baseTileSize))) - 2
	var startCol int = int(math.Round(c.camera.offset.X/(attributes.SCALE*c.baseTileSize))) - 2
	var endRow int = int(math.Round((c.camera.offset.Y+float64(screenHeight))/(attributes.SCALE*c.baseTileSize))) + 2
	var endCol int = int(math.Round((c.camera.offset.X+float64(screenWidth))/(attributes.SCALE*c.baseTileSize))) + 2

	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row > -1 && row < len(c.layers[0]) && col > -1 && col < len(c.layers[0][0]) {
				if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
					if c.layers[c.activeLayer][row][col].Rect().CollidePoint(cursor.Center()) && cursor.Right()-c.camera.Offset().X+float64(screenWidth)*0.25 > float64(screenWidth)*0.25 {
						if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
							c.layers[c.activeLayer][row][col].SetImage(currentImageName, currentImage)
						}
						if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
							c.layers[c.activeLayer][row][col].SetImage("air", nil)
						}
					}
					continue
				}
				if cursor.CollideRect(c.layers[c.activeLayer][row][col].Rect()) && cursor.Right()-c.camera.Offset().X+float64(screenWidth)*0.25 > float64(screenWidth)*0.25 {
					if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
						c.layers[c.activeLayer][row][col].SetImage(currentImageName, currentImage)
					}
					if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
						c.layers[c.activeLayer][row][col].SetImage("air", nil)
					}
				}
			}
		}
	}
}

func (c *Canvas) drawLayers() {
	var screenWidth, screenHeight int = ebiten.Monitor().Size()

	var startRow int = int(math.Round(c.camera.offset.Y/(attributes.SCALE*c.baseTileSize))) - 2
	var startCol int = int(math.Round(c.camera.offset.X/(attributes.SCALE*c.baseTileSize))) - 2
	var endRow int = int(math.Round((c.camera.offset.Y+float64(screenHeight))/(attributes.SCALE*c.baseTileSize))) + 2
	var endCol int = int(math.Round((c.camera.offset.X+float64(screenWidth))/(attributes.SCALE*c.baseTileSize))) + 2

	for layerIndex := len(c.layers) - 1; layerIndex >= 0; layerIndex-- {
		if layerIndex == c.activeLayer {
			continue
		}
		for row := startRow; row <= endRow; row++ {
			for col := startCol; col <= endCol; col++ {
				if row > -1 && row < len(c.layers[0]) && col > -1 && col < len(c.layers[0]) {
					if c.layers[layerIndex][row][col].Image() != nil && c.camera.Rect().CollideRect(c.layers[layerIndex][row][col].Rect()) {
						var opacity float64 = 1 / (float64(layerIndex) + 0.2) // non-active layers are going to have at least by 0.2 lower opacity than the active layer
						c.layers[layerIndex][row][col].Draw(c.canvas, opacity, c.camera.offset)
					}
				}
			}
		}
	}
	// draw active layer last
	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row > -1 && row < len(c.layers[0]) && col > -1 && col < len(c.layers[0]) {
				if c.layers[c.activeLayer][row][col].Image() != nil && c.camera.Rect().CollideRect(c.layers[c.activeLayer][row][col].Rect()) {
					c.layers[c.activeLayer][row][col].Draw(c.canvas, 1, c.camera.offset)
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
