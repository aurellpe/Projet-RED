package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DamageText struct {
	X float64
	Y float64

	Damage int

	Timer int
}

func NewDamageText(x, y float64, damage int) *DamageText {
	return &DamageText{
		X:      x,
		Y:      y,
		Damage: damage,
		Timer:  45,
	}
}

func (d *DamageText) Update() {
	// Le texte monte doucement.
	d.Y -= 0.35

	d.Timer--
}

func (d *DamageText) Alive() bool {
	return d.Timer > 0
}

func (d *DamageText) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("-%d", d.Damage),
		int(d.X),
		int(d.Y),
	)
}