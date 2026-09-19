package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --------------------------------------------------
// CONSTANTES
// --------------------------------------------------

const (
	screenWidth = 320
	gravity     = 0.35
	jumpForce   = -6.5
)

var groundY = 156.0

// --------------------------------------------------
// MARCHE
// --------------------------------------------------

//go:embed assets/player/walk1.png
var playerWalk1Data []byte

//go:embed assets/player/walk2.png
var playerWalk2Data []byte

//go:embed assets/player/walk3.png
var playerWalk3Data []byte

//go:embed assets/player/walk4.png
var playerWalk4Data []byte

// --------------------------------------------------
// ATTAQUE
// --------------------------------------------------

//go:embed assets/player/attack1.png
var playerAttack1Data []byte

//go:embed assets/player/attack2.png
var playerAttack2Data []byte

//go:embed assets/player/attack3.png
var playerAttack3Data []byte

//go:embed assets/player/attack4.png
var playerAttack4Data []byte

type Player struct {
	X float64
	Y float64

	Speed float64

	VelocityY float64
	OnGround  bool

	KnockbackX float64

	Scale float64

	Width  float64
	Height float64

	// --------------------------------------------------
	// MARCHE
	// --------------------------------------------------

	WalkFrames []*ebiten.Image

	CurrentWalkFrame int
	WalkFrameTimer   int
	WalkFrameSpeed   int

	Moving bool

	// --------------------------------------------------
	// DIRECTION
	// --------------------------------------------------

	FacingLeft bool

	// --------------------------------------------------
	// ATTAQUE
	// --------------------------------------------------

	AttackFrames []*ebiten.Image

	CurrentAttackFrame int

	Attacking bool

	AttackTimer int

	AttackDuration int

	AttackFrameSpeed int

	AttackHit bool

	// --------------------------------------------------
	// DEGATS / STUN
	// --------------------------------------------------

	HurtTimer int

	HurtDuration int

	// --------------------------------------------------
	// VIE
	// --------------------------------------------------

	HP    int
	MaxHP int

	Alive bool

	InvincibleTimer int

	FlashTimer int
}

// --------------------------------------------------
// CHARGEMENT IMAGE
// --------------------------------------------------

func loadPlayerImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// --------------------------------------------------
// CREATION PLAYER
// --------------------------------------------------

func NewPlayer() *Player {
	walkFrames := []*ebiten.Image{
		loadPlayerImage(playerWalk1Data),
		loadPlayerImage(playerWalk2Data),
		loadPlayerImage(playerWalk3Data),
		loadPlayerImage(playerWalk4Data),
	}

	attackFrames := []*ebiten.Image{
		loadPlayerImage(playerAttack1Data),
		loadPlayerImage(playerAttack2Data),
		loadPlayerImage(playerAttack3Data),
		loadPlayerImage(playerAttack4Data),
	}

	scale := 0.08

	idle := walkFrames[0]

	width := float64(idle.Bounds().Dx()) * scale
	height := float64(idle.Bounds().Dy()) * scale

	return &Player{
		X: 70,
		Y: groundY - height,

		Speed: 1.5,

		VelocityY: 0,
		OnGround:  true,

		KnockbackX: 0,

		Scale: scale,

		Width:  width,
		Height: height,

		WalkFrames: walkFrames,

		CurrentWalkFrame: 0,

		WalkFrameTimer: 0,

		WalkFrameSpeed: 6,

		Moving: false,

		FacingLeft: false,

		AttackFrames: attackFrames,

		CurrentAttackFrame: 0,

		Attacking: false,

		AttackTimer: 0,

		AttackDuration: 24,

		AttackFrameSpeed: 6,

		AttackHit: false,

		HurtTimer: 0,

		HurtDuration: 12,

		HP:    100,
		MaxHP: 100,

		Alive: true,

		InvincibleTimer: 0,

		FlashTimer: 0,
	}
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func (p *Player) Update() {
	// ==================================================
	// TIMERS
	// ==================================================

	if p.InvincibleTimer > 0 {
		p.InvincibleTimer--
	}

	if p.FlashTimer > 0 {
		p.FlashTimer--
	}

	if p.HurtTimer > 0 {
		p.HurtTimer--
	}

	if !p.Alive {
		return
	}

	p.Moving = false

	// ==================================================
	// RECUL
	// ==================================================

	p.X += p.KnockbackX

	p.KnockbackX *= 0.82

	if p.KnockbackX > -0.05 &&
		p.KnockbackX < 0.05 {

		p.KnockbackX = 0
	}

	// ==================================================
	// JOUEUR BLESSE
	//
	// Pendant quelques frames :
	// pas de marche
	// pas d'attaque
	// mais la gravité continue.
	// ==================================================

	if p.HurtTimer > 0 {
		p.Attacking = false

		p.AttackTimer = 0

		p.CurrentAttackFrame = 0

		p.AttackHit = false

		p.CurrentWalkFrame = 0

		p.WalkFrameTimer = 0

		p.UpdateGravity()

		p.KeepInsideScreen()

		return
	}

	// ==================================================
	// ATTAQUE EN COURS
	// ==================================================

	if p.Attacking {
		p.UpdateAttackAnimation()

		p.UpdateGravity()

		p.KeepInsideScreen()

		return
	}

	// ==================================================
	// DEPLACEMENT
	// ==================================================

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= p.Speed

		p.Moving = true

		p.FacingLeft = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += p.Speed

		p.Moving = true

		p.FacingLeft = false
	}

	// ==================================================
	// SAUT
	// ==================================================

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) &&
		p.OnGround {

		p.VelocityY = jumpForce

		p.OnGround = false
	}

	// ==================================================
	// ATTAQUE
	// ==================================================

	if inpututil.IsMouseButtonJustPressed(
		ebiten.MouseButtonLeft,
	) {
		p.StartAttack()
	}

	// ==================================================
	// GRAVITE
	// ==================================================

	p.UpdateGravity()

	// ==================================================
	// ANIMATION MARCHE
	// ==================================================

	p.UpdateWalkAnimation()

	// ==================================================
	// LIMITES
	// ==================================================

	p.KeepInsideScreen()
}

// --------------------------------------------------
// COMMENCER ATTAQUE
// --------------------------------------------------

func (p *Player) StartAttack() {
	if p.Attacking {
		return
	}

	if p.HurtTimer > 0 {
		return
	}

	if !p.Alive {
		return
	}

	p.Attacking = true

	p.AttackTimer = 0

	p.CurrentAttackFrame = 0

	p.AttackHit = false

	p.Moving = false

	p.CurrentWalkFrame = 0

	p.WalkFrameTimer = 0
}

// --------------------------------------------------
// ANIMATION ATTAQUE
// --------------------------------------------------

func (p *Player) UpdateAttackAnimation() {
	p.AttackTimer++

	frame := p.AttackTimer /
		p.AttackFrameSpeed

	if frame >= len(p.AttackFrames) {
		frame = len(p.AttackFrames) - 1
	}

	p.CurrentAttackFrame = frame

	if p.AttackTimer >= p.AttackDuration {
		p.Attacking = false

		p.AttackTimer = 0

		p.CurrentAttackFrame = 0

		p.AttackHit = false
	}
}

// --------------------------------------------------
// ANIMATION MARCHE
// --------------------------------------------------

func (p *Player) UpdateWalkAnimation() {
	if p.Moving && p.OnGround {
		p.WalkFrameTimer++

		if p.WalkFrameTimer >= p.WalkFrameSpeed {
			p.WalkFrameTimer = 0

			p.CurrentWalkFrame++

			if p.CurrentWalkFrame >= len(p.WalkFrames) {
				p.CurrentWalkFrame = 0
			}
		}
	} else {
		p.CurrentWalkFrame = 0

		p.WalkFrameTimer = 0
	}
}

// --------------------------------------------------
// GRAVITE
// --------------------------------------------------

func (p *Player) UpdateGravity() {
	p.VelocityY += gravity

	p.Y += p.VelocityY

	if p.Y+p.Height >= groundY {
		p.Y = groundY - p.Height

		p.VelocityY = 0

		p.OnGround = true
	} else {
		p.OnGround = false
	}
}

// --------------------------------------------------
// LIMITES
// --------------------------------------------------

func (p *Player) KeepInsideScreen() {
	if p.X < 0 {
		p.X = 0
	}

	if p.X+p.Width > screenWidth {
		p.X = screenWidth - p.Width
	}
}

// --------------------------------------------------
// HITBOX
// --------------------------------------------------

func (p *Player) HitBox() Rect {
	bodyWidth := 28.0
	bodyHeight := 55.0

	playerCenter := p.X + p.Width/2

	playerBottom := p.Y + p.Height

	return Rect{
		X: playerCenter - bodyWidth/2,

		Y: playerBottom - bodyHeight,

		W: bodyWidth,

		H: bodyHeight,
	}
}

// --------------------------------------------------
// HITBOX ATTAQUE
// --------------------------------------------------

func (p *Player) AttackBox() Rect {
	body := p.HitBox()

	attackWidth := 70.0

	attackHeight := 55.0

	if p.FacingLeft {
		return Rect{
			X: body.X - attackWidth,

			Y: body.Y,

			W: attackWidth,

			H: attackHeight,
		}
	}

	return Rect{
		X: body.X + body.W,

		Y: body.Y,

		W: attackWidth,

		H: attackHeight,
	}
}

// --------------------------------------------------
// PLAYER PREND DES DEGATS
// --------------------------------------------------

func (p *Player) Hit(
	damage int,
	sourceX float64,
) bool {

	if !p.Alive {
		return false
	}

	if p.InvincibleTimer > 0 {
		return false
	}

	// --------------------------------------------------
	// DEGATS
	// --------------------------------------------------

	p.HP -= damage

	// --------------------------------------------------
	// INVINCIBILITE TEMPORAIRE
	// --------------------------------------------------

	p.InvincibleTimer = 45

	// --------------------------------------------------
	// FLASH
	// --------------------------------------------------

	p.FlashTimer = 7

	// --------------------------------------------------
	// ANIMATION DE DEGATS
	// --------------------------------------------------

	p.HurtTimer = p.HurtDuration

	// --------------------------------------------------
	// ANNULE L'ATTAQUE EN COURS
	// --------------------------------------------------

	p.Attacking = false

	p.AttackTimer = 0

	p.CurrentAttackFrame = 0

	p.AttackHit = false

	// --------------------------------------------------
	// RECUL
	// --------------------------------------------------

	playerCenter := p.X + p.Width/2

	if playerCenter < sourceX {
		p.KnockbackX = -3.2

		p.FacingLeft = false
	} else {
		p.KnockbackX = 3.2

		p.FacingLeft = true
	}

	// Petit saut au moment du choc.
	p.VelocityY = -2.3

	// --------------------------------------------------
	// MORT
	// --------------------------------------------------

	if p.HP <= 0 {
		p.HP = 0

		p.Alive = false
	}

	return true
}

// --------------------------------------------------
// IMAGE ACTUELLE
// --------------------------------------------------

func (p *Player) CurrentImage() *ebiten.Image {
	if p.Attacking {
		return p.AttackFrames[p.CurrentAttackFrame]
	}

	return p.WalkFrames[p.CurrentWalkFrame]
}

// --------------------------------------------------
// DRAW
// --------------------------------------------------

func (p *Player) Draw(
	screen *ebiten.Image,
) {
	if !p.Alive {
		return
	}

	img := p.CurrentImage()

	frameWidth := float64(
		img.Bounds().Dx(),
	) * p.Scale

	frameHeight := float64(
		img.Bounds().Dy(),
	) * p.Scale

	centerX := p.X + p.Width/2

	playerBottom := p.Y + p.Height

	drawX := centerX - frameWidth/2

	drawY := playerBottom - frameHeight

	// ==================================================
	// PETIT TREMBLEMENT VISUEL QUAND IL EST TOUCHE
	// ==================================================

	if p.HurtTimer > 0 {
		if p.HurtTimer%4 < 2 {
			drawX -= 1.5
		} else {
			drawX += 1.5
		}
	}

	options := &ebiten.DrawImageOptions{}

	// ==================================================
	// ORIENTATION
	// ==================================================

	if p.FacingLeft {
		options.GeoM.Scale(
			-p.Scale,
			p.Scale,
		)

		options.GeoM.Translate(
			drawX+frameWidth,
			drawY,
		)
	} else {
		options.GeoM.Scale(
			p.Scale,
			p.Scale,
		)

		options.GeoM.Translate(
			drawX,
			drawY,
		)
	}

	// ==================================================
	// FLASH DEGATS
	// ==================================================

	if p.FlashTimer > 0 {
		options.ColorScale.Scale(
			2,
			2,
			2,
			1,
		)
	}

	screen.DrawImage(
		img,
		options,
	)
}
