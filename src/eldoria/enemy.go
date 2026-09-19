package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets/enemy/enemy1.png
var enemy1Data []byte

//go:embed assets/enemy/enemy2.png
var enemy2Data []byte

//go:embed assets/enemy/enemy3.png
var enemy3Data []byte

//go:embed assets/enemy/enemy4.png
var enemy4Data []byte

type Enemy struct {
	X float64
	Y float64

	Width  float64
	Height float64

	Scale float64
	Speed float64

	SpriteYOffset float64

	HP    int
	MaxHP int

	Alive bool

	// --------------------------------------------------
	// MORT
	// --------------------------------------------------

	Dying bool

	DeathTimer int

	DeathDuration int

	DeathRewarded bool

	// --------------------------------------------------
	// ANIMATION
	// --------------------------------------------------

	Frames []*ebiten.Image

	CurrentFrame int
	FrameTimer   int
	FrameSpeed   int

	FacingLeft bool
	Moving     bool

	KnockbackX float64

	FlashTimer int

	// --------------------------------------------------
	// ATTAQUE
	// --------------------------------------------------

	Attacking bool

	AttackTimer int

	AttackCooldown int

	AttackHit bool
}

func loadEnemyImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(
		bytes.NewReader(data),
	)

	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func NewEnemy(x float64) *Enemy {
	frames := []*ebiten.Image{
		loadEnemyImage(enemy1Data),
		loadEnemyImage(enemy2Data),
		loadEnemyImage(enemy3Data),
		loadEnemyImage(enemy4Data),
	}

	return &Enemy{
		X: x,
		Y: groundY - 55,

		Width:  28,
		Height: 55,

		Scale: 0.055,
		Speed: 0.45,

		SpriteYOffset: -12,

		HP:    100,
		MaxHP: 100,

		Alive: true,

		Dying: false,

		DeathTimer: 0,

		DeathDuration: 35,

		DeathRewarded: false,

		Frames: frames,

		CurrentFrame: 0,

		FrameTimer: 0,

		FrameSpeed: 8,

		FacingLeft: true,

		Moving: false,

		KnockbackX: 0,

		FlashTimer: 0,

		Attacking: false,

		AttackTimer: 0,

		AttackCooldown: 0,

		AttackHit: false,
	}
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func (e *Enemy) Update(player *Player) int {
	if !e.Alive {
		return 0
	}

	// ==================================================
	// ANIMATION DE MORT
	// ==================================================

	if e.Dying {
		e.UpdateDeath()

		return 0
	}

	if e.FlashTimer > 0 {
		e.FlashTimer--
	}

	if e.AttackCooldown > 0 {
		e.AttackCooldown--
	}

	// --------------------------------------------------
	// RECUL
	// --------------------------------------------------

	e.X += e.KnockbackX

	e.KnockbackX *= 0.82

	if e.KnockbackX > -0.05 &&
		e.KnockbackX < 0.05 {

		e.KnockbackX = 0
	}

	// --------------------------------------------------
	// JOUEUR MORT
	// --------------------------------------------------

	if !player.Alive {
		e.Attacking = false

		e.Moving = false

		e.CurrentFrame = 0

		return 0
	}

	// --------------------------------------------------
	// ATTAQUE EN COURS
	// --------------------------------------------------

	if e.Attacking {
		e.Moving = false

		return e.UpdateAttack(player)
	}

	playerCenter := player.X +
		player.Width/2

	enemyCenter := e.X +
		e.Width/2

	distance := playerCenter -
		enemyCenter

	// --------------------------------------------------
	// DIRECTION
	// --------------------------------------------------

	if distance < 0 {
		e.FacingLeft = true
	} else {
		e.FacingLeft = false
	}

	e.Moving = false

	// --------------------------------------------------
	// POURSUITE
	// --------------------------------------------------

	if distance > 28 {
		e.X += e.Speed

		e.Moving = true
	}

	if distance < -28 {
		e.X -= e.Speed

		e.Moving = true
	}

	// --------------------------------------------------
	// ANIMATION MARCHE
	// --------------------------------------------------

	if e.Moving {
		e.FrameTimer++

		if e.FrameTimer >= e.FrameSpeed {
			e.FrameTimer = 0

			if e.CurrentFrame == 1 {
				e.CurrentFrame = 2
			} else {
				e.CurrentFrame = 1
			}
		}
	} else {
		e.CurrentFrame = 0

		e.FrameTimer = 0
	}

	// --------------------------------------------------
	// COMMENCE UNE ATTAQUE
	// --------------------------------------------------

	if distance <= 32 &&
		distance >= -32 &&
		e.AttackCooldown <= 0 {

		e.Attacking = true

		e.AttackTimer = 0

		e.AttackHit = false

		e.CurrentFrame = 3
	}

	// --------------------------------------------------
	// LIMITES
	// --------------------------------------------------

	if e.X < 0 {
		e.X = 0
	}

	if e.X+e.Width > screenWidth {
		e.X = screenWidth - e.Width
	}

	e.Y = groundY - e.Height

	return 0
}

// --------------------------------------------------
// MORT
// --------------------------------------------------

func (e *Enemy) UpdateDeath() {
	e.DeathTimer++

	e.Attacking = false

	e.Moving = false

	e.CurrentFrame = 0

	e.AttackTimer = 0

	e.AttackHit = false

	// Le recul continue légèrement pendant la mort.
	e.X += e.KnockbackX

	e.KnockbackX *= 0.90

	if e.DeathTimer >= e.DeathDuration {
		e.Alive = false
	}
}

// --------------------------------------------------
// ATTAQUE
// --------------------------------------------------

func (e *Enemy) UpdateAttack(player *Player) int {
	e.AttackTimer++

	e.CurrentFrame = 3

	if e.AttackTimer >= 7 &&
		e.AttackTimer <= 12 &&
		!e.AttackHit {

		attackBox := e.AttackBox()

		playerBox := player.HitBox()

		if Intersects(
			attackBox,
			playerBox,
		) {
			e.AttackHit = true

			enemyCenter := e.X +
				e.Width/2

			if player.Hit(
				10,
				enemyCenter,
			) {
				return 10
			}
		}
	}

	if e.AttackTimer > 18 {
		e.Attacking = false

		e.AttackTimer = 0

		e.AttackHit = false

		e.AttackCooldown = 45

		e.CurrentFrame = 0
	}

	return 0
}

// --------------------------------------------------
// HITBOX
// --------------------------------------------------

func (e *Enemy) HitBox() Rect {
	return Rect{
		X: e.X,
		Y: e.Y,
		W: e.Width,
		H: e.Height,
	}
}

// --------------------------------------------------
// HITBOX ATTAQUE
// --------------------------------------------------

func (e *Enemy) AttackBox() Rect {
	attackWidth := 32.0

	attackHeight := 55.0

	if e.FacingLeft {
		return Rect{
			X: e.X - attackWidth,

			Y: groundY - attackHeight,

			W: attackWidth,

			H: attackHeight,
		}
	}

	return Rect{
		X: e.X + e.Width,

		Y: groundY - attackHeight,

		W: attackWidth,

		H: attackHeight,
	}
}

// --------------------------------------------------
// PREND DES DEGATS
// --------------------------------------------------

func (e *Enemy) Hit(
	damage int,
	attackerX float64,
) bool {

	if !e.Alive {
		return false
	}

	if e.Dying {
		return false
	}

	e.HP -= damage

	e.FlashTimer = 6

	enemyCenter := e.X +
		e.Width/2

	// --------------------------------------------------
	// RECUL
	// --------------------------------------------------

	if enemyCenter < attackerX {
		e.KnockbackX = -1.8
	} else {
		e.KnockbackX = 1.8
	}

	// --------------------------------------------------
	// COMMENCE LA MORT
	// --------------------------------------------------

	if e.HP <= 0 {
		e.HP = 0

		e.Dying = true

		e.DeathTimer = 0

		e.Attacking = false

		e.Moving = false

		e.CurrentFrame = 0
	}

	return true
}

// --------------------------------------------------
// DRAW
// --------------------------------------------------

func (e *Enemy) Draw(
	screen *ebiten.Image,
) {
	if !e.Alive {
		return
	}

	img := e.Frames[e.CurrentFrame]

	if e.Dying {
		img = e.Frames[0]
	}

	frameWidth := float64(
		img.Bounds().Dx(),
	) * e.Scale

	frameHeight := float64(
		img.Bounds().Dy(),
	) * e.Scale

	centerX := e.X +
		e.Width/2

	drawX := centerX -
		frameWidth/2

	drawY := groundY -
		frameHeight +
		e.SpriteYOffset

	options := &ebiten.DrawImageOptions{}

	// ==================================================
	// ANIMATION DE MORT
	// ==================================================

	deathAlpha := float32(1)

	if e.Dying {
		progress := float64(
			e.DeathTimer,
		) / float64(
			e.DeathDuration,
		)

		if progress > 1 {
			progress = 1
		}

		// L'ennemi s'enfonce progressivement.
		drawY += progress * 14

		deathAlpha = float32(
			1 - progress,
		)
	}

	// --------------------------------------------------
	// ORIENTATION
	// --------------------------------------------------

	if e.FacingLeft {
		options.GeoM.Scale(
			-e.Scale,
			e.Scale,
		)

		options.GeoM.Translate(
			drawX+frameWidth,
			drawY,
		)
	} else {
		options.GeoM.Scale(
			e.Scale,
			e.Scale,
		)

		options.GeoM.Translate(
			drawX,
			drawY,
		)
	}

	// --------------------------------------------------
	// FLASH
	// --------------------------------------------------

	if e.FlashTimer > 0 {
		options.ColorScale.Scale(
			2,
			2,
			2,
			1,
		)
	}

	if e.Dying {
		options.ColorScale.Scale(
			1,
			1,
			1,
			deathAlpha,
		)
	}

	screen.DrawImage(
		img,
		options,
	)

	// Pas de barre de vie pendant la mort.
	if e.Dying {
		return
	}

	// --------------------------------------------------
	// BARRE DE VIE
	// --------------------------------------------------

	barY := e.Y +
		e.SpriteYOffset -
		5

	ebitenutil.DrawRect(
		screen,
		e.X,
		barY,
		e.Width,
		3,
		color.RGBA{
			R: 30,
			G: 30,
			B: 30,
			A: 255,
		},
	)

	ratio := float64(e.HP) /
		float64(e.MaxHP)

	if ratio < 0 {
		ratio = 0
	}

	ebitenutil.DrawRect(
		screen,
		e.X,
		barY,
		e.Width*ratio,
		3,
		color.RGBA{
			R: 220,
			G: 40,
			B: 40,
			A: 255,
		},
	)
}
