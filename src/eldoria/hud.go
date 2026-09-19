package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const renderWidth = 1280
const renderHeight = 720
const hudTextScale = 1.20

var hudTextBuffer *ebiten.Image

func hudRect(screen *ebiten.Image, x float64, y float64, w float64, h float64, c color.RGBA) {
	ebitenutil.DrawRect(screen, x, y, w, h, c)
}

func drawHUDText(screen *ebiten.Image, text string, x float64, y float64) {
	if hudTextBuffer == nil {
		hudTextBuffer = ebiten.NewImage(700, 24)
	}

	hudTextBuffer.Clear()

	ebitenutil.DebugPrintAt(hudTextBuffer, text, 0, 0)

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(hudTextScale, hudTextScale)
	options.GeoM.Translate(x, y)

	options.Filter = ebiten.FilterNearest

	screen.DrawImage(hudTextBuffer, options)
}

func hudCenteredText(screen *ebiten.Image, text string, centerX float64, y float64) {
	textWidth := float64(len(text)*6) * hudTextScale

	x := centerX - textWidth/2

	drawHUDText(screen, text, x, y)
}

func drawMinimalBar(screen *ebiten.Image, x float64, y float64, w float64, h float64, ratio float64) {
	if ratio < 0 {
		ratio = 0
	}

	if ratio > 1 {
		ratio = 1
	}

	hudRect(
		screen,
		x,
		y,
		w,
		h,
		color.RGBA{
			R: 15,
			G: 15,
			B: 20,
			A: 210,
		},
	)

	hudRect(
		screen,
		x+1,
		y+1,
		w-2,
		h-2,
		color.RGBA{
			R: 55,
			G: 15,
			B: 20,
			A: 255,
		},
	)

	if ratio > 0 {
		fillWidth := (w - 2) * ratio

		hudRect(
			screen,
			x+1,
			y+1,
			fillWidth,
			h-2,
			color.RGBA{
				R: 210,
				G: 40,
				B: 50,
				A: 255,
			},
		)

		hudRect(
			screen,
			x+1,
			y+1,
			fillWidth,
			2,
			color.RGBA{
				R: 255,
				G: 100,
				B: 105,
				A: 255,
			},
		)
	}
}

func (g *Game) DrawHUD(screen *ebiten.Image) {
	if g.Player == nil {
		return
	}

	g.drawPlayerHUD(screen)
	g.drawLevelHUD(screen)

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		g.drawBossHUD(screen)
	} else {
		g.drawEnemyHUD(screen)
	}

	g.drawExitHUD(screen)
	g.drawHealHUD(screen)
	g.drawUpgradeHUD(screen)

	if !g.Player.Alive {
		g.drawGameOverHUD(screen)
		return
	}

	if g.Boss != nil && !g.Boss.Alive {
		g.drawVictoryHUD(screen)
	}

	if g.Transitioning {
		g.drawTransitionHUD(screen)
	}
}

func (g *Game) drawPlayerHUD(screen *ebiten.Image) {
	ratio := float64(g.Player.HP) / float64(g.Player.MaxHP)

	hudRect(
		screen,
		18,
		18,
		205,
		47,
		color.RGBA{
			R: 5,
			G: 8,
			B: 14,
			A: 125,
		},
	)

	drawHUDText(screen, "PV", 28, 27)

	drawMinimalBar(screen, 57, 28, 148, 9, ratio)

	drawHUDText(
		screen,
		fmt.Sprintf("%d/%d", g.Player.HP, g.Player.MaxHP),
		57,
		45,
	)
}

func (g *Game) drawLevelHUD(screen *ebiten.Image) {
	level := g.Levels[g.CurrentLevel]

	hudCenteredText(
		screen,
		level.Name,
		float64(renderWidth)/2,
		20,
	)

	hudCenteredText(
		screen,
		fmt.Sprintf("NIVEAU %d", level.Number),
		float64(renderWidth)/2,
		39,
	)
}

func (g *Game) drawEnemyHUD(screen *ebiten.Image) {
	text := fmt.Sprintf("ENNEMIS  %d", g.EnemiesRemaining())

	textWidth := float64(len(text)*6) * hudTextScale

	x := float64(renderWidth) - textWidth - 25

	drawHUDText(screen, text, x, 27)
}

func (g *Game) drawExitHUD(screen *ebiten.Image) {
	if g.Exit == nil {
		return
	}

	if g.Exit.Active {
		text := "PORTAIL OUVERT"

		textWidth := float64(len(text)*6) * hudTextScale

		x := float64(renderWidth) - textWidth - 25

		drawHUDText(screen, text, x, 48)
	}

	if !g.Exit.PlayerInside(g.Player) {
		return
	}

	text := "PORTAIL VERROUILLE"

	if g.Exit.Active {
		text = "E  -  ENTRER"
	}

	textWidth := float64(len(text)*6) * hudTextScale

	width := textWidth + 32

	x := float64(renderWidth)/2 - width/2

	hudRect(
		screen,
		x,
		660,
		width,
		30,
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 145,
		},
	)

	hudCenteredText(
		screen,
		text,
		float64(renderWidth)/2,
		668,
	)
}

func (g *Game) drawHealHUD(screen *ebiten.Image) {
	if g.HealMessageTimer <= 0 {
		return
	}

	drawHUDText(
		screen,
		fmt.Sprintf("+%d PV", g.LastHealAmount),
		25,
		78,
	)
}

func (g *Game) drawUpgradeHUD(screen *ebiten.Image) {
	if g.UpgradeMessageTimer <= 0 {
		return
	}

	if g.UpgradeMessage == "" {
		return
	}

	textWidth := float64(len(g.UpgradeMessage)*6) * hudTextScale

	width := textWidth + 40

	x := float64(renderWidth)/2 - width/2

	hudRect(
		screen,
		x,
		90,
		width,
		30,
		color.RGBA{
			R: 12,
			G: 20,
			B: 30,
			A: 165,
		},
	)

	hudCenteredText(
		screen,
		g.UpgradeMessage,
		float64(renderWidth)/2,
		99,
	)
}

func (g *Game) drawBossHUD(screen *ebiten.Image) {
	if g.Boss == nil {
		return
	}

	if !g.Boss.Alive {
		return
	}

	ratio := float64(g.Boss.HP) / float64(g.Boss.MaxHP)

	hudCenteredText(
		screen,
		"GARDIEN D'ELDORIA",
		float64(renderWidth)/2,
		20,
	)

	drawMinimalBar(
		screen,
		470,
		45,
		340,
		10,
		ratio,
	)

	hudCenteredText(
		screen,
		fmt.Sprintf("%d/%d", g.Boss.HP, g.Boss.MaxHP),
		float64(renderWidth)/2,
		62,
	)

	if g.Boss.Phase >= 2 {
		hudCenteredText(
			screen,
			"PHASE 2",
			float64(renderWidth)/2,
			80,
		)
	}
}

func (g *Game) drawGameOverHUD(screen *ebiten.Image) {
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
			A: 160,
		},
	)

	hudCenteredText(
		screen,
		"GAME OVER",
		float64(renderWidth)/2,
		300,
	)

	hudCenteredText(
		screen,
		"LE HEROS EST TOMBE",
		float64(renderWidth)/2,
		330,
	)

	hudCenteredText(
		screen,
		"ENTREE  -  REJOUER",
		float64(renderWidth)/2,
		370,
	)

	hudCenteredText(
		screen,
		"ECHAP  -  MENU",
		float64(renderWidth)/2,
		400,
	)
}

func (g *Game) drawVictoryHUD(screen *ebiten.Image) {
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
			A: 175,
		},
	)

	hudCenteredText(
		screen,
		"VICTOIRE",
		float64(renderWidth)/2,
		260,
	)

	hudCenteredText(
		screen,
		"LE GARDIEN EST VAINCU",
		float64(renderWidth)/2,
		305,
	)

	hudCenteredText(
		screen,
		"ELDORIA EST LIBEREE",
		float64(renderWidth)/2,
		345,
	)

	hudCenteredText(
		screen,
		"LA LUMIERE REVIENT SUR LE ROYAUME",
		float64(renderWidth)/2,
		385,
	)

	hudCenteredText(
		screen,
		"TON AVENTURE S'ACHEVE ICI...",
		float64(renderWidth)/2,
		425,
	)

	hudCenteredText(
		screen,
		"ENTREE  -  RETOUR AU MENU",
		float64(renderWidth)/2,
		485,
	)
}

func (g *Game) drawTransitionHUD(screen *ebiten.Image) {
	alpha := 0.0

	if g.TransitionTimer <= 20 {
		alpha = float64(g.TransitionTimer) / 20
	} else {
		alpha = float64(40-g.TransitionTimer) / 20
	}

	if alpha < 0 {
		alpha = 0
	}

	if alpha > 1 {
		alpha = 1
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
			A: uint8(alpha * 255),
		},
	)

	if g.TransitionTimer >= 17 && g.TransitionTimer <= 27 {
		levelIndex := g.NextLevel

		if g.TransitionTimer >= 20 {
			levelIndex = g.CurrentLevel
		}

		if levelIndex >= 0 && levelIndex < len(g.Levels) {
			level := g.Levels[levelIndex]

			hudCenteredText(
				screen,
				"NOUVELLE ZONE",
				float64(renderWidth)/2,
				315,
			)

			hudCenteredText(
				screen,
				level.Name,
				float64(renderWidth)/2,
				350,
			)

			hudCenteredText(
				screen,
				fmt.Sprintf("NIVEAU %d", level.Number),
				float64(renderWidth)/2,
				380,
			)
		}
	}
}
