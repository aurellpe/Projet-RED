package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth = 320
	gravity     = 0.35
	jumpForce   = -6.5
)

var groundY = 156.0

//go:embed assets/player/walk1.png
var playerWalk1Data []byte

//go:embed assets/player/walk2.png
var playerWalk2Data []byte

//go:embed assets/player/walk3.png
var playerWalk3Data []byte

//go:embed assets/player/walk4.png
var playerWalk4Data []byte

//go:embed assets/player/attack1.png
var playerAttack1Data []byte

//go:embed assets/player/attack2.png
var playerAttack2Data []byte

//go:embed assets/player/attack3.png
var playerAttack3Data []byte

//go:embed assets/player/attack4.png
var playerAttack4Data []byte

//go:embed assets/player/energy/energy1.png
var playerEnergy1Data []byte

//go:embed assets/player/energy/energy2.png
var playerEnergy2Data []byte

//go:embed assets/player/energy/energy3.png
var playerEnergy3Data []byte

//go:embed assets/player/energy/energy4.png
var playerEnergy4Data []byte

//go:embed assets/player/energy/energy5.png
var playerEnergy5Data []byte

//go:embed assets/player/energy/energy6.png
var playerEnergy6Data []byte

//go:embed assets/player/energy/energy7.png
var playerEnergy7Data []byte

//go:embed assets/player/energy/energy8.png
var playerEnergy8Data []byte

type DodgeDust struct {
	X       float64
	Y       float64
	VX      float64
	VY      float64
	Life    int
	MaxLife int
}

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

	WalkFrames      []*ebiten.Image
	WalkFrameBounds []image.Rectangle
	WalkFootY       []float64

	AttackFrames      []*ebiten.Image
	AttackFrameBounds []image.Rectangle
	AttackFootY       []float64

	EnergyFrames      []*ebiten.Image
	EnergyFrameBounds []image.Rectangle
	EnergyFootY       []float64

	EnergyAttacking    bool
	CurrentEnergyFrame int

	CurrentWalkFrame int
	WalkFrameTimer   int
	WalkFrameSpeed   int
	WalkSequenceStep int

	Moving bool

	FacingLeft bool

	CurrentAttackFrame int
	Attacking          bool
	AttackTimer        int
	AttackDuration     int
	AttackFrameSpeed   int
	AttackHit          bool
	AttackDamage       int

	HurtTimer    int
	HurtDuration int

	HP    int
	MaxHP int

	Alive bool

	InvincibleTimer int
	FlashTimer      int

	VisualHeight       float64
	VisualGroundOffset float64

	Dodging            bool
	DodgeTimer         int
	DodgeDuration      int
	DodgeCooldown      int
	DodgeCooldownTimer int
	DodgeSpeed         float64
	DodgeDirection     float64

	DodgeDustParticles []DodgeDust
}

func findPlayerContentBounds(img image.Image) image.Rectangle {
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

func findPlayerFootY(img image.Image, content image.Rectangle) float64 {
	if content.Empty() {
		return float64(img.Bounds().Max.Y)
	}

	contentWidth := content.Dx()

	left := content.Min.X + contentWidth*20/100
	right := content.Min.X + contentWidth*80/100

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

func loadPlayerFrame(data []byte) (*ebiten.Image, image.Rectangle, float64) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	contentBounds := findPlayerContentBounds(img)
	footY := findPlayerFootY(img, contentBounds)

	return ebiten.NewImageFromImage(img), contentBounds, footY
}

func NewPlayer() *Player {
	walk1, walkBounds1, walkFoot1 := loadPlayerFrame(playerWalk1Data)
	walk2, walkBounds2, walkFoot2 := loadPlayerFrame(playerWalk2Data)
	walk3, walkBounds3, walkFoot3 := loadPlayerFrame(playerWalk3Data)
	walk4, walkBounds4, walkFoot4 := loadPlayerFrame(playerWalk4Data)

	attack1, attackBounds1, attackFoot1 := loadPlayerFrame(playerAttack1Data)
	attack2, attackBounds2, attackFoot2 := loadPlayerFrame(playerAttack2Data)
	attack3, attackBounds3, attackFoot3 := loadPlayerFrame(playerAttack3Data)
	attack4, attackBounds4, attackFoot4 := loadPlayerFrame(playerAttack4Data)

	energy1, energyBounds1, energyFoot1 := loadPlayerFrame(playerEnergy1Data)
	energy2, energyBounds2, energyFoot2 := loadPlayerFrame(playerEnergy2Data)
	energy3, energyBounds3, energyFoot3 := loadPlayerFrame(playerEnergy3Data)
	energy4, energyBounds4, energyFoot4 := loadPlayerFrame(playerEnergy4Data)
	energy5, energyBounds5, energyFoot5 := loadPlayerFrame(playerEnergy5Data)
	energy6, energyBounds6, energyFoot6 := loadPlayerFrame(playerEnergy6Data)
	energy7, energyBounds7, energyFoot7 := loadPlayerFrame(playerEnergy7Data)
	energy8, energyBounds8, energyFoot8 := loadPlayerFrame(playerEnergy8Data)

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

	energyFrames := []*ebiten.Image{
		energy1,
		energy2,
		energy3,
		energy4,
		energy5,
		energy6,
		energy7,
		energy8,
	}

	energyBounds := []image.Rectangle{
		energyBounds1,
		energyBounds2,
		energyBounds3,
		energyBounds4,
		energyBounds5,
		energyBounds6,
		energyBounds7,
		energyBounds8,
	}

	energyFeet := []float64{
		energyFoot1,
		energyFoot2,
		energyFoot3,
		energyFoot4,
		energyFoot5,
		energyFoot6,
		energyFoot7,
		energyFoot8,
	}

	scale := 0.08

	idleWidth := float64(walk1.Bounds().Dx()) * scale
	idleHeight := float64(walk1.Bounds().Dy()) * scale

	visualHeight := float64(walkBounds1.Dy()) * scale

	bottomPaddingPixels := float64(walk1.Bounds().Max.Y) - walkFoot1
	visualGroundOffset := bottomPaddingPixels * scale

	if visualGroundOffset < 0 {
		visualGroundOffset = 0
	}

	return &Player{
		X: 70,
		Y: groundY - idleHeight,

		Speed: 1.5,

		VelocityY: 0,
		OnGround:  true,

		KnockbackX: 0,

		Scale: scale,

		Width:  idleWidth,
		Height: idleHeight,

		WalkFrames:      walkFrames,
		WalkFrameBounds: walkBounds,
		WalkFootY:       walkFeet,

		AttackFrames:      attackFrames,
		AttackFrameBounds: attackBounds,
		AttackFootY:       attackFeet,

		EnergyFrames:      energyFrames,
		EnergyFrameBounds: energyBounds,
		EnergyFootY:       energyFeet,

		EnergyAttacking:    false,
		CurrentEnergyFrame: 0,

		CurrentWalkFrame: 0,
		WalkFrameTimer:   0,
		WalkFrameSpeed:   7,
		WalkSequenceStep: 0,

		Moving: false,

		FacingLeft: false,

		CurrentAttackFrame: 0,
		Attacking:          false,
		AttackTimer:        0,
		AttackDuration:     24,
		AttackFrameSpeed:   6,
		AttackHit:          false,
		AttackDamage:       50,

		HurtTimer:    0,
		HurtDuration: 12,

		HP:    100,
		MaxHP: 100,

		Alive: true,

		InvincibleTimer: 0,
		FlashTimer:      0,

		VisualHeight:       visualHeight,
		VisualGroundOffset: visualGroundOffset,

		Dodging:            false,
		DodgeTimer:         0,
		DodgeDuration:      14,
		DodgeCooldown:      38,
		DodgeCooldownTimer: 0,
		DodgeSpeed:         4,
		DodgeDirection:     1,

		DodgeDustParticles: []DodgeDust{},
	}
}

func (p *Player) Update() {
	p.UpdateDodgeDust()

	if p.InvincibleTimer > 0 {
		p.InvincibleTimer--
	}

	if p.FlashTimer > 0 {
		p.FlashTimer--
	}

	if p.HurtTimer > 0 {
		p.HurtTimer--
	}

	if p.DodgeCooldownTimer > 0 {
		p.DodgeCooldownTimer--
	}

	if !p.Alive {
		return
	}

	if p.EnergyAttacking {
		p.UpdateGravity()
		p.KeepInsideScreen()
		return
	}

	if p.Dodging {
		p.UpdateDodge()
		return
	}

	p.Moving = false

	p.X += p.KnockbackX
	p.KnockbackX *= 0.82

	if p.KnockbackX > -0.05 && p.KnockbackX < 0.05 {
		p.KnockbackX = 0
	}

	if p.HurtTimer > 0 {
		p.Attacking = false
		p.AttackTimer = 0
		p.CurrentAttackFrame = 0
		p.AttackHit = false

		p.CurrentWalkFrame = 0
		p.WalkFrameTimer = 0
		p.WalkSequenceStep = 0

		p.UpdateGravity()
		p.KeepInsideScreen()

		return
	}

	if p.Attacking {
		p.UpdateAttackAnimation()
		p.UpdateGravity()
		p.KeepInsideScreen()

		return
	}

	leftPressed := ebiten.IsKeyPressed(ebiten.KeyA)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyD)

	if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) || inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
		if p.OnGround && p.DodgeCooldownTimer <= 0 {
			p.StartDodge(leftPressed, rightPressed)
			return
		}
	}

	if leftPressed && !rightPressed {
		p.X -= p.Speed
		p.Moving = true
		p.FacingLeft = true
	}

	if rightPressed && !leftPressed {
		p.X += p.Speed
		p.Moving = true
		p.FacingLeft = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && p.OnGround {
		p.VelocityY = jumpForce
		p.OnGround = false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		p.StartAttack()
	}

	p.UpdateGravity()
	p.UpdateWalkAnimation()
	p.KeepInsideScreen()
}

func (p *Player) StartDodge(leftPressed bool, rightPressed bool) {
	if p.EnergyAttacking {
		return
	}

	p.Dodging = true
	p.DodgeTimer = p.DodgeDuration
	p.DodgeCooldownTimer = p.DodgeCooldown

	p.Attacking = false
	p.AttackTimer = 0
	p.CurrentAttackFrame = 0
	p.AttackHit = false

	p.Moving = false

	if leftPressed && !rightPressed {
		p.DodgeDirection = -1
		p.FacingLeft = true
	} else if rightPressed && !leftPressed {
		p.DodgeDirection = 1
		p.FacingLeft = false
	} else if p.FacingLeft {
		p.DodgeDirection = -1
	} else {
		p.DodgeDirection = 1
	}

	p.SpawnDodgeDust()
}

func (p *Player) UpdateDodge() {
	p.DodgeTimer--

	p.InvincibleTimer = 2

	p.X += p.DodgeDirection * p.DodgeSpeed

	if p.DodgeTimer%3 == 0 {
		p.SpawnDodgeDust()
	}

	p.UpdateGravity()
	p.KeepInsideScreen()

	if p.DodgeTimer <= 0 {
		p.Dodging = false
		p.DodgeTimer = 0
	}
}

func (p *Player) SpawnDodgeDust() {
	baseX := p.X + p.Width/2
	baseY := p.VisualGroundY() - 2

	for i := 0; i < 4; i++ {
		direction := -p.DodgeDirection

		particle := DodgeDust{
			X:       baseX + direction*float64(i*2),
			Y:       baseY - float64(i%2),
			VX:      direction * (0.25 + float64(i)*0.08),
			VY:      -0.12 - float64(i%2)*0.07,
			Life:    18 + i*2,
			MaxLife: 18 + i*2,
		}

		p.DodgeDustParticles = append(p.DodgeDustParticles, particle)
	}
}

func (p *Player) UpdateDodgeDust() {
	active := []DodgeDust{}

	for _, particle := range p.DodgeDustParticles {
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.VY += 0.015
		particle.Life--

		if particle.Life > 0 {
			active = append(active, particle)
		}
	}

	p.DodgeDustParticles = active
}

func (p *Player) StartAttack() {
	if p.Attacking {
		return
	}

	if p.EnergyAttacking {
		return
	}

	if p.Dodging {
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
	p.WalkSequenceStep = 0
}

func (p *Player) StartEnergyAnimation() {
	p.EnergyAttacking = true
	p.CurrentEnergyFrame = 0

	p.Attacking = false
	p.AttackTimer = 0
	p.CurrentAttackFrame = 0
	p.AttackHit = false

	p.Dodging = false
	p.DodgeTimer = 0

	p.Moving = false
	p.CurrentWalkFrame = 0
	p.WalkFrameTimer = 0
	p.WalkSequenceStep = 0

	p.KnockbackX = 0
	p.VelocityY = 0

	p.Y = groundY - p.Height
	p.OnGround = true
}

func (p *Player) SetEnergyFrame(frame int) {
	if len(p.EnergyFrames) == 0 {
		return
	}

	if frame < 0 {
		frame = 0
	}

	if frame >= len(p.EnergyFrames) {
		frame = len(p.EnergyFrames) - 1
	}

	p.CurrentEnergyFrame = frame
}

func (p *Player) StopEnergyAnimation() {
	p.EnergyAttacking = false
	p.CurrentEnergyFrame = 0
}

func (p *Player) UpdateAttackAnimation() {
	p.AttackTimer++

	frame := p.AttackTimer / p.AttackFrameSpeed

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

func (p *Player) UpdateWalkAnimation() {
	if !p.Moving || !p.OnGround {
		p.CurrentWalkFrame = 0
		p.WalkFrameTimer = 0
		p.WalkSequenceStep = 0
		return
	}

	p.WalkFrameTimer++

	if p.WalkFrameTimer < p.WalkFrameSpeed {
		return
	}

	p.WalkFrameTimer = 0

	walkSequence := []int{
		0,
		1,
		2,
		3,
		2,
		1,
	}

	p.WalkSequenceStep++

	if p.WalkSequenceStep >= len(walkSequence) {
		p.WalkSequenceStep = 0
	}

	p.CurrentWalkFrame = walkSequence[p.WalkSequenceStep]
}

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

func (p *Player) KeepInsideScreen() {
	if p.X < 0 {
		p.X = 0
	}

	if p.X+p.Width > screenWidth {
		p.X = screenWidth - p.Width
	}
}

func (p *Player) VisualGroundY() float64 {
	return p.Y + p.Height - p.VisualGroundOffset
}

func (p *Player) HitBox() Rect {
	bodyWidth := 28.0
	bodyHeight := 55.0

	playerCenter := p.X + p.Width/2
	playerBottom := p.VisualGroundY()

	return Rect{
		X: playerCenter - bodyWidth/2,
		Y: playerBottom - bodyHeight,
		W: bodyWidth,
		H: bodyHeight,
	}
}

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

func (p *Player) Hit(damage int, sourceX float64) bool {
	if !p.Alive {
		return false
	}

	if p.Dodging {
		return false
	}

	if p.InvincibleTimer > 0 {
		return false
	}

	p.StopEnergyAnimation()

	p.HP -= damage

	p.InvincibleTimer = 45
	p.FlashTimer = 7
	p.HurtTimer = p.HurtDuration

	p.Attacking = false
	p.AttackTimer = 0
	p.CurrentAttackFrame = 0
	p.AttackHit = false

	playerCenter := p.X + p.Width/2

	if playerCenter < sourceX {
		p.KnockbackX = -3.2
		p.FacingLeft = false
	} else {
		p.KnockbackX = 3.2
		p.FacingLeft = true
	}

	p.VelocityY = -2.3

	if p.HP <= 0 {
		p.HP = 0
		p.Alive = false
	}

	return true
}

func (p *Player) CurrentImage() *ebiten.Image {
	if p.EnergyAttacking {
		return p.EnergyFrames[p.CurrentEnergyFrame]
	}

	if p.Attacking {
		return p.AttackFrames[p.CurrentAttackFrame]
	}

	return p.WalkFrames[p.CurrentWalkFrame]
}

func (p *Player) CurrentContentBounds() image.Rectangle {
	if p.EnergyAttacking {
		return p.EnergyFrameBounds[p.CurrentEnergyFrame]
	}

	if p.Attacking {
		return p.AttackFrameBounds[p.CurrentAttackFrame]
	}

	return p.WalkFrameBounds[p.CurrentWalkFrame]
}

func (p *Player) CurrentFootY() float64 {
	if p.EnergyAttacking {
		return p.EnergyFootY[p.CurrentEnergyFrame]
	}

	if p.Attacking {
		return p.AttackFootY[p.CurrentAttackFrame]
	}

	return p.WalkFootY[p.CurrentWalkFrame]
}

func (p *Player) DrawDodgeDust(screen *ebiten.Image) {
	for _, particle := range p.DodgeDustParticles {
		ratio := float64(particle.Life) / float64(particle.MaxLife)

		alpha := uint8(140 * ratio)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			2,
			1,
			color.RGBA{
				R: 185,
				G: 175,
				B: 145,
				A: alpha,
			},
		)
	}
}

func (p *Player) DrawAttackTrail(screen *ebiten.Image) {
	if !p.Attacking {
		return
	}

	if p.CurrentAttackFrame < 1 || p.CurrentAttackFrame > 2 {
		return
	}

	centerX := p.X + p.Width/2
	centerY := p.VisualGroundY() - 30

	direction := 1.0

	if p.FacingLeft {
		direction = -1
	}

	for i := 0; i < 7; i++ {
		distance := 18 + float64(i)*5

		x := centerX + direction*distance
		y := centerY - 12 + float64(i)*3

		alpha := uint8(220 - i*25)

		ebitenutil.DrawRect(
			screen,
			x,
			y,
			3,
			1,
			color.RGBA{
				R: 210,
				G: 230,
				B: 255,
				A: alpha,
			},
		)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.DrawDodgeDust(screen)

	if !p.Alive {
		return
	}

	if !p.EnergyAttacking {
		p.DrawAttackTrail(screen)
	}

	img := p.CurrentImage()
	contentBounds := p.CurrentContentBounds()
	footY := p.CurrentFootY()

	contentHeight := float64(contentBounds.Dy())

	if contentHeight <= 0 {
		return
	}

	frameScale := p.VisualHeight / contentHeight

	anchorX := float64(contentBounds.Min.X+contentBounds.Max.X) / 2
	anchorY := footY

	centerX := p.X + p.Width/2
	visualGround := p.VisualGroundY()

	if p.HurtTimer > 0 && !p.EnergyAttacking {
		if p.HurtTimer%4 < 2 {
			centerX -= 1.5
		} else {
			centerX += 1.5
		}
	}

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Translate(-anchorX, -anchorY)

	if p.FacingLeft {
		options.GeoM.Scale(-frameScale, frameScale)
	} else {
		options.GeoM.Scale(frameScale, frameScale)
	}

	options.GeoM.Translate(centerX, visualGround)

	if p.Dodging {
		options.ColorScale.Scale(1, 1, 1, 0.72)
	}

	if p.FlashTimer > 0 {
		options.ColorScale.Scale(2, 2, 2, 1)
	}

	screen.DrawImage(img, options)
}
