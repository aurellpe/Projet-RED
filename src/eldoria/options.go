package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	OptionsActionNone = iota
	OptionsActionBack
)

type OptionsMenu struct {
	Selected      int
	MusicVolume   int
	EffectsVolume int
	Fullscreen    bool
	Timer         int
}

func NewOptionsMenu(music *MusicManager, sounds *SoundManager) *OptionsMenu {
	musicVolume := 20
	effectsVolume := 70

	if music != nil {
		musicVolume = int(music.Volume * 100)
	}

	if sounds != nil {
		effectsVolume = int(sounds.Volume * 100)
	}

	return &OptionsMenu{
		Selected:      0,
		MusicVolume:   musicVolume,
		EffectsVolume: effectsVolume,
		Fullscreen:    ebiten.IsFullscreen(),
		Timer:         0,
	}
}

func clampVolume(value int) int {
	if value < 0 {
		return 0
	}

	if value > 100 {
		return 100
	}

	return value
}

func (o *OptionsMenu) ApplyMusicVolume(music *MusicManager) {
	if music == nil {
		return
	}

	music.SetVolume(float64(o.MusicVolume) / 100)
}

func (o *OptionsMenu) ApplyEffectsVolume(sounds *SoundManager) {
	if sounds == nil {
		return
	}

	sounds.SetVolume(float64(o.EffectsVolume) / 100)
}

func (o *OptionsMenu) ToggleFullscreen() {
	o.Fullscreen = !o.Fullscreen
	ebiten.SetFullscreen(o.Fullscreen)
}

func (o *OptionsMenu) Update(music *MusicManager, sounds *SoundManager) int {
	o.Timer++

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return OptionsActionBack
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		o.Selected--

		if o.Selected < 0 {
			o.Selected = 2
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		o.Selected++

		if o.Selected > 2 {
			o.Selected = 0
		}
	}

	moveLeft := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft)
	moveRight := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight)

	if o.Selected == 0 {
		if moveLeft {
			o.MusicVolume = clampVolume(o.MusicVolume - 10)
			o.ApplyMusicVolume(music)
		}

		if moveRight {
			o.MusicVolume = clampVolume(o.MusicVolume + 10)
			o.ApplyMusicVolume(music)
		}
	}

	if o.Selected == 1 {
		if moveLeft {
			o.EffectsVolume = clampVolume(o.EffectsVolume - 10)
			o.ApplyEffectsVolume(sounds)
		}

		if moveRight {
			o.EffectsVolume = clampVolume(o.EffectsVolume + 10)
			o.ApplyEffectsVolume(sounds)
		}
	}

	if o.Selected == 2 {
		if moveLeft || moveRight || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			o.ToggleFullscreen()
		}
	}

	return OptionsActionNone
}

func drawOptionRow(screen *ebiten.Image, label string, value string, y float64, selected bool, timer int) {
	x := 55.0
	width := 210.0
	height := 25.0

	if selected {
		pulseValue := 175 + int(55*math.Sin(float64(timer)*0.08))

		if pulseValue < 0 {
			pulseValue = 0
		}

		if pulseValue > 255 {
			pulseValue = 255
		}

		pulse := uint8(pulseValue)

		ebitenutil.DrawRect(screen, x-3, y-3, width+6, height+6, color.RGBA{R: 20, G: 100, B: 255, A: 55})
		ebitenutil.DrawRect(screen, x-2, y-2, width+4, height+4, color.RGBA{R: 50, G: 170, B: 255, A: pulse})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 25, B: 55, A: 240})
	} else {
		ebitenutil.DrawRect(screen, x-1, y-1, width+2, height+2, color.RGBA{R: 180, G: 140, B: 65, A: 170})
		ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{R: 5, G: 12, B: 25, A: 225})
	}

	ebitenutil.DebugPrintAt(screen, label, int(x)+12, int(y)+8)

	valueWidth := len(value) * 6
	valueX := int(x+width) - valueWidth - 12

	ebitenutil.DebugPrintAt(screen, value, valueX, int(y)+8)
}

func (o *OptionsMenu) Draw(screen *ebiten.Image, background *ebiten.Image, fromPause bool) {
	drawMenuBackgroundCover(screen, background)

	ebitenutil.DrawRect(screen, 0, 0, 320, 180, color.RGBA{R: 0, G: 5, B: 15, A: 100})

	ebitenutil.DrawRect(screen, 40, 10, 240, 160, color.RGBA{R: 3, G: 10, B: 25, A: 225})

	ebitenutil.DrawRect(screen, 39, 9, 242, 2, color.RGBA{R: 210, G: 165, B: 70, A: 255})
	ebitenutil.DrawRect(screen, 39, 169, 242, 2, color.RGBA{R: 210, G: 165, B: 70, A: 255})

	drawCenteredMenuText(screen, "OPTIONS", 20)

	ebitenutil.DrawRect(screen, 105, 35, 110, 1, color.RGBA{R: 60, G: 165, B: 255, A: 255})

	musicValue := fmt.Sprintf("< %d%% >", o.MusicVolume)
	effectsValue := fmt.Sprintf("< %d%% >", o.EffectsVolume)

	fullscreenValue := "< NON >"

	if o.Fullscreen {
		fullscreenValue = "< OUI >"
	}

	drawOptionRow(screen, "MUSIQUE", musicValue, 48, o.Selected == 0, o.Timer)
	drawOptionRow(screen, "EFFETS", effectsValue, 82, o.Selected == 1, o.Timer)
	drawOptionRow(screen, "PLEIN ECRAN", fullscreenValue, 116, o.Selected == 2, o.Timer)

	ebitenutil.DebugPrintAt(screen, "HAUT / BAS : CHOISIR", 99, 148)
	ebitenutil.DebugPrintAt(screen, "GAUCHE / DROITE : MODIFIER", 76, 158)

	if fromPause {
		ebitenutil.DebugPrintAt(screen, "ECHAP : RETOUR PAUSE", 100, 168)
	} else {
		ebitenutil.DebugPrintAt(screen, "ECHAP : RETOUR MENU", 103, 168)
	}
}
