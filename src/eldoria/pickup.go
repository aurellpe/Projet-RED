package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Pickup struct {
	X float64
	Y float64

	Width  float64
	Height float64

	HealAmount int

	Collected bool

	BobTimer int
}

func NewHealthPotion(x float64) *Pickup {
	width := 12.0
	height := 16.0

	return &Pickup{
		X: x,

		Y: groundY - height,

		Width:  width,
		Height: height,

		HealAmount: 30,

		Collected: false,

		BobTimer: 0,
	}
}

func (p *Pickup) HitBox() Rect {
	return Rect{
		X: p.X,
		Y: p.Y,
		W: p.Width,
		H: p.Height,
	}
}

func (p *Pickup) Update(player *Player) int {
	if p.Collected {
		return 0
	}

	p.BobTimer++

	if p.BobTimer >= 360 {
		p.BobTimer = 0
	}

	p.Y = groundY - p.Height

	if !player.Alive {
		return 0
	}

	if !Intersects(
		p.HitBox(),
		player.HitBox(),
	) {
		return 0
	}

	if player.HP >= player.MaxHP {
		return 0
	}

	oldHP := player.HP

	player.HP += p.HealAmount

	if player.HP > player.MaxHP {
		player.HP = player.MaxHP
	}

	healed := player.HP - oldHP

	p.Collected = true

	return healed
}

func (p *Pickup) Draw(screen *ebiten.Image) {
	if p.Collected {
		return
	}

	bob := math.Sin(
		float64(p.BobTimer)*0.08,
	) * 2

	drawY := p.Y + bob

	// Ombre
	ebitenutil.DrawRect(
		screen,
		p.X+2,
		groundY-2,
		p.Width-4,
		2,
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 100,
		},
	)

	// Bouchon
	ebitenutil.DrawRect(
		screen,
		p.X+4,
		drawY,
		4,
		3,
		color.RGBA{
			R: 190,
			G: 160,
			B: 80,
			A: 255,
		},
	)

	// Haut de la fiole
	ebitenutil.DrawRect(
		screen,
		p.X+3,
		drawY+3,
		6,
		3,
		color.RGBA{
			R: 220,
			G: 220,
			B: 235,
			A: 255,
		},
	)

	// Corps
	ebitenutil.DrawRect(
		screen,
		p.X+1,
		drawY+6,
		10,
		9,
		color.RGBA{
			R: 200,
			G: 30,
			B: 60,
			A: 255,
		},
	)

	// Reflet
	ebitenutil.DrawRect(
		screen,
		p.X+3,
		drawY+7,
		2,
		5,
		color.RGBA{
			R: 255,
			G: 150,
			B: 170,
			A: 255,
		},
	)
}
