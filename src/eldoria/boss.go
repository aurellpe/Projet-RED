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

	Scale float64

	HitWidth  float64
	HitHeight float64

	HP    int
	MaxHP int

	Alive bool

	Dying         bool
	DeathTimer    int
	DeathDuration int

	AttackFrames      []*ebiten.Image
	AttackFrameBounds []image.Rectangle
	AttackFootY       []float64

	WalkFrames      []*ebiten.Image
	WalkFrameBounds []image.Rectangle
	WalkFootY       []float64

	CurrentAttackFrame int
	CurrentWalkFrame   int

	WalkFrameTimer   int
	WalkFrameSpeed   int
	WalkSequenceStep int

	FacingLeft bool
	Moving     bool

	Speed      float64
	KnockbackX float64

	Attacking       bool
	AttackTimer     int
	AttackCooldown  int
	AttackHit       bool
	AttackLungeDone bool

	FlashTimer  int
	AttackRange float64

	Phase int

	PhaseTransition      bool
	PhaseTransitionTimer int
	PhaseTransitionMax   int

	CombatTimer int
	AuraTimer   int

	VisualHeight       float64
	VisualGroundOffset float64
}

func findBossContentBounds(img image.Image) image.Rectangle {
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

func findBossFootY(img image.Image, content image.Rectangle) float64 {
	if content.Empty() {
		return float64(img.Bounds().Max.Y)
	}

	contentWidth := content.Dx()

	left := content.Min.X + contentWidth*25/100
	right := content.Min.X + contentWidth*75/100

	if left >= right {
		left = content.Min.X
		right = content.Max.X
	}

	searchStartY := content.Min.Y + content.Dy()/2

	footY := content.Min.Y
	found := false

	for y := searchStartY; y < content.Max.Y; y++ {
		for x := left; x < right; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()

			if alpha <= 1000 {
				continue
			}

			if y+1 > footY {
				footY = y + 1
			}

			found = true
		}
	}

	if !found {
		return float64(content.Max.Y)
	}

	return float64(footY)
}

func loadBossFrame(data []byte) (*ebiten.Image, image.Rectangle, float64) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	contentBounds := findBossContentBounds(img)
	footY := findBossFootY(img, contentBounds)

	return ebiten.NewImageFromImage(img), contentBounds, footY
}

func NewBoss() *Boss {
	attack1, attackBounds1, attackFoot1 := loadBossFrame(boss1Data)
	attack2, attackBounds2, attackFoot2 := loadBossFrame(boss2Data)
	attack3, attackBounds3, attackFoot3 := loadBossFrame(boss3Data)
	attack4, attackBounds4, attackFoot4 := loadBossFrame(boss4Data)

	walk1, walkBounds1, walkFoot1 := loadBossFrame(bossWalk1Data)
	walk2, walkBounds2, walkFoot2 := loadBossFrame(bossWalk2Data)
	walk3, walkBounds3, walkFoot3 := loadBossFrame(bossWalk3Data)
	walk4, walkBounds4, walkFoot4 := loadBossFrame(bossWalk4Data)

	attackFrames := []*ebiten.Image{
		attack1,
		attack2,
		attack3,
		attack4,
	}

	attackBounds := []image.Rectangle{
		attackBounds1,
		attackBounds2,
		attackBounds3,
		attackBounds4,
	}

	attackFeet := []float64{
		attackFoot1,
		attackFoot2,
		attackFoot3,
		attackFoot4,
	}

	walkFrames := []*ebiten.Image{
		walk1,
		walk2,
		walk3,
		walk4,
	}

	walkBounds := []image.Rectangle{
		walkBounds1,
		walkBounds2,
		walkBounds3,
		walkBounds4,
	}

	walkFeet := []float64{
		walkFoot1,
		walkFoot2,
		walkFoot3,
		walkFoot4,
	}

	scale := 0.11

	visualHeight := float64(walkBounds1.Dy()) * scale

	bottomPadding := float64(walk1.Bounds().Max.Y) - walkFoot1
	visualGroundOffset := bottomPadding * scale

	if visualGroundOffset < 0 {
		visualGroundOffset = 0
	}

	return &Boss{
		X: 225,
		Y: groundY - 92,

		Scale: scale,

		HitWidth:  56,
		HitHeight: 92,

		HP:    500,
		MaxHP: 500,

		Alive: true,

		Dying:         false,
		DeathTimer:    0,
		DeathDuration: 70,

		AttackFrames:      attackFrames,
		AttackFrameBounds: attackBounds,
		AttackFootY:       attackFeet,

		WalkFrames:      walkFrames,
		WalkFrameBounds: walkBounds,
		WalkFootY:       walkFeet,

		CurrentAttackFrame: 0,
		CurrentWalkFrame:   0,

		WalkFrameTimer:   0,
		WalkFrameSpeed:   8,
		WalkSequenceStep: 0,

		FacingLeft: true,
		Moving:     false,

		Speed:      0.30,
		KnockbackX: 0,

		Attacking:       false,
		AttackTimer:     0,
		AttackCooldown:  0,
		AttackHit:       false,
		AttackLungeDone: false,

		FlashTimer:  0,
		AttackRange: 60,

		Phase: 1,

		PhaseTransition:      false,
		PhaseTransitionTimer: 0,
		PhaseTransitionMax:   55,

		CombatTimer: 0,
		AuraTimer:   0,

		VisualHeight:       visualHeight,
		VisualGroundOffset: visualGroundOffset,
	}
}

func (b *Boss) Update(player *Player) int {
	if !b.Alive {
		return 0
	}

	b.CombatTimer++
	b.AuraTimer++

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

	if b.Phase == 1 && b.HP <= b.MaxHP/2 {
		b.StartPhaseTwo()
	}

	if b.PhaseTransition {
		b.UpdatePhaseTransition()
		return 0
	}

	b.X += b.KnockbackX
	b.KnockbackX *= 0.84

	if b.KnockbackX > -0.05 && b.KnockbackX < 0.05 {
		b.KnockbackX = 0
	}

	b.KeepInsideScreen()

	if player == nil {
		return 0
	}

	if !player.Alive {
		b.Moving = false
		b.Attacking = false
		b.AttackTimer = 0
		b.CurrentAttackFrame = 0
		b.CurrentWalkFrame = 0
		b.WalkFrameTimer = 0
		b.WalkSequenceStep = 0
		return 0
	}

	if b.Attacking {
		b.Moving = false
		return b.UpdateAttack(player)
	}

	playerCenter := player.X + player.Width/2
	bossCenter := b.X + b.HitWidth/2

	distance := playerCenter - bossCenter
	absoluteDistance := math.Abs(distance)

	if distance < 0 {
		b.FacingLeft = true
	} else {
		b.FacingLeft = false
	}

	moveSpeed := b.GetMoveSpeed()
	attackRange := b.GetAttackRange()

	if absoluteDistance > attackRange {
		b.Moving = true

		if distance > 0 {
			b.X += moveSpeed
		} else {
			b.X -= moveSpeed
		}

		b.UpdateWalkAnimation()
		b.KeepInsideScreen()

		return 0
	}

	b.Moving = false
	b.WalkFrameTimer = 0
	b.CurrentWalkFrame = 0
	b.WalkSequenceStep = 0

	if b.AttackCooldown > 0 {
		b.CurrentAttackFrame = 0
		return 0
	}

	b.StartAttack()

	return 0
}

func (b *Boss) StartPhaseTwo() {
	if b.Phase >= 2 {
		return
	}

	b.Phase = 2

	b.PhaseTransition = true
	b.PhaseTransitionTimer = 0

	b.Attacking = false
	b.AttackTimer = 0
	b.AttackHit = false
	b.AttackLungeDone = false

	b.Moving = false

	b.CurrentAttackFrame = 0
	b.CurrentWalkFrame = 0

	b.WalkFrameTimer = 0
	b.WalkSequenceStep = 0

	b.KnockbackX = 0
	b.AttackCooldown = 20
}

func (b *Boss) UpdatePhaseTransition() {
	b.PhaseTransitionTimer++

	b.Moving = false
	b.Attacking = false

	b.CurrentAttackFrame = 0
	b.CurrentWalkFrame = 0

	if b.PhaseTransitionTimer >= b.PhaseTransitionMax {
		b.PhaseTransition = false
		b.PhaseTransitionTimer = 0
		b.AttackCooldown = 25
	}
}

func (b *Boss) GetMoveSpeed() float64 {
	if b.Phase >= 2 {
		return 0.48
	}

	return b.Speed
}

func (b *Boss) GetAttackRange() float64 {
	if b.Phase >= 2 {
		return 66
	}

	return b.AttackRange
}

func (b *Boss) GetAttackDamage() int {
	if b.Phase >= 2 {
		return 35
	}

	return 25
}

func (b *Boss) GetAttackCooldown() int {
	if b.Phase >= 2 {
		return 48
	}

	return 80
}

func (b *Boss) GetWalkFrameSpeed() int {
	if b.Phase >= 2 {
		return 5
	}

	return b.WalkFrameSpeed
}

func (b *Boss) UpdateWalkAnimation() {
	b.WalkFrameTimer++

	if b.WalkFrameTimer < b.GetWalkFrameSpeed() {
		return
	}

	b.WalkFrameTimer = 0

	walkSequence := []int{
		0,
		1,
		2,
		3,
		2,
		1,
	}

	b.WalkSequenceStep++

	if b.WalkSequenceStep >= len(walkSequence) {
		b.WalkSequenceStep = 0
	}

	b.CurrentWalkFrame = walkSequence[b.WalkSequenceStep]
}

func (b *Boss) StartAttack() {
	b.Attacking = true

	b.AttackTimer = 0
	b.AttackHit = false
	b.AttackLungeDone = false

	b.CurrentAttackFrame = 1

	b.CurrentWalkFrame = 0
	b.WalkFrameTimer = 0
	b.WalkSequenceStep = 0
}

func (b *Boss) UpdateAttack(player *Player) int {
	if b.Phase >= 2 {
		return b.UpdatePhaseTwoAttack(player)
	}

	return b.UpdatePhaseOneAttack(player)
}

func (b *Boss) UpdatePhaseOneAttack(player *Player) int {
	b.AttackTimer++

	if b.AttackTimer <= 14 {
		b.CurrentAttackFrame = 1
		return 0
	}

	if b.AttackTimer <= 22 {
		b.CurrentAttackFrame = 2

		if !b.AttackHit {
			return b.TryDamagePlayer(player)
		}

		return 0
	}

	if b.AttackTimer <= 36 {
		b.CurrentAttackFrame = 3
		return 0
	}

	b.EndAttack()

	return 0
}

func (b *Boss) UpdatePhaseTwoAttack(player *Player) int {
	b.AttackTimer++

	if b.AttackTimer <= 9 {
		b.CurrentAttackFrame = 1
		return 0
	}

	if b.AttackTimer <= 17 {
		b.CurrentAttackFrame = 2

		if !b.AttackLungeDone {
			b.PerformAttackLunge()
			b.AttackLungeDone = true
		}

		if !b.AttackHit {
			return b.TryDamagePlayer(player)
		}

		return 0
	}

	if b.AttackTimer <= 28 {
		b.CurrentAttackFrame = 3
		return 0
	}

	b.EndAttack()

	return 0
}

func (b *Boss) PerformAttackLunge() {
	lungeDistance := 8.0

	if b.FacingLeft {
		b.X -= lungeDistance
	} else {
		b.X += lungeDistance
	}

	b.KeepInsideScreen()
}

func (b *Boss) TryDamagePlayer(player *Player) int {
	attackBox := b.AttackBox()
	playerBox := player.HitBox()

	if !Intersects(attackBox, playerBox) {
		return 0
	}

	b.AttackHit = true

	bossCenter := b.X + b.HitWidth/2
	damage := b.GetAttackDamage()

	if player.Hit(damage, bossCenter) {
		return damage
	}

	return 0
}

func (b *Boss) EndAttack() {
	b.Attacking = false

	b.AttackTimer = 0
	b.AttackHit = false
	b.AttackLungeDone = false

	b.CurrentAttackFrame = 0

	b.AttackCooldown = b.GetAttackCooldown()
}

func (b *Boss) UpdateDeath() {
	b.DeathTimer++

	b.Moving = false
	b.Attacking = false

	b.AttackTimer = 0
	b.AttackHit = false
	b.AttackLungeDone = false

	b.CurrentAttackFrame = 0
	b.CurrentWalkFrame = 0

	b.X += b.KnockbackX
	b.KnockbackX *= 0.92

	b.KeepInsideScreen()

	if b.DeathTimer >= b.DeathDuration {
		b.Alive = false
	}
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

func (b *Boss) VisualGroundY() float64 {
	return groundY - b.VisualGroundOffset
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
	attackWidth := 70.0
	attackHeight := 92.0

	if b.Phase >= 2 {
		attackWidth = 82
	}

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
	if !b.Alive {
		return false
	}

	if b.Dying {
		return false
	}

	if b.PhaseTransition {
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

func (b *Boss) CurrentImage() *ebiten.Image {
	if b.Attacking {
		return b.AttackFrames[b.CurrentAttackFrame]
	}

	if b.Moving {
		return b.WalkFrames[b.CurrentWalkFrame]
	}

	return b.AttackFrames[0]
}

func (b *Boss) CurrentContentBounds() image.Rectangle {
	if b.Attacking {
		return b.AttackFrameBounds[b.CurrentAttackFrame]
	}

	if b.Moving {
		return b.WalkFrameBounds[b.CurrentWalkFrame]
	}

	return b.AttackFrameBounds[0]
}

func (b *Boss) CurrentFootY() float64 {
	if b.Attacking {
		return b.AttackFootY[b.CurrentAttackFrame]
	}

	if b.Moving {
		return b.WalkFootY[b.CurrentWalkFrame]
	}

	return b.AttackFootY[0]
}

func (b *Boss) Draw(screen *ebiten.Image) {
	if !b.Alive {
		return
	}

	b.DrawShadow(screen)

	if b.Phase >= 2 || b.PhaseTransition {
		b.DrawAura(screen)
	}

	if b.Attacking && b.CurrentAttackFrame == 2 {
		b.DrawSlashEffect(screen)
	}

	img := b.CurrentImage()
	contentBounds := b.CurrentContentBounds()
	footY := b.CurrentFootY()

	if b.Dying {
		img = b.AttackFrames[0]
		contentBounds = b.AttackFrameBounds[0]
		footY = b.AttackFootY[0]
	}

	contentHeight := float64(contentBounds.Dy())

	if contentHeight <= 0 {
		return
	}

	frameScale := b.VisualHeight / contentHeight

	anchorX := float64(contentBounds.Min.X+contentBounds.Max.X) / 2
	anchorY := footY

	centerX := b.X + b.HitWidth/2
	bottomY := b.VisualGroundY()

	deathAlpha := float32(1)

	if b.Dying {
		progress := float64(b.DeathTimer) / float64(b.DeathDuration)

		if progress > 1 {
			progress = 1
		}

		bottomY += progress * 25
		deathAlpha = float32(1 - progress)
	}

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Translate(-anchorX, -anchorY)

	if b.Moving && !b.Attacking && !b.Dying {
		if b.FacingLeft {
			options.GeoM.Scale(frameScale, frameScale)
		} else {
			options.GeoM.Scale(-frameScale, frameScale)
		}
	} else {
		if b.FacingLeft {
			options.GeoM.Scale(-frameScale, frameScale)
		} else {
			options.GeoM.Scale(frameScale, frameScale)
		}
	}

	options.GeoM.Translate(centerX, bottomY)

	if b.Phase >= 2 && !b.Dying {
		pulse := float32((math.Sin(float64(b.AuraTimer)*0.08) + 1) / 2)
		options.ColorScale.Scale(1.10+0.12*pulse, 0.88, 0.88, 1)
	}

	if b.PhaseTransition {
		pulse := float32((math.Sin(float64(b.PhaseTransitionTimer)*0.45) + 1) / 2)
		options.ColorScale.Scale(1.35+0.45*pulse, 0.70, 0.70, 1)
	}

	if b.FlashTimer > 0 {
		options.ColorScale.Scale(2, 2, 2, 1)
	}

	if b.Dying {
		options.ColorScale.Scale(1, 1, 1, deathAlpha)
	}

	screen.DrawImage(img, options)
}

func (b *Boss) DrawShadow(screen *ebiten.Image) {
	shadowWidth := 46.0
	shadowHeight := 4.0

	centerX := b.X + b.HitWidth/2
	shadowY := b.VisualGroundY() - 2

	ebitenutil.DrawRect(screen, centerX-shadowWidth/2, shadowY, shadowWidth, shadowHeight, color.RGBA{R: 0, G: 0, B: 0, A: 95})
}

func (b *Boss) DrawAura(screen *ebiten.Image) {
	centerX := b.X + b.HitWidth/2
	centerY := b.VisualGroundY() - b.HitHeight/2

	timer := float64(b.AuraTimer)

	intensity := 1.0

	if b.PhaseTransition {
		intensity = 1.8
	}

	pulse := (math.Sin(timer*0.08) + 1) / 2

	glowWidth := 48 + pulse*12*intensity

	alpha := uint8(55)

	if b.PhaseTransition {
		alpha = 90
	}

	ebitenutil.DrawRect(screen, centerX-glowWidth/2, b.VisualGroundY()-4, glowWidth, 3, color.RGBA{R: 255, G: 35, B: 20, A: alpha})

	particleCount := 16

	if b.PhaseTransition {
		particleCount = 24
	}

	for i := 0; i < particleCount; i++ {
		angle := timer*0.025 + float64(i)*0.8

		radiusX := 30 + float64(i%4)*5
		radiusY := 36 + float64(i%5)*6

		x := centerX + math.Cos(angle)*radiusX
		y := centerY + math.Sin(angle*1.3)*radiusY

		flicker := (math.Sin(timer*0.12+float64(i)*1.7) + 1) / 2
		particleAlpha := uint8(80 + flicker*130)

		size := 1.0

		if i%5 == 0 {
			size = 2
		}

		var particleColor color.RGBA

		if i%3 == 0 {
			particleColor = color.RGBA{R: 255, G: 150, B: 45, A: particleAlpha}
		} else {
			particleColor = color.RGBA{R: 255, G: 45, B: 30, A: particleAlpha}
		}

		ebitenutil.DrawRect(screen, x, y, size, size, particleColor)
	}
}

func (b *Boss) DrawSlashEffect(screen *ebiten.Image) {
	centerX := b.X + b.HitWidth/2
	centerY := b.VisualGroundY() - 48

	direction := 1.0

	if b.FacingLeft {
		direction = -1
	}

	if b.Phase >= 2 {
		for i := 0; i < 7; i++ {
			distance := 24 + float64(i)*6
			x := centerX + direction*distance
			y := centerY - 15 + float64(i)*4

			ebitenutil.DrawRect(screen, x, y, 4, 1, color.RGBA{R: 255, G: 80, B: 45, A: uint8(230 - i*24)})
		}

		return
	}

	for i := 0; i < 5; i++ {
		distance := 22 + float64(i)*6
		x := centerX + direction*distance
		y := centerY - 10 + float64(i)*4

		ebitenutil.DrawRect(screen, x, y, 3, 1, color.RGBA{R: 220, G: 220, B: 235, A: uint8(210 - i*30)})
	}
}
