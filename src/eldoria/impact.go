package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Impact struct {
	X float64
	Y float64

	Timer    int
	MaxTimer int

	PlayerHit bool
}

func NewImpact(
	x float64,
	y float64,
	playerHit bool,
) *Impact {
	return &Impact{
		X: x,
		Y: y,

		Timer:    12,
		MaxTimer: 12,

		PlayerHit: playerHit,
	}
}

func (i *Impact) Update() {
	i.Timer--
}

func (i *Impact) Alive() bool {
	return i.Timer > 0
}

func (i *Impact) Draw(screen *ebiten.Image) {
	if !i.Alive() {
		return
	}

	progress := float64(i.MaxTimer-i.Timer) /
		float64(i.MaxTimer)

	size := 3.0 + progress*9.0

	alpha := uint8(
		float64(255) *
			(1.0 - progress),
	)

	var mainColor color.RGBA

	if i.PlayerHit {
		mainColor = color.RGBA{
			R: 255,
			G: 60,
			B: 60,
			A: alpha,
		}
	} else {
		mainColor = color.RGBA{
			R: 255,
			G: 230,
			B: 150,
			A: alpha,
		}
	}

	// Centre de l'impact.
	ebitenutil.DrawRect(
		screen,
		i.X-2,
		i.Y-2,
		4,
		4,
		mainColor,
	)

	// Trait horizontal.
	ebitenutil.DrawRect(
		screen,
		i.X-size,
		i.Y-1,
		size*2,
		2,
		mainColor,
	)

	// Trait vertical.
	ebitenutil.DrawRect(
		screen,
		i.X-1,
		i.Y-size,
		2,
		size*2,
		mainColor,
	)

	// Diagonale haut gauche.
	ebitenutil.DrawRect(
		screen,
		i.X-size*0.7,
		i.Y-size*0.7,
		3,
		3,
		mainColor,
	)

	// Diagonale bas droite.
	ebitenutil.DrawRect(
		screen,
		i.X+size*0.7,
		i.Y+size*0.7,
		3,
		3,
		mainColor,
	)

	// Diagonale haut droite.
	ebitenutil.DrawRect(
		screen,
		i.X+size*0.7,
		i.Y-size*0.7,
		3,
		3,
		mainColor,
	)

	// Diagonale bas gauche.
	ebitenutil.DrawRect(
		screen,
		i.X-size*0.7,
		i.Y+size*0.7,
		3,
		3,
		mainColor,
	)
}
