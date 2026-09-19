package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/menu/menu_background.png
var menuImageData []byte

const (
	GameStateMenu = iota
	GameStatePlaying
	GameStateOptions
)

const (
	MenuActionNone = iota
	MenuActionPlay
	MenuActionOptions
	MenuActionQuit
)

const (
	menuWidth  = 320
	menuHeight = 180
)

type MainMenu struct {
	Selected int
	Items    []string
	Image    *ebiten.Image
	Timer    int
}

func loadMenuImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func NewMainMenu() *MainMenu {
	return &MainMenu{
		Selected: 0,
		Items: []string{
			"JOUER",
			"OPTIONS",
			"QUITTER",
		},
		Image: loadMenuImage(menuImageData),
		Timer: 0,
	}
}

func (m *MainMenu) Update() int {
	m.Timer++

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.Selected--

		if m.Selected < 0 {
			m.Selected = len(m.Items) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.Selected++

		if m.Selected >= len(m.Items) {
			m.Selected = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.Selected {
		case 0:
			return MenuActionPlay
		case 1:
			return MenuActionOptions
		case 2:
			return MenuActionQuit
		}
	}

	return MenuActionNone
}

func drawMenuBackgroundCover(screen *ebiten.Image, background *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 8, B: 18, A: 255})

	if background == nil {
		return
	}

	imageWidth := float64(background.Bounds().Dx())
	imageHeight := float64(background.Bounds().Dy())

	scaleX := float64(menuWidth) / imageWidth
	scaleY := float64(menuHeight) / imageHeight

	scale := scaleX

	if scaleY > scale {
		scale = scaleY
	}

	drawWidth := imageWidth * scale
	drawHeight := imageHeight * scale

	drawX := (float64(menuWidth) - drawWidth) / 2
	drawY := (float64(menuHeight) - drawHeight) / 2

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(drawX, drawY)
	options.Filter = ebiten.FilterLinear

	screen.DrawImage(background, options)
}

func drawCenteredMenuText(screen *ebiten.Image, text string, y int) {
	textWidth := len(text) * 6
	x := menuWidth/2 - textWidth/2

	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func drawMenuButton(screen *ebiten.Image, text string, x float64, y float64, width float64, height float64, selected bool, timer int) {
	if selected {
		pulseValue := 180 + int(50*math.Sin(float64(timer)*0.08))

		if pulseValue < 0 {
			pulseValue = 0
		}

		if pulseValue > 255 {
			pulseValue = 255
		}

		pulse := uint8(pulseValue)

		ebitenutil.DrawRect(screen, x-3, y-3, width+6, height+6, color.RGBA{R: 20, G: 100, B: 255, A: 60})
		ebitenutil.DrawRect(screen, x-2, y-2, width+4, height+4, color.RGBA{R: 50, G: 170, B: 255, A: pulse})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 25, B: 55, A: 230})

		ebitenutil.DrawRect(screen, x+2, y+2, width-4, 1, color.RGBA{R: 120, G: 220, B: 255, A: 255})
		ebitenutil.DrawRect(screen, x+2, y+height-3, width-4, 1, color.RGBA{R: 30, G: 110, B: 220, A: 255})

		ebitenutil.DebugPrintAt(screen, ">", int(x)+12, int(y)+7)
		ebitenutil.DebugPrintAt(screen, "<", int(x+width)-17, int(y)+7)
	} else {
		ebitenutil.DrawRect(screen, x-1, y-1, width+2, height+2, color.RGBA{R: 180, G: 140, B: 65, A: 190})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 12, B: 25, A: 220})
	}

	textWidth := len(text) * 6
	textX := int(x) + int(width)/2 - textWidth/2
	textY := int(y) + 7

	ebitenutil.DebugPrintAt(screen, text, textX, textY)
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	drawMenuBackgroundCover(screen, m.Image)

	ebitenutil.DrawRect(screen, 0, 0, menuWidth, menuHeight, color.RGBA{R: 0, G: 5, B: 15, A: 40})

	panelX := 78.0
	panelY := 7.0
	panelWidth := 164.0
	panelHeight := 160.0

	ebitenutil.DrawRect(screen, panelX, panelY, panelWidth, panelHeight, color.RGBA{R: 2, G: 8, B: 20, A: 125})

	ebitenutil.DrawRect(screen, 87, 15, 46, 1, color.RGBA{R: 220, G: 175, B: 70, A: 255})
	ebitenutil.DrawRect(screen, 187, 15, 46, 1, color.RGBA{R: 220, G: 175, B: 70, A: 255})

	drawCenteredMenuText(screen, "E L D O R I A", 19)

	ebitenutil.DrawRect(screen, 105, 31, 110, 1, color.RGBA{R: 220, G: 175, B: 70, A: 220})

	drawCenteredMenuText(screen, "LE ROYAUME OUBLIE", 37)

	buttonX := 88.0
	buttonWidth := 144.0
	buttonHeight := 22.0

	drawMenuButton(screen, "JOUER", buttonX, 54, buttonWidth, buttonHeight, m.Selected == 0, m.Timer)
	drawMenuButton(screen, "OPTIONS", buttonX, 84, buttonWidth, buttonHeight, m.Selected == 1, m.Timer)
	drawMenuButton(screen, "QUITTER", buttonX, 114, buttonWidth, buttonHeight, m.Selected == 2, m.Timer)

	drawCenteredMenuText(screen, "HAUT / BAS : CHOISIR", 145)
	drawCenteredMenuText(screen, "ENTREE : VALIDER", 157)
}

func DrawOptionsScreen(screen *ebiten.Image, background *ebiten.Image) {
	drawMenuBackgroundCover(screen, background)

	ebitenutil.DrawRect(screen, 0, 0, menuWidth, menuHeight, color.RGBA{R: 0, G: 5, B: 15, A: 80})

	panelX := 48.0
	panelY := 12.0
	panelWidth := 224.0
	panelHeight := 150.0

	ebitenutil.DrawRect(screen, panelX, panelY, panelWidth, panelHeight, color.RGBA{R: 3, G: 10, B: 25, A: 225})

	ebitenutil.DrawRect(screen, 47, 11, 226, 2, color.RGBA{R: 210, G: 165, B: 70, A: 255})
	ebitenutil.DrawRect(screen, 47, 162, 226, 2, color.RGBA{R: 210, G: 165, B: 70, A: 255})

	drawCenteredMenuText(screen, "OPTIONS", 24)

	ebitenutil.DrawRect(screen, 105, 38, 110, 1, color.RGBA{R: 60, G: 165, B: 255, A: 255})

	ebitenutil.DebugPrintAt(screen, "DEPLACEMENT", 65, 53)
	ebitenutil.DebugPrintAt(screen, "Q / D", 210, 53)

	ebitenutil.DebugPrintAt(screen, "SAUT", 65, 70)
	ebitenutil.DebugPrintAt(screen, "ESPACE", 210, 70)

	ebitenutil.DebugPrintAt(screen, "ATTAQUE", 65, 87)
	ebitenutil.DebugPrintAt(screen, "CLIC GAUCHE", 185, 87)

	ebitenutil.DebugPrintAt(screen, "PORTAIL", 65, 104)
	ebitenutil.DebugPrintAt(screen, "E", 210, 104)

	ebitenutil.DebugPrintAt(screen, "MENU", 65, 121)
	ebitenutil.DebugPrintAt(screen, "ECHAP", 210, 121)

	drawCenteredMenuText(screen, "ENTREE OU ECHAP : RETOUR", 145)
}
