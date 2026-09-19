package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/boss/boss1.png
var boss1Data []byte

//go:embed assets/boss/boss2.png
var boss2Data []byte

//go:embed assets/boss/boss3.png
var boss3Data []byte

//go:embed assets/boss/boss4.png
var boss4Data []byte

//go:embed assets/boss/walk1.png
var bossWalk1Data []byte

//go:embed assets/boss/walk2.png
var bossWalk2Data []byte

//go:embed assets/boss/walk3.png
var bossWalk3Data []byte

//go:embed assets/boss/walk4.png
var bossWalk4Data []byte

type Boss struct {
	X float64
	Y float64

	Scale         float64
	SpriteYOffset float64

	HitWidth  float64
	HitHeight float64

	HP    int
	MaxHP int

	Alive bool

	Dying         bool
	DeathTimer    int
	DeathDuration int

	AttackFrames []*ebiten.Image
	WalkFrames   []*ebiten.Image

	CurrentAttackFrame int
	CurrentWalkFrame   int
	WalkFrameTimer     int
	WalkFrameSpeed     int

	FacingLeft bool
	Moving     bool

	Speed      float64
	KnockbackX float64

	Attacking      bool
	AttackTimer    int
	AttackCooldown int
	AttackHit      bool
	FlashTimer     int
	AttackRange    float64
}

func loadBossImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func NewBoss() *Boss {
	attackFrames := []*ebiten.Image{
		loadBossImage(boss1Data),
		loadBossImage(boss2Data),
		loadBossImage(boss3Data),
		loadBossImage(boss4Data),
	}

	walkFrames := []*ebiten.Image{
		loadBossImage(bossWalk1Data),
		loadBossImage(bossWalk2Data),
		loadBossImage(bossWalk3Data),
		loadBossImage(bossWalk4Data),
	}

	return &Boss{
		X:             230,
		Y:             groundY - 72,
		Scale:         0.075,
		SpriteYOffset: 0,

		HitWidth:  45,
		HitHeight: 72,

		HP:    500,
		MaxHP: 500,

		Alive: true,

		Dying:         false,
		DeathTimer:    0,
		DeathDuration: 70,

		AttackFrames: attackFrames,
		WalkFrames:   walkFrames,

		CurrentAttackFrame: 0,
		CurrentWalkFrame:   0,
		WalkFrameTimer:     0,
		WalkFrameSpeed:     8,

		FacingLeft: true,
		Moving:     false,

		Speed:      0.30,
		KnockbackX: 0,

		Attacking:      false,
		AttackTimer:    0,
		AttackCooldown: 0,
		AttackHit:      false,
		FlashTimer:     0,
		AttackRange:    48,
	}
}

func (b *Boss) Update(player *Player) int {
	if !b.Alive {
		return 0
	}

	if b.Dying {
		b.UpdateDeath()
		return 0
	}

	if b.FlashTimer > 0 {
		b.FlashTimer--
	}

	if b.AttackCooldown > 0 {
		b.AttackCooldown--
	}

	b.X += b.KnockbackX
	b.KnockbackX *= 0.84

	if b.KnockbackX > -0.05 && b.KnockbackX < 0.05 {
		b.KnockbackX = 0
	}

	b.KeepInsideScreen()

	if !player.Alive {
		b.Moving = false
		b.Attacking = false
		b.AttackTimer = 0
		b.CurrentAttackFrame = 0
		b.CurrentWalkFrame = 0
		return 0
	}

	if b.Attacking {
		b.Moving = false
		return b.UpdateAttack(player)
	}

	playerCenter := player.X + player.Width/2
	bossCenter := b.X + b.HitWidth/2
	distance := playerCenter - bossCenter

	absoluteDistance := distance
	if absoluteDistance < 0 {
		absoluteDistance = -absoluteDistance
	}

	if distance < 0 {
		b.FacingLeft = true
	} else {
		b.FacingLeft = false
	}

	if absoluteDistance > b.AttackRange {
		b.Moving = true

		if distance > 0 {
			b.X += b.Speed
		} else {
			b.X -= b.Speed
		}

		b.WalkFrameTimer++

		if b.WalkFrameTimer >= b.WalkFrameSpeed {
			b.WalkFrameTimer = 0
			b.CurrentWalkFrame++

			if b.CurrentWalkFrame >= len(b.WalkFrames) {
				b.CurrentWalkFrame = 0
			}
		}

		b.KeepInsideScreen()
		return 0
	}

	b.Moving = false
	b.WalkFrameTimer = 0
	b.CurrentWalkFrame = 0

	if b.AttackCooldown > 0 {
		b.CurrentAttackFrame = 0
		return 0
	}

	b.Attacking = true
	b.AttackTimer = 0
	b.AttackHit = false
	b.CurrentAttackFrame = 1

	return 0
}

func (b *Boss) UpdateDeath() {
	b.DeathTimer++

	b.Moving = false
	b.Attacking = false
	b.AttackTimer = 0
	b.AttackHit = false
	b.CurrentAttackFrame = 0
	b.CurrentWalkFrame = 0

	b.X += b.KnockbackX
	b.KnockbackX *= 0.92

	b.KeepInsideScreen()

	if b.DeathTimer >= b.DeathDuration {
		b.Alive = false
	}
}

func (b *Boss) UpdateAttack(player *Player) int {
	b.AttackTimer++

	if b.AttackTimer <= 14 {
		b.CurrentAttackFrame = 1
		return 0
	}

	if b.AttackTimer <= 22 {
		b.CurrentAttackFrame = 2

		if !b.AttackHit {
			attackBox := b.AttackBox()
			playerBox := player.HitBox()

			if Intersects(attackBox, playerBox) {
				b.AttackHit = true
				bossCenter := b.X + b.HitWidth/2

				if player.Hit(25, bossCenter) {
					return 25
				}
			}
		}

		return 0
	}

	if b.AttackTimer <= 36 {
		b.CurrentAttackFrame = 3
		return 0
	}

	b.Attacking = false
	b.AttackTimer = 0
	b.AttackHit = false
	b.CurrentAttackFrame = 0
	b.AttackCooldown = 80

	return 0
}

func (b *Boss) KeepInsideScreen() {
	if b.X < 0 {
		b.X = 0
	}

	if b.X+b.HitWidth > screenWidth {
		b.X = screenWidth - b.HitWidth
	}

	b.Y = groundY - b.HitHeight
}

func (b *Boss) HitBox() Rect {
	return Rect{
		X: b.X,
		Y: groundY - b.HitHeight,
		W: b.HitWidth,
		H: b.HitHeight,
	}
}

func (b *Boss) AttackBox() Rect {
	attackWidth := 58.0
	attackHeight := 72.0

	if b.FacingLeft {
		return Rect{
			X: b.X - attackWidth,
			Y: groundY - attackHeight,
			W: attackWidth,
			H: attackHeight,
		}
	}

	return Rect{
		X: b.X + b.HitWidth,
		Y: groundY - attackHeight,
		W: attackWidth,
		H: attackHeight,
	}
}

func (b *Boss) Hit(damage int, attackerX float64) bool {
	if !b.Alive || b.Dying {
		return false
	}

	b.HP -= damage
	b.FlashTimer = 6

	bossCenter := b.X + b.HitWidth/2

	if bossCenter < attackerX {
		b.KnockbackX = -2
	} else {
		b.KnockbackX = 2
	}

	if b.HP <= 0 {
		b.HP = 0
		b.Dying = true
		b.DeathTimer = 0
		b.Attacking = false
		b.Moving = false
		b.CurrentAttackFrame = 0
		b.CurrentWalkFrame = 0
	}

	return true
}

func (b *Boss) Draw(screen *ebiten.Image) {
	if !b.Alive {
		return
	}

	var img *ebiten.Image

	if b.Dying {
		img = b.AttackFrames[0]
	} else if b.Attacking {
		img = b.AttackFrames[b.CurrentAttackFrame]
	} else if b.Moving {
		img = b.WalkFrames[b.CurrentWalkFrame]
	} else {
		img = b.AttackFrames[0]
	}

	frameWidth := float64(img.Bounds().Dx()) * b.Scale
	frameHeight := float64(img.Bounds().Dy()) * b.Scale

	centerX := b.X + b.HitWidth/2
	drawX := centerX - frameWidth/2
	drawY := groundY - frameHeight + b.SpriteYOffset

	deathAlpha := float32(1)

	if b.Dying {
		progress := float64(b.DeathTimer) / float64(b.DeathDuration)

		if progress > 1 {
			progress = 1
		}

		drawY += progress * 25
		deathAlpha = float32(1 - progress)
	}

	options := &ebiten.DrawImageOptions{}

	if b.Moving && !b.Attacking && !b.Dying {
		if b.FacingLeft {
			options.GeoM.Scale(b.Scale, b.Scale)
			options.GeoM.Translate(drawX, drawY)
		} else {
			options.GeoM.Scale(-b.Scale, b.Scale)
			options.GeoM.Translate(drawX+frameWidth, drawY)
		}
	} else {
		if b.FacingLeft {
			options.GeoM.Scale(-b.Scale, b.Scale)
			options.GeoM.Translate(drawX+frameWidth, drawY)
		} else {
			options.GeoM.Scale(b.Scale, b.Scale)
			options.GeoM.Translate(drawX, drawY)
		}
	}

	if b.FlashTimer > 0 {
		options.ColorScale.Scale(2, 2, 2, 1)
	}

	if b.Dying {
		options.ColorScale.Scale(1, 1, 1, deathAlpha)
	}

	screen.DrawImage(img, options)
}
