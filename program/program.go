package program

import (
	"log"
	"os"
	"platformer_world_maker_v2/attributes"
	"platformer_world_maker_v2/canvas"
	"platformer_world_maker_v2/guicomponents"
	"platformer_world_maker_v2/loader"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Program struct {
	images           map[string]*ebiten.Image
	canv             canvas.Canvas
	dimensionChanger guicomponents.DimensionChanger
	toolbar          guicomponents.Toolbar
	buttons          map[string]*guicomponents.Button
	dropDown         guicomponents.DropDown
	puwLoad          guicomponents.PopUpWindow
	puwSave          guicomponents.PopUpWindow
	cursor           canvas.Cursor
}

func (p *Program) Init() {
	ebiten.SetWindowTitle("Platformer World Maker")
	ebiten.SetFullscreen(true)
	p.images = loader.LoadImages()

	p.canv.Contruct(p.images)
	p.dimensionChanger.Construct(attributes.Spatial{X: 30, Y: 600})
	p.toolbar.Construct(p.images)

	p.buttons = make(map[string]*guicomponents.Button)
	p.buttons["exit"] = &guicomponents.Button{}
	p.buttons["exit"].Construct(attributes.Spatial{X: 30, Y: 990}, "Exit")
	p.buttons["load"] = &guicomponents.Button{}
	p.buttons["load"].Construct(attributes.Spatial{X: 30, Y: 900}, "Load")
	p.buttons["save"] = &guicomponents.Button{}
	p.buttons["save"].Construct(attributes.Spatial{X: 260, Y: 900}, "Save")
	p.buttons["polish"] = &guicomponents.Button{}
	p.buttons["polish"].Construct(attributes.Spatial{X: 30, Y: 810}, "Polish")
	p.buttons["fill"] = &guicomponents.Button{}
	p.buttons["fill"].Construct(attributes.Spatial{X: 260, Y: 810}, "Fill")

	p.dropDown = guicomponents.DropDown{}
	p.dropDown.Construct(attributes.Spatial{X: 260, Y: 600})

	p.puwLoad = guicomponents.PopUpWindow{}
	p.puwLoad.Construct("Enter project name:", "./worlds/")
	p.puwSave = guicomponents.PopUpWindow{}
	p.puwSave.Construct("Save project as:", "")

	p.cursor.Construct(attributes.Spatial{X: 70, Y: 70})
}

func (p *Program) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if p.puwLoad.Active() {
		p.puwLoad.Update()
		if p.puwLoad.Confirmed() {
			p.buttons["load"].Load(p.puwLoad.Value(), &p.canv.Layers, p.images, &p.dropDown, &p.dimensionChanger)
		}
		return nil
	}

	if p.puwSave.Active() {
		p.puwSave.Update()
		if p.puwSave.Confirmed() {
			p.buttons["save"].Save(p.puwSave.Value(), p.canv.Layers)
		}
		return nil
	}

	p.dropDown.Update()
	p.cursor.Update()

	if p.dropDown.JustAddedNewRow() {
		p.canv.AddLayer()
	}

	if p.dropDown.Active() {
		p.canv.SetActiveLayer(p.dropDown.ActiveRow())
		return nil
	}

	p.canv.Update(p.toolbar.SelectedName(), p.toolbar.SelectedImage(), p.cursor.Rect())
	p.dimensionChanger.Update()
	if p.dimensionChanger.ValidX() && p.dimensionChanger.ValidY() && (p.dimensionChanger.ChangeableX() || p.dimensionChanger.ChangeableY()) {
		var data [2]int
		data[0], _ = strconv.Atoi(p.dimensionChanger.ValueX())
		data[1], _ = strconv.Atoi(p.dimensionChanger.ValueY())
		p.canv.ChangeDimensions(data)
		p.dimensionChanger.ChangeUsed()
	}
	p.toolbar.Update()

	for key := range p.buttons {
		if p.buttons[key].Pressed() {
			switch key {
			case "exit":
				os.Exit(0)
			case "load":
				p.puwLoad.Activate()
			case "polish":
				p.buttons[key].Polish(&p.canv.Layers[p.canv.ActiveLayer()], p.images, strings.Split(p.toolbar.SelectedName(), "_")[0])
			case "save":
				p.puwSave.Activate()
			case "fill":
				p.buttons[key].Fill(&p.canv.Layers[p.canv.ActiveLayer()], p.toolbar.SelectedName(), p.toolbar.SelectedImage())
			}
		}
		p.buttons[key].HighLight()
	}

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

	if p.puwLoad.Active() {
		p.puwLoad.Draw(screen)
	}

	if p.puwSave.Active() {
		p.puwSave.Draw(screen)
	}
	p.dropDown.Draw(screen)
}

func (p *Program) Run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}

func (p *Program) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
