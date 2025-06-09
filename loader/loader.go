package loader

import (
	"fmt"
	"log"
	"os"
	"platformer_world_maker_v2/attributes"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func loadImageNamesFromCSV(projectLayerName string) []string {
	var imageNames []string
	fileContent, err := os.ReadFile(fmt.Sprintf("worlds/%s", projectLayerName))
	if err != nil {
		log.Fatal(err)
	}
	imageNames = append(imageNames, strings.Split(string(fileContent), ";")...)
	return imageNames
}

func transformData(data []string, dimensions [2]int) [][]string {
	var transformedData [][]string = make([][]string, dimensions[1])
	var imgIndex int = 0
	for i := 0; i < dimensions[1]; i++ {
		transformedData[i] = make([]string, dimensions[0])
		for j := 0; j < dimensions[0]; j++ {
			transformedData[i][j] = data[imgIndex]
			imgIndex++
		}
	}
	return transformedData
}

func getImageNames() []string {
	dir, err := os.ReadDir("textures")
	var imageNames []string
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range dir {
		imageNames = append(imageNames, image.Name())
	}
	return imageNames
}

func LoadTextures() map[string]*ebiten.Image {
	var tileImages map[string]*ebiten.Image = make(map[string]*ebiten.Image)
	for _, imageName := range getImageNames() {
		image, _, err := ebitenutil.NewImageFromFile("textures/" + imageName)
		tileImages[strings.Split(imageName, ".")[0]] = image
		if err != nil {
			log.Fatal(err)
		}
	}
	return tileImages
}

func CalcBaseSize(loadedImages map[string]*ebiten.Image) float64 {
	for key := range loadedImages {
		var splitTileName []string = strings.Split(key, "_")
		if len(splitTileName) > 1 && splitTileName[1] == "tile5" {
			return float64(loadedImages[key].Bounds().Max.X) * attributes.SCALE
		}
	}
	return 40
}

func GenerateLevelTiles(loadedImages map[string]*ebiten.Image, projectLayerName string, dimensions [2]int) attributes.Layer {
	var baseSize float64 = CalcBaseSize(loadedImages)
	var imageNames [][]string = transformData(loadImageNamesFromCSV(fmt.Sprintf("worlds/%s", projectLayerName)), dimensions)
	var layer attributes.Layer = make(attributes.Layer, dimensions[1])

	for row := range imageNames {
		layer[row] = make([]attributes.Tile, dimensions[0])
		for col := range imageNames[row] {
			if imageNames[row][col] != "air" {
				layer[row][col] = attributes.Tile{}
				layer[row][col].Construct(
					attributes.Spatial{X: float64(col) * baseSize, Y: float64(row) * baseSize},
					baseSize,
					imageNames[row][col],
					loadedImages[imageNames[row][col]],
				)
				continue
			}
			layer[row][col] = attributes.Tile{}
			layer[row][col].Construct(
				attributes.Spatial{X: float64(col) * baseSize, Y: float64(row) * baseSize},
				baseSize,
				"air",
				nil,
			)
		}
	}
	return layer
}

func LoadDimensions(projectName string) [2]int {
	var data [2]int
	fileContent, err := os.ReadFile(fmt.Sprintf("worlds/%s/properties.txt", projectName))
	if err != nil {
		log.Fatal(err)
	}
	var rawData []string = strings.Split(string(fileContent), "x")
	data[0], _ = strconv.Atoi(rawData[0])
	data[1], _ = strconv.Atoi(rawData[1])
	return data
}
