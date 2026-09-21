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
var menuBackgroundData []byte

const (
	GameStateMenu = iota
	GameStatePlaying
	GameStateOptions
)

const (
	MenuActionNone = iota
	MenuActionContinue
	MenuActionPlay
	MenuActionDungeon
	MenuActionOptions
	MenuActionQuit
)

type MainMenu struct {
	Selected int
	Items    []string
	Image    *ebiten.Image
	Timer    int

	ForgeOpen bool
	Forge     *ForgeMenu
}

func NewMainMenu() *MainMenu {
	img, _, err := image.Decode(bytes.NewReader(menuBackgroundData))
	if err != nil {
		panic(err)
	}

	selected := 0

	if !HasSave() {
		selected = 1
	}

	return &MainMenu{
		Selected: selected,
		Items: []string{
			"CONTINUER",
			"NOUVELLE PARTIE",
			"DONJON",
			"FORGE",
			"OPTIONS",
			"QUITTER",
		},
		Image:     ebiten.NewImageFromImage(img),
		ForgeOpen: false,
		Forge:     NewForgeMenu(),
	}
}

func (m *MainMenu) moveSelection(direction int) {
	for i := 0; i < len(m.Items); i++ {
		m.Selected += direction

		if m.Selected < 0 {
			m.Selected = len(m.Items) - 1
		}

		if m.Selected >= len(m.Items) {
			m.Selected = 0
		}

		if m.Selected == 0 && !HasSave() {
			continue
		}

		break
	}
}

func (m *MainMenu) Update() int {
	m.Timer++

	if m.ForgeOpen {
		if m.Forge == nil {
			m.Forge = NewForgeMenu()
		}

		if m.Forge.Update() {
			m.ForgeOpen = false
		}

		return MenuActionNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.moveSelection(-1)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.moveSelection(1)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch m.Selected {
		case 0:
			if HasSave() {
				return MenuActionContinue
			}

		case 1:
			return MenuActionPlay

		case 2:
			return MenuActionDungeon

		case 3:
			if m.Forge == nil {
				m.Forge = NewForgeMenu()
			}

			m.Forge.Reload()
			m.ForgeOpen = true

		case 4:
			return MenuActionOptions

		case 5:
			return MenuActionQuit
		}
	}

	return MenuActionNone
}

func drawMenuBackgroundCover(screen *ebiten.Image, background *ebiten.Image) {
	if background == nil {
		screen.Fill(
			color.RGBA{
				R: 5,
				G: 7,
				B: 12,
				A: 255,
			},
		)

		return
	}

	imageWidth := float64(background.Bounds().Dx())
	imageHeight := float64(background.Bounds().Dy())

	screenWidth := 320.0
	screenHeight := 180.0

	scaleX := screenWidth / imageWidth
	scaleY := screenHeight / imageHeight

	scale := math.Max(scaleX, scaleY)

	drawWidth := imageWidth * scale
	drawHeight := imageHeight * scale

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(
		scale,
		scale,
	)

	options.GeoM.Translate(
		(screenWidth-drawWidth)/2,
		(screenHeight-drawHeight)/2,
	)

	options.Filter = ebiten.FilterLinear

	screen.DrawImage(
		background,
		options,
	)
}

func drawCenteredMenuText(screen *ebiten.Image, text string, y int) {
	width := len(text) * 6
	x := (320 - width) / 2

	ebitenutil.DebugPrintAt(
		screen,
		text,
		x,
		y,
	)
}

func drawMenuButton(screen *ebiten.Image, text string, x float64, y float64, width float64, height float64, selected bool) {
	backgroundColor := color.RGBA{
		R: 12,
		G: 17,
		B: 28,
		A: 220,
	}

	borderColor := color.RGBA{
		R: 70,
		G: 85,
		B: 105,
		A: 255,
	}

	if selected {
		backgroundColor = color.RGBA{
			R: 20,
			G: 55,
			B: 95,
			A: 235,
		}

		borderColor = color.RGBA{
			R: 80,
			G: 180,
			B: 255,
			A: 255,
		}
	}

	ebitenutil.DrawRect(
		screen,
		x-1,
		y-1,
		width+2,
		height+2,
		borderColor,
	)

	ebitenutil.DrawRect(
		screen,
		x,
		y,
		width,
		height,
		backgroundColor,
	)

	textWidth := len(text) * 6

	textX := int(x) + (int(width)-textWidth)/2
	textY := int(y) + int(height)/2 - 4

	ebitenutil.DebugPrintAt(
		screen,
		text,
		textX,
		textY,
	)
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	drawMenuBackgroundCover(
		screen,
		m.Image,
	)

	if m.ForgeOpen {
		if m.Forge != nil {
			m.Forge.Draw(screen)
		}

		return
	}

	panelColor := color.RGBA{
		R: 4,
		G: 8,
		B: 15,
		A: 205,
	}

	borderColor := color.RGBA{
		R: 65,
		G: 95,
		B: 125,
		A: 255,
	}

	ebitenutil.DrawRect(
		screen,
		77,
		4,
		166,
		172,
		borderColor,
	)

	ebitenutil.DrawRect(
		screen,
		78,
		5,
		164,
		170,
		panelColor,
	)

	drawCenteredMenuText(
		screen,
		"E L D O R I A",
		10,
	)

	drawCenteredMenuText(
		screen,
		"LE ROYAUME OUBLIE",
		22,
	)

	buttonX := 91.0
	buttonWidth := 138.0
	buttonHeight := 16.0

	positions := []float64{
		36,
		56,
		76,
		96,
		116,
		136,
	}

	for i, item := range m.Items {
		selected := m.Selected == i

		if i == 0 && !HasSave() {
			ebitenutil.DrawRect(
				screen,
				buttonX-1,
				positions[i]-1,
				buttonWidth+2,
				buttonHeight+2,
				color.RGBA{
					R: 45,
					G: 45,
					B: 50,
					A: 255,
				},
			)

			ebitenutil.DrawRect(
				screen,
				buttonX,
				positions[i],
				buttonWidth,
				buttonHeight,
				color.RGBA{
					R: 15,
					G: 15,
					B: 18,
					A: 210,
				},
			)

			textWidth := len(item) * 6
			textX := int(buttonX) + (int(buttonWidth)-textWidth)/2

			ebitenutil.DebugPrintAt(
				screen,
				item,
				textX,
				int(positions[i])+4,
			)

			continue
		}

		drawMenuButton(
			screen,
			item,
			buttonX,
			positions[i],
			buttonWidth,
			buttonHeight,
			selected,
		)
	}

	drawCenteredMenuText(
		screen,
		"FLECHES + ENTREE",
		159,
	)

	if !HasSave() {
		drawCenteredMenuText(
			screen,
			"Aucune sauvegarde",
			169,
		)
	} else {
		drawCenteredMenuText(
			screen,
			"Sauvegarde disponible",
			169,
		)
	}
}
