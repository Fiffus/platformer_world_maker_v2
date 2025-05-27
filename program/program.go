package program

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Program struct{}

func (p *Program) Init() {
	ebiten.SetMaxTPS(60)
	ebiten.SetWindowTitle("Platformer World Maker")
	ebiten.SetFullscreen(true)
}

func (p *Program) Update() error {
	return nil
}

func (p *Program) Draw(screen *ebiten.Image) {

}

func (p *Program) Run() {
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
}

func (p *Program) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
