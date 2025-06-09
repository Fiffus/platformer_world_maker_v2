package guicomponents

import (
	"bytes"
	"fmt"
	"image/color"
	"io/fs"
	"log"
	"os"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/fonts"
	"platformer_world_maker_v2/loader"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	rect     attributes.Rect
	clr      attributes.Color
	text     string
	fontFace *text.GoTextFace
}

func (b *Button) Construct(position attributes.Spatial, buttonText string) {
	b.rect = attributes.Rect{
		Position: position,
	}
	b.text = buttonText
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	b.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   40,
	}
	b.clr = attributes.Color{
		Current:   color.RGBA{58, 55, 94, 255},
		Normal:    color.RGBA{58, 55, 94, 255},
		Highlight: color.RGBA{77, 73, 122, 255},
	}
}

func (b *Button) HighLight() {
	x, y := ebiten.CursorPosition()
	if b.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)}) {
		b.clr.Current = b.clr.Highlight
		return
	}
	b.clr.Current = b.clr.Normal
}

func (b *Button) Rect() attributes.Rect {
	return b.rect
}

func (b *Button) Pressed() bool {
	var x, y int = ebiten.CursorPosition()
	if b.rect.CollidePoint(attributes.Spatial{X: float64(x), Y: float64(y)}) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	return false
}

func (b *Button) Draw(surface *ebiten.Image) {
	b.rect.Draw(surface, b.clr.Current, attributes.Spatial{X: 0, Y: 0})

	options := &text.DrawOptions{}
	options.GeoM.Translate(b.rect.Left()+b.rect.Size.X/2-float64(len(b.text))*150/18, b.rect.Top()+b.rect.Size.Y/2+10)
	options.ColorScale.Scale(1, 1, 1, 1)
	text.Draw(surface, b.text, b.fontFace, options)
}

func (b *Button) Fill(layer *attributes.Layer, currentImageName string, currentImage *ebiten.Image) {
	for row := range *layer {
		for col := range (*layer)[0] {
			(*layer)[row][col].SetImage(currentImageName, currentImage)
		}
	}
}

func (b *Button) Load(projectName string, layers *[]attributes.Layer, images map[string]*ebiten.Image, dimensionChanger *DimensionChanger) {
	if projectName == "" {
		return
	}
	var dimensions [2]int = loader.LoadDimensions(fmt.Sprintf("%s/properties.txt", projectName))

	for i := range *layers {
		(*layers)[i] = make(attributes.Layer, dimensions[1])
		for j := range dimensions[1] {
			(*layers)[i][j] = make([]attributes.Tile, dimensions[0])
		}
		(*layers)[i] = loader.GenerateLevelTiles(images, fmt.Sprintf("%s/layer_%d.csv", projectName, i), dimensions)
	}

	dimensionChanger.entry.text = fmt.Sprintf("%dx%d", dimensions[0], dimensions[1])
}

func (b *Button) Save(projectName string, layers []attributes.Layer) {
	if projectName == "" {
		return
	}
	var err error
	var worlds []fs.DirEntry
	if worlds, err = os.ReadDir("worlds"); err != nil {
		log.Fatal(err)
	}
	for _, folder := range worlds {
		if projectName == folder.Name() {
			return
		}
	}
	os.Mkdir("worlds/"+projectName, os.ModePerm)

	for i := range layers {
		var layerFile *os.File
		if layerFile, err = os.OpenFile(fmt.Sprintf("worlds/%s/layer_1.csv", projectName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		}
		defer layerFile.Close()
		for row := range layers[i] {
			for col := range layers[i][row] {
				if row == len(layers[i])-1 && col == len(layers[i][row])-1 {
					layerFile.WriteString(layers[i][row][col].ImageName())
				} else {
					layerFile.WriteString(layers[i][row][col].ImageName() + ";")
				}
			}
		}
	}

	var propertiesFile *os.File
	if propertiesFile, err = os.OpenFile(fmt.Sprintf("worlds/%s/properties.txt", projectName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	}
	defer propertiesFile.Close()

	propertiesFile.WriteString(fmt.Sprintf("%dx%d", len(layers[0][0]), len(layers[0])))
}

func (b *Button) Polish(layer *attributes.Layer, images map[string]*ebiten.Image, tileType string) {
	// rules that set where in the current layer in relation to tile's {row, col} has to be nil value for it to be replaced by specific tile
	rules := map[string][][2]int{
		"tile1":  {{0, -1}, {-1, 0}},
		"tile2":  {{-1, 0}},
		"tile3":  {{-1, 0}, {0, 1}},
		"tile4":  {{0, -1}},
		"tile6":  {{0, 1}},
		"tile7":  {{0, -1}, {1, 0}},
		"tile8":  {{1, 0}},
		"tile9":  {{1, 0}, {0, 1}},
		"tile10": {{-1, -1}},
		"tile11": {{-1, 1}},
		"tile12": {{1, -1}},
		"tile13": {{1, 1}},
	}
	var tileTypes []string = []string{"tile10", "tile11", "tile12", "tile13", "tile2", "tile4", "tile6", "tile8", "tile1", "tile3", "tile7", "tile9"}
	for row := range *layer {
		for col := range (*layer)[0] {
			var followedRules map[string]int = make(map[string]int)
			if (*layer)[row][col].ImageName() == tileType+"_tile5" {
				for tile := range rules {
					for _, rule := range rules[tile] {
						if row+rule[0] > -1 && row+rule[0] < len(*layer) && col+rule[1] > -1 && col+rule[1] < len((*layer)[0]) {
							if strings.Split((*layer)[row+rule[0]][col+rule[1]].ImageName(), "_")[0] != tileType {
								followedRules[tile]++
							}
						}
					}
				}
				var useTile string = tileType + "_tile5"
				for _, tile := range tileTypes {
					if images[tileType+"_"+tile] != nil {
						if followedRules[tile] == len(rules[tile]) {
							useTile = tileType + "_" + tile
						}
					}
				}
				if useTile != "tile5" {
					(*layer)[row][col].SetImage(useTile, images[useTile])
					continue
				}
				if strings.Split((*layer)[row][col].ImageName(), "_")[0] != tileType {
					(*layer)[row][col].SetImage(useTile, images[useTile])
				}
			}
		}
	}
}
