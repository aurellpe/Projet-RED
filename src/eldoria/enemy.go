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
)

//go:embed assets/enemy/enemy1.png
var enemy1Data []byte

//go:embed assets/enemy/enemy2.png
var enemy2Data []byte

//go:embed assets/enemy/enemy3.png
var enemy3Data []byte

//go:embed assets/enemy/enemy4.png
var enemy4Data []byte

const (
	EnemyTypeNormal = iota
	EnemyTypeFast
	EnemyTypeHeavy
	EnemyTypeRanged
)

type EnemyDeathParticle struct {
	X       float64
	Y       float64
	VX      float64
	VY      float64
	Life    int
	MaxLife int
}

type Enemy struct {
	X float64
	Y float64

	Width  float64
	Height float64

	Scale float64
	Speed float64

	SpriteYOffset float64

	Type int

	HP    int
	MaxHP int

	Damage int

	Alive bool

	Dying         bool
	DeathTimer    int
	DeathDuration int
	DeathRewarded bool

	Frames      []*ebiten.Image
	FrameBounds []image.Rectangle

	CurrentFrame int
	FrameTimer   int
	FrameSpeed   int

	FacingLeft bool
	Moving     bool

	KnockbackX float64

	FlashTimer int

	Attacking      bool
	AttackTimer    int
	AttackCooldown int
	AttackHit      bool

	AttackRange int

	ProjectileActive  bool
	ProjectileSpawned bool
	ProjectileX       float64
	ProjectileY       float64
	ProjectileVX      float64

	DeathParticles []EnemyDeathParticle

	VisualHeight float64
}

func findEnemyContentBounds(img image.Image) image.Rectangle {
	bounds := img.Bounds()

	minX := bounds.Max.X
	minY := bounds.Max.Y
	maxX := bounds.Min.X
	maxY := bounds.Min.Y

	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()

			if alpha <= 1000 {
				continue
			}

			found = true

			if x < minX {
				minX = x
			}

			if y < minY {
				minY = y
			}

			if x > maxX {
				maxX = x
			}

			if y > maxY {
				maxY = y
			}
		}
	}

	if !found {
		return bounds
	}

	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func loadEnemyFrame(data []byte) (*ebiten.Image, image.Rectangle) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	contentBounds := findEnemyContentBounds(img)

	return ebiten.NewImageFromImage(img), contentBounds
}

func newEnemyWithType(x float64, enemyType int) *Enemy {
	enemy1, bounds1 := loadEnemyFrame(enemy1Data)
	enemy2, bounds2 := loadEnemyFrame(enemy2Data)
	enemy3, bounds3 := loadEnemyFrame(enemy3Data)
	enemy4, bounds4 := loadEnemyFrame(enemy4Data)

	frames := []*ebiten.Image{
		enemy1,
		enemy2,
		enemy3,
		enemy4,
	}

	frameBounds := []image.Rectangle{
		bounds1,
		bounds2,
		bounds3,
		bounds4,
	}

	scale := 0.055
	width := 28.0
	height := 55.0
	speed := 0.45
	hp := 100
	damage := 10
	frameSpeed := 7
	attackRange := 32

	switch enemyType {
	case EnemyTypeFast:
		speed = 0.72
		hp = 70
		damage = 8
		frameSpeed = 5

	case EnemyTypeHeavy:
		scale = 0.065
		width = 34
		height = 63
		speed = 0.27
		hp = 180
		damage = 20
		frameSpeed = 9
		attackRange = 38

	case EnemyTypeRanged:
		speed = 0.34
		hp = 80
		damage = 12
		frameSpeed = 7
		attackRange = 100
	}

	visualHeight := float64(bounds1.Dy()) * scale

	return &Enemy{
		X: x,
		Y: groundY - height,

		Width:  width,
		Height: height,

		Scale: scale,
		Speed: speed,

		SpriteYOffset: -22,

		Type: enemyType,

		HP:    hp,
		MaxHP: hp,

		Damage: damage,

		Alive: true,

		Dying:         false,
		DeathTimer:    0,
		DeathDuration: 40,
		DeathRewarded: false,

		Frames:      frames,
		FrameBounds: frameBounds,

		CurrentFrame: 0,
		FrameTimer:   0,
		FrameSpeed:   frameSpeed,

		FacingLeft: true,
		Moving:     false,

		KnockbackX: 0,

		FlashTimer: 0,

		Attacking:      false,
		AttackTimer:    0,
		AttackCooldown: 0,
		AttackHit:      false,

		AttackRange: attackRange,

		ProjectileActive:  false,
		ProjectileSpawned: false,

		DeathParticles: []EnemyDeathParticle{},

		VisualHeight: visualHeight,
	}
}

func NewEnemy(x float64) *Enemy {
	return newEnemyWithType(x, EnemyTypeNormal)
}

func NewFastEnemy(x float64) *Enemy {
	return newEnemyWithType(x, EnemyTypeFast)
}

func NewHeavyEnemy(x float64) *Enemy {
	return newEnemyWithType(x, EnemyTypeHeavy)
}

func NewRangedEnemy(x float64) *Enemy {
	return newEnemyWithType(x, EnemyTypeRanged)
}

func (e *Enemy) Update(player *Player) int {
	e.UpdateDeathParticles()

	if !e.Alive {
		return 0
	}

	if e.ProjectileActive {
		damage := e.UpdateProjectile(player)

		if damage > 0 {
			return damage
		}
	}

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

	e.X += e.KnockbackX
	e.KnockbackX *= 0.82

	if e.KnockbackX > -0.05 && e.KnockbackX < 0.05 {
		e.KnockbackX = 0
	}

	if player == nil || !player.Alive {
		e.Moving = false
		e.Attacking = false
		e.CurrentFrame = 0
		return 0
	}

	if e.Attacking {
		e.Moving = false
		return e.UpdateAttack(player)
	}

	playerCenter := player.X + player.Width/2
	enemyCenter := e.X + e.Width/2

	distance := playerCenter - enemyCenter
	absoluteDistance := math.Abs(distance)

	if distance < 0 {
		e.FacingLeft = true
	} else {
		e.FacingLeft = false
	}

	if e.Type == EnemyTypeRanged {
		e.UpdateRangedMovement(distance, absoluteDistance)

		if absoluteDistance >= 65 && absoluteDistance <= 125 && e.AttackCooldown <= 0 {
			e.StartAttack()
		}

		e.UpdateWalkAnimation()
		e.KeepInsideScreen()

		return 0
	}

	e.Moving = false

	if absoluteDistance > float64(e.AttackRange) {
		e.Moving = true

		if distance > 0 {
			e.X += e.Speed
		} else {
			e.X -= e.Speed
		}
	}

	e.UpdateWalkAnimation()

	if absoluteDistance <= float64(e.AttackRange) && e.AttackCooldown <= 0 {
		e.StartAttack()
	}

	e.KeepInsideScreen()

	return 0
}

func (e *Enemy) UpdateRangedMovement(distance float64, absoluteDistance float64) {
	e.Moving = false

	if absoluteDistance > 125 {
		e.Moving = true

		if distance > 0 {
			e.X += e.Speed
		} else {
			e.X -= e.Speed
		}

		return
	}

	if absoluteDistance < 60 {
		e.Moving = true

		if distance > 0 {
			e.X -= e.Speed
		} else {
			e.X += e.Speed
		}
	}
}

func (e *Enemy) StartAttack() {
	e.Attacking = true
	e.AttackTimer = 0
	e.AttackHit = false
	e.ProjectileSpawned = false
	e.CurrentFrame = 3
}

func (e *Enemy) UpdateAttack(player *Player) int {
	e.AttackTimer++
	e.CurrentFrame = 3

	if e.Type == EnemyTypeRanged {
		if e.AttackTimer == 8 && !e.ProjectileSpawned {
			e.SpawnProjectile(player)
		}

		if e.AttackTimer > 20 {
			e.Attacking = false
			e.AttackTimer = 0
			e.AttackHit = false
			e.ProjectileSpawned = false
			e.AttackCooldown = 65
			e.CurrentFrame = 0
		}

		return 0
	}

	if e.AttackTimer >= 7 && e.AttackTimer <= 12 && !e.AttackHit {
		attackBox := e.AttackBox()
		playerBox := player.HitBox()

		if Intersects(attackBox, playerBox) {
			e.AttackHit = true

			enemyCenter := e.X + e.Width/2

			if player.Hit(e.Damage, enemyCenter) {
				return e.Damage
			}
		}
	}

	endTimer := 18
	cooldown := 45

	if e.Type == EnemyTypeFast {
		endTimer = 15
		cooldown = 30
	}

	if e.Type == EnemyTypeHeavy {
		endTimer = 24
		cooldown = 65
	}

	if e.AttackTimer > endTimer {
		e.Attacking = false
		e.AttackTimer = 0
		e.AttackHit = false
		e.AttackCooldown = cooldown
		e.CurrentFrame = 0
	}

	return 0
}

func (e *Enemy) SpawnProjectile(player *Player) {
	e.ProjectileSpawned = true
	e.ProjectileActive = true

	e.ProjectileX = e.X + e.Width/2
	e.ProjectileY = groundY - e.Height*0.60

	playerCenter := player.X + player.Width/2

	if playerCenter < e.ProjectileX {
		e.ProjectileVX = -1.8
	} else {
		e.ProjectileVX = 1.8
	}
}

func (e *Enemy) UpdateProjectile(player *Player) int {
	e.ProjectileX += e.ProjectileVX

	if e.ProjectileX < -10 || e.ProjectileX > screenWidth+10 {
		e.ProjectileActive = false
		return 0
	}

	if player == nil || !player.Alive {
		return 0
	}

	projectileBox := Rect{
		X: e.ProjectileX - 2,
		Y: e.ProjectileY - 2,
		W: 4,
		H: 4,
	}

	if !Intersects(projectileBox, player.HitBox()) {
		return 0
	}

	e.ProjectileActive = false

	if player.Hit(e.Damage, e.X+e.Width/2) {
		return e.Damage
	}

	return 0
}

func (e *Enemy) UpdateWalkAnimation() {
	if !e.Moving {
		e.CurrentFrame = 0
		e.FrameTimer = 0
		return
	}

	e.FrameTimer++

	if e.FrameTimer < e.FrameSpeed {
		return
	}

	e.FrameTimer = 0

	if e.CurrentFrame == 1 {
		e.CurrentFrame = 2
	} else {
		e.CurrentFrame = 1
	}
}

func (e *Enemy) KeepInsideScreen() {
	if e.X < 0 {
		e.X = 0
	}

	if e.X+e.Width > screenWidth {
		e.X = screenWidth - e.Width
	}

	e.Y = groundY - e.Height
}

func (e *Enemy) HitBox() Rect {
	return Rect{
		X: e.X,
		Y: e.Y,
		W: e.Width,
		H: e.Height,
	}
}

func (e *Enemy) AttackBox() Rect {
	attackWidth := 32.0
	attackHeight := e.Height

	if e.Type == EnemyTypeHeavy {
		attackWidth = 40
	}

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

func (e *Enemy) Hit(damage int, attackerX float64) bool {
	if !e.Alive || e.Dying {
		return false
	}

	e.HP -= damage
	e.FlashTimer = 6

	enemyCenter := e.X + e.Width/2

	if enemyCenter < attackerX {
		e.KnockbackX = -1.8
	} else {
		e.KnockbackX = 1.8
	}

	if e.Type == EnemyTypeHeavy {
		e.KnockbackX *= 0.45
	}

	if e.HP <= 0 {
		e.HP = 0
		e.Dying = true
		e.DeathTimer = 0

		e.Attacking = false
		e.Moving = false
		e.ProjectileActive = false

		e.CurrentFrame = 0

		e.SpawnDeathParticles()
	}

	return true
}

func (e *Enemy) SpawnDeathParticles() {
	centerX := e.X + e.Width/2
	centerY := e.Y + e.Height/2

	for i := 0; i < 16; i++ {
		vx := float64((i%5)-2) * 0.18
		vy := -0.25 - float64(i%4)*0.10

		particle := EnemyDeathParticle{
			X:       centerX + float64((i%4)-2)*2,
			Y:       centerY + float64(i%5)*2,
			VX:      vx,
			VY:      vy,
			Life:    30 + i%8,
			MaxLife: 30 + i%8,
		}

		e.DeathParticles = append(e.DeathParticles, particle)
	}
}

func (e *Enemy) UpdateDeathParticles() {
	active := []EnemyDeathParticle{}

	for _, particle := range e.DeathParticles {
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.VY += 0.015
		particle.Life--

		if particle.Life > 0 {
			active = append(active, particle)
		}
	}

	e.DeathParticles = active
}

func (e *Enemy) UpdateDeath() {
	e.DeathTimer++

	e.Attacking = false
	e.Moving = false
	e.CurrentFrame = 0

	e.X += e.KnockbackX
	e.KnockbackX *= 0.90

	e.KeepInsideScreen()

	if e.DeathTimer >= e.DeathDuration {
		e.Alive = false
	}
}

func (e *Enemy) DrawDeathParticles(screen *ebiten.Image) {
	for _, particle := range e.DeathParticles {
		ratio := float64(particle.Life) / float64(particle.MaxLife)

		alpha := uint8(200 * ratio)

		particleColor := color.RGBA{
			R: 100,
			G: 70,
			B: 140,
			A: alpha,
		}

		if e.Type == EnemyTypeHeavy {
			particleColor = color.RGBA{
				R: 160,
				G: 55,
				B: 55,
				A: alpha,
			}
		}

		if e.Type == EnemyTypeFast {
			particleColor = color.RGBA{
				R: 70,
				G: 150,
				B: 190,
				A: alpha,
			}
		}

		if e.Type == EnemyTypeRanged {
			particleColor = color.RGBA{
				R: 150,
				G: 70,
				B: 210,
				A: alpha,
			}
		}

		ebitenutil.DrawRect(screen, particle.X, particle.Y, 2, 2, particleColor)
	}
}

func (e *Enemy) DrawProjectile(screen *ebiten.Image) {
	if !e.ProjectileActive {
		return
	}

	ebitenutil.DrawRect(
		screen,
		e.ProjectileX-3,
		e.ProjectileY-3,
		6,
		6,
		color.RGBA{
			R: 100,
			G: 40,
			B: 180,
			A: 80,
		},
	)

	ebitenutil.DrawRect(
		screen,
		e.ProjectileX-1,
		e.ProjectileY-1,
		3,
		3,
		color.RGBA{
			R: 210,
			G: 120,
			B: 255,
			A: 255,
		},
	)
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.DrawDeathParticles(screen)
	e.DrawProjectile(screen)

	if !e.Alive {
		return
	}

	frameIndex := e.CurrentFrame

	if e.Dying {
		frameIndex = 0
	}

	img := e.Frames[frameIndex]
	contentBounds := e.FrameBounds[frameIndex]

	contentHeight := float64(contentBounds.Dy())

	if contentHeight <= 0 {
		return
	}

	frameScale := e.VisualHeight / contentHeight

	anchorX := float64(contentBounds.Min.X+contentBounds.Max.X) / 2
	anchorY := float64(contentBounds.Max.Y)

	centerX := e.X + e.Width/2
	bottomY := groundY + e.SpriteYOffset

	deathAlpha := float32(1)

	if e.Dying {
		progress := float64(e.DeathTimer) / float64(e.DeathDuration)

		if progress > 1 {
			progress = 1
		}

		bottomY += progress * 10
		deathAlpha = float32(1 - progress)
	}

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Translate(-anchorX, -anchorY)

	if e.FacingLeft {
		options.GeoM.Scale(-frameScale, frameScale)
	} else {
		options.GeoM.Scale(frameScale, frameScale)
	}

	options.GeoM.Translate(centerX, bottomY)

	switch e.Type {
	case EnemyTypeFast:
		options.ColorScale.Scale(0.85, 1, 1.15, 1)

	case EnemyTypeHeavy:
		options.ColorScale.Scale(1.15, 0.78, 0.78, 1)

	case EnemyTypeRanged:
		options.ColorScale.Scale(1.05, 0.80, 1.20, 1)
	}

	if e.FlashTimer > 0 {
		options.ColorScale.Scale(2, 2, 2, 1)
	}

	if e.Dying {
		options.ColorScale.Scale(1, 1, 1, deathAlpha)
	}

	screen.DrawImage(img, options)

	if e.Dying {
		return
	}

	barY := bottomY - e.VisualHeight - 5

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

	ratio := float64(e.HP) / float64(e.MaxHP)

	if ratio < 0 {
		ratio = 0
	}

	if ratio > 1 {
		ratio = 1
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
