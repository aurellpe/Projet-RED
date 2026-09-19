package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	PuzzleRuneUp = iota
	PuzzleRuneRight
	PuzzleRuneDown
	PuzzleRuneLeft
)

const (
	PuzzleStateShowing = iota
	PuzzleStateInput
	PuzzleStateFailure
	PuzzleStateSuccess
)

type DoorPuzzle struct {
	Active bool
	Solved bool

	Level int

	State int

	Sequence []int

	ShowTimer int

	InputIndex int

	LastInput int

	Message string
}

func NewDoorPuzzle(level int) *DoorPuzzle {
	puzzle := &DoorPuzzle{
		Active:     false,
		Solved:     false,
		Level:      level,
		State:      PuzzleStateShowing,
		Sequence:   []int{},
		ShowTimer:  0,
		InputIndex: 0,
		LastInput:  -1,
		Message:    "",
	}

	puzzle.CreateSequence()

	return puzzle
}

func (p *DoorPuzzle) CreateSequence() {
	switch p.Level {
	case 1:
		p.Sequence = []int{
			PuzzleRuneUp,
			PuzzleRuneRight,
			PuzzleRuneDown,
			PuzzleRuneLeft,
		}

	case 2:
		p.Sequence = []int{
			PuzzleRuneLeft,
			PuzzleRuneUp,
			PuzzleRuneRight,
			PuzzleRuneUp,
			PuzzleRuneDown,
		}

	case 3:
		p.Sequence = []int{
			PuzzleRuneDown,
			PuzzleRuneRight,
			PuzzleRuneUp,
			PuzzleRuneLeft,
			PuzzleRuneDown,
			PuzzleRuneRight,
		}

	default:
		p.Sequence = []int{
			PuzzleRuneUp,
			PuzzleRuneRight,
			PuzzleRuneDown,
			PuzzleRuneLeft,
		}
	}
}

func (p *DoorPuzzle) Start() {
	if p.Solved {
		return
	}

	p.Active = true
	p.State = PuzzleStateShowing
	p.ShowTimer = 0
	p.InputIndex = 0
	p.LastInput = -1
	p.Message = "MEMORISE LES RUNES"
}

func (p *DoorPuzzle) RestartSequence() {
	p.State = PuzzleStateShowing
	p.ShowTimer = 0
	p.InputIndex = 0
	p.LastInput = -1
	p.Message = "MEMORISE A NOUVEAU"
}

func (p *DoorPuzzle) Update() {
	if !p.Active {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if !p.Solved {
			p.Active = false
			return
		}
	}

	switch p.State {
	case PuzzleStateShowing:
		p.UpdateShowing()

	case PuzzleStateInput:
		p.UpdateInput()

	case PuzzleStateFailure:
		p.UpdateFailure()

	case PuzzleStateSuccess:
		p.UpdateSuccess()
	}
}

func (p *DoorPuzzle) UpdateShowing() {
	p.ShowTimer++

	totalDuration := len(p.Sequence)*45 + 25

	if p.ShowTimer >= totalDuration {
		p.State = PuzzleStateInput
		p.InputIndex = 0
		p.LastInput = -1
		p.Message = "REPRODUIS LA SEQUENCE"
	}
}

func (p *DoorPuzzle) UpdateInput() {
	input := p.ReadDirectionInput()

	if input == -1 {
		return
	}

	p.LastInput = input

	if p.InputIndex >= len(p.Sequence) {
		return
	}

	expected := p.Sequence[p.InputIndex]

	if input != expected {
		p.State = PuzzleStateFailure
		p.ShowTimer = 0
		p.Message = "MAUVAISE RUNE"
		return
	}

	p.InputIndex++

	if p.InputIndex >= len(p.Sequence) {
		p.Solved = true
		p.State = PuzzleStateSuccess
		p.ShowTimer = 0
		p.Message = "SCEAU BRISE"
	}
}

func (p *DoorPuzzle) UpdateFailure() {
	p.ShowTimer++

	if p.ShowTimer >= 65 {
		p.RestartSequence()
	}
}

func (p *DoorPuzzle) UpdateSuccess() {
	p.ShowTimer++

	if p.ShowTimer >= 90 {
		p.Active = false
	}
}

func (p *DoorPuzzle) ReadDirectionInput() int {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return PuzzleRuneUp
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return PuzzleRuneRight
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return PuzzleRuneDown
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return PuzzleRuneLeft
	}

	return -1
}

func (p *DoorPuzzle) RuneName(runeValue int) string {
	switch runeValue {
	case PuzzleRuneUp:
		return "HAUT"

	case PuzzleRuneRight:
		return "DROITE"

	case PuzzleRuneDown:
		return "BAS"

	case PuzzleRuneLeft:
		return "GAUCHE"
	}

	return "?"
}

func (p *DoorPuzzle) CurrentShowingRune() int {
	if p.State != PuzzleStateShowing {
		return -1
	}

	if p.ShowTimer < 20 {
		return -1
	}

	adjustedTimer := p.ShowTimer - 20

	index := adjustedTimer / 45

	if index < 0 || index >= len(p.Sequence) {
		return -1
	}

	frameInsideRune := adjustedTimer % 45

	if frameInsideRune >= 31 {
		return -1
	}

	return p.Sequence[index]
}

func (p *DoorPuzzle) Draw(screen *ebiten.Image) {
	if !p.Active {
		return
	}

	hudRect(
		screen,
		0,
		0,
		renderWidth,
		renderHeight,
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 190,
		},
	)

	hudRect(
		screen,
		180,
		90,
		920,
		540,
		color.RGBA{
			R: 7,
			G: 10,
			B: 18,
			A: 245,
		},
	)

	hudRect(
		screen,
		180,
		90,
		920,
		4,
		color.RGBA{
			R: 120,
			G: 80,
			B: 200,
			A: 255,
		},
	)

	hudCenteredText(
		screen,
		"SCEAU DU PORTAIL",
		float64(renderWidth)/2,
		125,
	)

	hudCenteredText(
		screen,
		fmt.Sprintf("ENIGME %d", p.Level),
		float64(renderWidth)/2,
		155,
	)

	hudCenteredText(
		screen,
		p.Message,
		float64(renderWidth)/2,
		205,
	)

	switch p.State {
	case PuzzleStateShowing:
		p.DrawShowing(screen)

	case PuzzleStateInput:
		p.DrawInput(screen)

	case PuzzleStateFailure:
		p.DrawFailure(screen)

	case PuzzleStateSuccess:
		p.DrawSuccess(screen)
	}

	if p.State != PuzzleStateSuccess {
		hudCenteredText(
			screen,
			"ECHAP - QUITTER L'ENIGME",
			float64(renderWidth)/2,
			590,
		)
	}
}

func (p *DoorPuzzle) DrawShowing(screen *ebiten.Image) {
	currentRune := p.CurrentShowingRune()

	if currentRune == -1 {
		hudCenteredText(
			screen,
			"...",
			float64(renderWidth)/2,
			330,
		)

		return
	}

	p.DrawLargeRune(screen, currentRune)
}

func (p *DoorPuzzle) DrawLargeRune(screen *ebiten.Image, runeValue int) {
	centerX := float64(renderWidth) / 2

	hudRect(
		screen,
		centerX-145,
		270,
		290,
		150,
		color.RGBA{
			R: 20,
			G: 14,
			B: 38,
			A: 245,
		},
	)

	hudRect(
		screen,
		centerX-145,
		270,
		290,
		4,
		color.RGBA{
			R: 175,
			G: 100,
			B: 255,
			A: 255,
		},
	)

	hudCenteredText(
		screen,
		p.RuneName(runeValue),
		centerX,
		330,
	)

	p.DrawDirectionalSymbol(screen, runeValue, centerX, 380)
}

func (p *DoorPuzzle) DrawDirectionalSymbol(screen *ebiten.Image, runeValue int, centerX float64, centerY float64) {
	var symbol string

	switch runeValue {
	case PuzzleRuneUp:
		symbol = "^"

	case PuzzleRuneRight:
		symbol = ">"

	case PuzzleRuneDown:
		symbol = "V"

	case PuzzleRuneLeft:
		symbol = "<"
	}

	hudCenteredText(
		screen,
		symbol,
		centerX,
		centerY,
	)
}

func (p *DoorPuzzle) DrawInput(screen *ebiten.Image) {
	hudCenteredText(
		screen,
		"UTILISE LES FLECHES DU CLAVIER",
		float64(renderWidth)/2,
		260,
	)

	p.DrawProgress(screen)

	p.DrawDirectionButtons(screen)
}

func (p *DoorPuzzle) DrawProgress(screen *ebiten.Image) {
	count := len(p.Sequence)

	if count <= 0 {
		return
	}

	boxWidth := 70.0
	gap := 15.0

	totalWidth := float64(count)*boxWidth + float64(count-1)*gap
	startX := float64(renderWidth)/2 - totalWidth/2

	for i := 0; i < count; i++ {
		x := startX + float64(i)*(boxWidth+gap)

		boxColor := color.RGBA{
			R: 24,
			G: 25,
			B: 35,
			A: 255,
		}

		text := "?"

		if i < p.InputIndex {
			boxColor = color.RGBA{
				R: 35,
				G: 95,
				B: 55,
				A: 255,
			}

			text = "OK"
		}

		if i == p.InputIndex {
			boxColor = color.RGBA{
				R: 75,
				G: 45,
				B: 120,
				A: 255,
			}
		}

		hudRect(
			screen,
			x,
			320,
			boxWidth,
			55,
			boxColor,
		)

		hudCenteredText(
			screen,
			text,
			x+boxWidth/2,
			339,
		)
	}
}

func (p *DoorPuzzle) DrawDirectionButtons(screen *ebiten.Image) {
	centerX := float64(renderWidth) / 2

	p.DrawDirectionButton(
		screen,
		PuzzleRuneUp,
		centerX-70,
		420,
	)

	p.DrawDirectionButton(
		screen,
		PuzzleRuneLeft,
		centerX-160,
		485,
	)

	p.DrawDirectionButton(
		screen,
		PuzzleRuneDown,
		centerX-70,
		485,
	)

	p.DrawDirectionButton(
		screen,
		PuzzleRuneRight,
		centerX+20,
		485,
	)
}

func (p *DoorPuzzle) DrawDirectionButton(screen *ebiten.Image, runeValue int, x float64, y float64) {
	width := 140.0
	height := 48.0

	buttonColor := color.RGBA{
		R: 25,
		G: 28,
		B: 42,
		A: 255,
	}

	if p.LastInput == runeValue {
		buttonColor = color.RGBA{
			R: 80,
			G: 50,
			B: 125,
			A: 255,
		}
	}

	hudRect(
		screen,
		x,
		y,
		width,
		height,
		buttonColor,
	)

	hudCenteredText(
		screen,
		p.RuneName(runeValue),
		x+width/2,
		y+16,
	)
}

func (p *DoorPuzzle) DrawFailure(screen *ebiten.Image) {
	hudCenteredText(
		screen,
		"LE SCEAU REJETTE TA REPONSE",
		float64(renderWidth)/2,
		310,
	)

	hudCenteredText(
		screen,
		"LA SEQUENCE VA REAPPARAITRE",
		float64(renderWidth)/2,
		355,
	)
}

func (p *DoorPuzzle) DrawSuccess(screen *ebiten.Image) {
	hudCenteredText(
		screen,
		"SCEAU BRISE",
		float64(renderWidth)/2,
		300,
	)

	hudCenteredText(
		screen,
		"LE PORTAIL EST MAINTENANT OUVERT",
		float64(renderWidth)/2,
		350,
	)

	hudRect(
		screen,
		float64(renderWidth)/2-180,
		415,
		360,
		5,
		color.RGBA{
			R: 150,
			G: 90,
			B: 255,
			A: 255,
		},
	)
}
