package program

import (
	"log"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/canvas"
	"platformer_world_maker_v2/guicomponents"
	"platformer_world_maker_v2/loader"

	"github.com/hajimehoshi/ebiten/v2"
)

type Program struct {
	images           map[string]*ebiten.Image
	canv             canvas.Canvas
	dimensionChanger guicomponents.DimensionChanger
	toolbar          guicomponents.Toolbar
	buttons          map[string]*guicomponents.Button
	cursor           canvas.Cursor
}

func (p *Program) Init() {
	ebiten.SetWindowTitle("Platformer World Maker")
	ebiten.SetFullscreen(true)
	p.images = loader.LoadImages()

	p.canv.Contruct(p.images)
	p.dimensionChanger.Construct(attributes.Spatial{X: 30, Y: 500})
	p.toolbar.Construct(p.images)

	p.buttons = make(map[string]*guicomponents.Button)
	p.buttons["exit"].Construct(attributes.Spatial{X: 30, Y: 600}, "Exit")
	p.buttons["load"].Construct(attributes.Spatial{X: 30, Y: 600}, "Load")
	p.buttons["polish"].Construct(attributes.Spatial{X: 30, Y: 600}, "Polish")
	p.buttons["save"].Construct(attributes.Spatial{X: 30, Y: 600}, "Save")

	p.cursor.Construct(attributes.Spatial{X: 70, Y: 70})
}

func (p *Program) Update() error {
	return nil
}

func (p *Program) Draw(screen *ebiten.Image) {
	p.canv.Draw(screen)
	p.dimensionChanger.Draw(screen)
	p.toolbar.Draw(screen)

	for key := range p.buttons {
		p.buttons[key].Draw(screen)
	}

	p.cursor.Draw(screen)
}

func (p *Program) Run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}

func (p *Program) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
