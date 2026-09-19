package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	PauseActionNone = iota
	PauseActionResume
	PauseActionOptions
	PauseActionMenu
)

type PauseMenu struct {
	Selected int
	Items    []string
	Timer    int
}

func NewPauseMenu() *PauseMenu {
	return &PauseMenu{
		Selected: 0,
		Items: []string{
			"REPRENDRE",
			"OPTIONS",
			"RETOUR MENU",
		},
		Timer: 0,
	}
}

func (p *PauseMenu) Update() int {
	p.Timer++

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return PauseActionResume
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.Selected--

		if p.Selected < 0 {
			p.Selected = len(p.Items) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.Selected++

		if p.Selected >= len(p.Items) {
			p.Selected = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch p.Selected {
		case 0:
			return PauseActionResume

		case 1:
			return PauseActionOptions

		case 2:
			return PauseActionMenu
		}
	}

	return PauseActionNone
}

func drawPauseButton(screen *ebiten.Image, text string, y float64, selected bool, timer int) {
	x := 90.0
	width := 140.0
	height := 21.0

	if selected {
		pulseValue := 180 + int(55*math.Sin(float64(timer)*0.08))

		if pulseValue < 0 {
			pulseValue = 0
		}

		if pulseValue > 255 {
			pulseValue = 255
		}

		pulse := uint8(pulseValue)

		ebitenutil.DrawRect(screen, x-3, y-3, width+6, height+6, color.RGBA{R: 20, G: 100, B: 255, A: 70})
		ebitenutil.DrawRect(screen, x-2, y-2, width+4, height+4, color.RGBA{R: 60, G: 180, B: 255, A: pulse})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 25, B: 55, A: 245})

		ebitenutil.DrawRect(screen, x+2, y+2, width-4, 1, color.RGBA{R: 130, G: 230, B: 255, A: 255})
		ebitenutil.DrawRect(screen, x+2, y+height-3, width-4, 1, color.RGBA{R: 30, G: 110, B: 220, A: 255})

		ebitenutil.DebugPrintAt(screen, ">", int(x)+10, int(y)+7)
		ebitenutil.DebugPrintAt(screen, "<", int(x+width)-15, int(y)+7)
	} else {
		ebitenutil.DrawRect(screen, x-1, y-1, width+2, height+2, color.RGBA{R: 185, G: 145, B: 70, A: 190})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 12, B: 25, A: 235})
	}

	textWidth := len(text) * 6
	textX := int(x) + int(width)/2 - textWidth/2
	textY := int(y) + 7

	ebitenutil.DebugPrintAt(screen, text, textX, textY)
}

func (p *PauseMenu) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, 320, 180, color.RGBA{R: 0, G: 0, B: 0, A: 145})

	panelX := 70.0
	panelY := 18.0
	panelWidth := 180.0
	panelHeight := 145.0

	ebitenutil.DrawRect(screen, panelX-2, panelY-2, panelWidth+4, panelHeight+4, color.RGBA{R: 180, G: 140, B: 60, A: 220})
	ebitenutil.DrawRect(screen, panelX, panelY, panelWidth, panelHeight, color.RGBA{R: 3, G: 10, B: 25, A: 235})

	ebitenutil.DrawRect(screen, 95, 34, 45, 1, color.RGBA{R: 220, G: 175, B: 70, A: 255})
	ebitenutil.DrawRect(screen, 180, 34, 45, 1, color.RGBA{R: 220, G: 175, B: 70, A: 255})

	ebitenutil.DebugPrintAt(screen, "PAUSE", 145, 29)

	drawPauseButton(screen, "REPRENDRE", 52, p.Selected == 0, p.Timer)
	drawPauseButton(screen, "OPTIONS", 82, p.Selected == 1, p.Timer)
	drawPauseButton(screen, "RETOUR MENU", 112, p.Selected == 2, p.Timer)

	ebitenutil.DebugPrintAt(screen, "HAUT / BAS : CHOISIR", 99, 145)
	ebitenutil.DebugPrintAt(screen, "ENTREE : VALIDER", 111, 156)
}
