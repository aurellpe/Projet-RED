package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	AtmosphereParticleDust = iota
	AtmosphereParticleGlow
	AtmosphereParticleLeaf
	AtmosphereParticleEmber
	AtmosphereParticleCold
)

type AtmosphereParticle struct {
	X float64
	Y float64

	VX float64
	VY float64

	Size float64

	Alpha float64

	Phase float64

	Life    int
	MaxLife int

	Kind int

	Depth float64
}

type Atmosphere struct {
	Particles []AtmosphereParticle

	Level int

	Timer int

	Random *rand.Rand
}

func NewAtmosphere() *Atmosphere {
	return &Atmosphere{
		Particles: []AtmosphereParticle{},
		Level:     1,
		Timer:     0,
		Random: rand.New(
			rand.NewSource(
				time.Now().UnixNano(),
			),
		),
	}
}

func (a *Atmosphere) Reset(level int) {
	a.Level = level
	a.Timer = 0
	a.Particles = []AtmosphereParticle{}

	count := a.ParticleCountForLevel(level)

	for i := 0; i < count; i++ {
		particle := a.NewParticle(level, true)

		a.Particles = append(
			a.Particles,
			particle,
		)
	}
}

func (a *Atmosphere) ParticleCountForLevel(level int) int {
	switch level {
	case 1:
		return 85

	case 2:
		return 90

	case 3:
		return 105

	case 4:
		return 130
	}

	return 80
}

func (a *Atmosphere) Update(level int) {
	if a.Random == nil {
		a.Random = rand.New(
			rand.NewSource(
				time.Now().UnixNano(),
			),
		)
	}

	if a.Level != level || len(a.Particles) == 0 {
		a.Reset(level)
	}

	a.Timer++

	for i := range a.Particles {
		a.UpdateParticle(
			&a.Particles[i],
			level,
		)
	}
}

func (a *Atmosphere) UpdateParticle(particle *AtmosphereParticle, level int) {
	particle.Phase += 0.02 + particle.Depth*0.02

	switch level {
	case 1:
		a.UpdateForestParticle(particle)

	case 2:
		a.UpdateRuinsParticle(particle)

	case 3:
		a.UpdateCastleParticle(particle)

	case 4:
		a.UpdateBossParticle(particle)

	default:
		a.UpdateRuinsParticle(particle)
	}

	particle.Life--

	if a.ShouldRespawn(particle) {
		newParticle := a.NewParticle(level, false)
		*particle = newParticle
	}
}

func (a *Atmosphere) UpdateForestParticle(particle *AtmosphereParticle) {
	switch particle.Kind {
	case AtmosphereParticleLeaf:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.09
		particle.Y += math.Cos(particle.Phase*0.7) * 0.025

	case AtmosphereParticleGlow:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.04
		particle.Y += math.Sin(particle.Phase*1.4) * 0.06

	default:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.02
	}
}

func (a *Atmosphere) UpdateRuinsParticle(particle *AtmosphereParticle) {
	switch particle.Kind {
	case AtmosphereParticleGlow:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.04

	default:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.025
	}
}

func (a *Atmosphere) UpdateCastleParticle(particle *AtmosphereParticle) {
	switch particle.Kind {
	case AtmosphereParticleCold:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.055

	case AtmosphereParticleGlow:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.Y += math.Sin(particle.Phase) * 0.04

	default:
		particle.X += particle.VX
		particle.Y += particle.VY
	}
}

func (a *Atmosphere) UpdateBossParticle(particle *AtmosphereParticle) {
	switch particle.Kind {
	case AtmosphereParticleEmber:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase) * 0.09

		particle.VY -= 0.0007

	case AtmosphereParticleGlow:
		particle.X += particle.VX
		particle.Y += particle.VY

		particle.X += math.Sin(particle.Phase*1.4) * 0.05

	default:
		particle.X += particle.VX
		particle.Y += particle.VY
	}
}

func (a *Atmosphere) ShouldRespawn(particle *AtmosphereParticle) bool {
	if particle.Life <= 0 {
		return true
	}

	if particle.X < -20 {
		return true
	}

	if particle.X > 340 {
		return true
	}

	if particle.Y < -20 {
		return true
	}

	if particle.Y > 200 {
		return true
	}

	return false
}

func (a *Atmosphere) NewParticle(level int, anywhere bool) AtmosphereParticle {
	switch level {
	case 1:
		return a.NewForestParticle(anywhere)

	case 2:
		return a.NewRuinsParticle(anywhere)

	case 3:
		return a.NewCastleParticle(anywhere)

	case 4:
		return a.NewBossParticle(anywhere)
	}

	return a.NewRuinsParticle(anywhere)
}

func (a *Atmosphere) NewForestParticle(anywhere bool) AtmosphereParticle {
	randomKind := a.Random.Intn(100)

	kind := AtmosphereParticleLeaf

	if randomKind < 32 {
		kind = AtmosphereParticleGlow
	} else if randomKind < 52 {
		kind = AtmosphereParticleDust
	}

	x := a.Random.Float64() * 320
	y := a.Random.Float64() * 175

	if !anywhere {
		switch kind {
		case AtmosphereParticleLeaf:
			x = -8
			y = a.Random.Float64() * 150

		case AtmosphereParticleGlow:
			x = a.Random.Float64() * 320
			y = 185

		default:
			x = a.Random.Float64() * 320
			y = 185
		}
	}

	depth := 0.35 + a.Random.Float64()*0.95

	if kind == AtmosphereParticleGlow {
		return AtmosphereParticle{
			X: x,
			Y: y,

			VX: -0.04 + a.Random.Float64()*0.08,
			VY: -0.025 - a.Random.Float64()*0.055,

			Size: 0.7 + a.Random.Float64()*1.0,

			Alpha: 0.45 + a.Random.Float64()*0.45,

			Phase: a.Random.Float64() * math.Pi * 2,

			Life:    550 + a.Random.Intn(450),
			MaxLife: 1000,

			Kind: kind,

			Depth: depth,
		}
	}

	if kind == AtmosphereParticleDust {
		return AtmosphereParticle{
			X: x,
			Y: y,

			VX: -0.015 + a.Random.Float64()*0.03,
			VY: -0.015 - a.Random.Float64()*0.035,

			Size: 0.5 + a.Random.Float64()*0.8,

			Alpha: 0.18 + a.Random.Float64()*0.30,

			Phase: a.Random.Float64() * math.Pi * 2,

			Life:    650 + a.Random.Intn(450),
			MaxLife: 1100,

			Kind: kind,

			Depth: depth,
		}
	}

	return AtmosphereParticle{
		X: x,
		Y: y,

		VX: 0.07 + a.Random.Float64()*0.22,
		VY: 0.03 + a.Random.Float64()*0.10,

		Size: 0.8 + a.Random.Float64()*1.5,

		Alpha: 0.30 + a.Random.Float64()*0.40,

		Phase: a.Random.Float64() * math.Pi * 2,

		Life:    650 + a.Random.Intn(450),
		MaxLife: 1100,

		Kind: kind,

		Depth: depth,
	}
}

func (a *Atmosphere) NewRuinsParticle(anywhere bool) AtmosphereParticle {
	randomKind := a.Random.Intn(100)

	kind := AtmosphereParticleDust

	if randomKind < 16 {
		kind = AtmosphereParticleGlow
	}

	x := a.Random.Float64() * 320
	y := a.Random.Float64() * 180

	if !anywhere {
		x = a.Random.Float64() * 320
		y = 190
	}

	depth := 0.3 + a.Random.Float64()*1.0

	if kind == AtmosphereParticleGlow {
		return AtmosphereParticle{
			X: x,
			Y: y,

			VX: -0.025 + a.Random.Float64()*0.05,
			VY: -0.045 - a.Random.Float64()*0.07,

			Size: 0.6 + a.Random.Float64()*0.9,

			Alpha: 0.30 + a.Random.Float64()*0.35,

			Phase: a.Random.Float64() * math.Pi * 2,

			Life:    450 + a.Random.Intn(400),
			MaxLife: 850,

			Kind: kind,

			Depth: depth,
		}
	}

	return AtmosphereParticle{
		X: x,
		Y: y,

		VX: -0.03 + a.Random.Float64()*0.06,
		VY: -0.02 - a.Random.Float64()*0.065,

		Size: 0.5 + a.Random.Float64()*1.35,

		Alpha: 0.20 + a.Random.Float64()*0.36,

		Phase: a.Random.Float64() * math.Pi * 2,

		Life:    650 + a.Random.Intn(500),
		MaxLife: 1150,

		Kind: kind,

		Depth: depth,
	}
}

func (a *Atmosphere) NewCastleParticle(anywhere bool) AtmosphereParticle {
	randomKind := a.Random.Intn(100)

	kind := AtmosphereParticleCold

	if randomKind < 20 {
		kind = AtmosphereParticleGlow
	}

	x := a.Random.Float64() * 320
	y := a.Random.Float64() * 175

	if !anywhere {
		x = 328
		y = a.Random.Float64() * 175
	}

	depth := 0.3 + a.Random.Float64()*1.0

	if kind == AtmosphereParticleGlow {
		return AtmosphereParticle{
			X: x,
			Y: y,

			VX: -0.06 - a.Random.Float64()*0.09,
			VY: -0.02 + a.Random.Float64()*0.04,

			Size: 0.7 + a.Random.Float64()*0.8,

			Alpha: 0.40 + a.Random.Float64()*0.42,

			Phase: a.Random.Float64() * math.Pi * 2,

			Life:    500 + a.Random.Intn(450),
			MaxLife: 950,

			Kind: kind,

			Depth: depth,
		}
	}

	speedMultiplier := 0.65 + depth*0.55

	return AtmosphereParticle{
		X: x,
		Y: y,

		VX: (-0.075 - a.Random.Float64()*0.20) * speedMultiplier,
		VY: 0.01 + a.Random.Float64()*0.055,

		Size: 0.5 + depth*a.Random.Float64()*1.35,

		Alpha: 0.22 + depth*0.23 + a.Random.Float64()*0.20,

		Phase: a.Random.Float64() * math.Pi * 2,

		Life:    750 + a.Random.Intn(500),
		MaxLife: 1250,

		Kind: kind,

		Depth: depth,
	}
}

func (a *Atmosphere) NewBossParticle(anywhere bool) AtmosphereParticle {
	randomKind := a.Random.Intn(100)

	kind := AtmosphereParticleEmber

	if randomKind < 24 {
		kind = AtmosphereParticleGlow
	}

	x := a.Random.Float64() * 320
	y := a.Random.Float64() * 180

	if !anywhere {
		x = a.Random.Float64() * 320
		y = 190
	}

	depth := 0.4 + a.Random.Float64()*1.0

	if kind == AtmosphereParticleGlow {
		return AtmosphereParticle{
			X: x,
			Y: y,

			VX: -0.035 + a.Random.Float64()*0.07,
			VY: -0.10 - a.Random.Float64()*0.14,

			Size: 0.9 + a.Random.Float64()*1.5,

			Alpha: 0.40 + a.Random.Float64()*0.48,

			Phase: a.Random.Float64() * math.Pi * 2,

			Life:    350 + a.Random.Intn(400),
			MaxLife: 750,

			Kind: kind,

			Depth: depth,
		}
	}

	return AtmosphereParticle{
		X: x,
		Y: y,

		VX: -0.075 + a.Random.Float64()*0.15,
		VY: -0.14 - a.Random.Float64()*0.25,

		Size: 0.65 + a.Random.Float64()*1.7,

		Alpha: 0.45 + a.Random.Float64()*0.45,

		Phase: a.Random.Float64() * math.Pi * 2,

		Life:    400 + a.Random.Intn(450),
		MaxLife: 850,

		Kind: kind,

		Depth: depth,
	}
}

func (a *Atmosphere) Draw(screen *ebiten.Image) {
	if screen == nil {
		return
	}

	for i := range a.Particles {
		particle := &a.Particles[i]

		switch a.Level {
		case 1:
			a.DrawForestParticle(screen, particle)

		case 2:
			a.DrawRuinsParticle(screen, particle)

		case 3:
			a.DrawCastleParticle(screen, particle)

		case 4:
			a.DrawBossParticle(screen, particle)

		default:
			a.DrawRuinsParticle(screen, particle)
		}
	}
}

func atmosphereAlpha(value float64) uint8 {
	if value < 0 {
		value = 0
	}

	if value > 1 {
		value = 1
	}

	return uint8(value * 255)
}

func (a *Atmosphere) ParticleFade(particle *AtmosphereParticle) float64 {
	if particle.MaxLife <= 0 {
		return 1
	}

	ratio := float64(particle.Life) / float64(particle.MaxLife)

	if ratio < 0 {
		ratio = 0
	}

	if ratio > 1 {
		ratio = 1
	}

	fade := 1.0

	if ratio < 0.15 {
		fade = ratio / 0.15
	}

	return fade
}

func (a *Atmosphere) DrawForestParticle(screen *ebiten.Image, particle *AtmosphereParticle) {
	fade := a.ParticleFade(particle)

	if particle.Kind == AtmosphereParticleGlow {
		pulse := (math.Sin(particle.Phase*2) + 1) / 2

		alpha := atmosphereAlpha(
			particle.Alpha * fade * (0.50 + pulse*0.50),
		)

		size := particle.Size

		ebitenutil.DrawRect(
			screen,
			particle.X-size,
			particle.Y-size,
			size*2,
			size*2,
			color.RGBA{
				R: 120,
				G: 230,
				B: 150,
				A: alpha / 4,
			},
		)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			1,
			1,
			color.RGBA{
				R: 210,
				G: 255,
				B: 180,
				A: alpha,
			},
		)

		return
	}

	if particle.Kind == AtmosphereParticleDust {
		alpha := atmosphereAlpha(
			particle.Alpha * fade,
		)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			particle.Size,
			particle.Size,
			color.RGBA{
				R: 180,
				G: 205,
				B: 150,
				A: alpha,
			},
		)

		return
	}

	alpha := atmosphereAlpha(
		particle.Alpha * fade,
	)

	size := particle.Size

	ebitenutil.DrawRect(
		screen,
		particle.X,
		particle.Y,
		size,
		size*0.55,
		color.RGBA{
			R: 105,
			G: 150,
			B: 65,
			A: alpha,
		},
	)
}

func (a *Atmosphere) DrawRuinsParticle(screen *ebiten.Image, particle *AtmosphereParticle) {
	fade := a.ParticleFade(particle)

	if particle.Kind == AtmosphereParticleGlow {
		pulse := (math.Sin(particle.Phase*2.3) + 1) / 2

		alpha := atmosphereAlpha(
			particle.Alpha * fade * (0.55 + pulse*0.45),
		)

		ebitenutil.DrawRect(
			screen,
			particle.X-1,
			particle.Y-1,
			3,
			3,
			color.RGBA{
				R: 210,
				G: 180,
				B: 125,
				A: alpha / 5,
			},
		)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			1,
			1,
			color.RGBA{
				R: 245,
				G: 220,
				B: 170,
				A: alpha,
			},
		)

		return
	}

	alpha := atmosphereAlpha(
		particle.Alpha * fade,
	)

	size := particle.Size

	ebitenutil.DrawRect(
		screen,
		particle.X,
		particle.Y,
		size,
		size,
		color.RGBA{
			R: 180,
			G: 170,
			B: 145,
			A: alpha,
		},
	)

	if particle.Depth > 0.95 {
		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y+size,
			1,
			1,
			color.RGBA{
				R: 205,
				G: 195,
				B: 170,
				A: alpha / 2,
			},
		)
	}
}

func (a *Atmosphere) DrawCastleParticle(screen *ebiten.Image, particle *AtmosphereParticle) {
	fade := a.ParticleFade(particle)

	if particle.Kind == AtmosphereParticleGlow {
		pulse := (math.Sin(particle.Phase*2.4) + 1) / 2

		alpha := atmosphereAlpha(
			particle.Alpha * fade * (0.55 + pulse*0.45),
		)

		ebitenutil.DrawRect(
			screen,
			particle.X-1,
			particle.Y-1,
			3,
			3,
			color.RGBA{
				R: 110,
				G: 175,
				B: 255,
				A: alpha / 5,
			},
		)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			1,
			1,
			color.RGBA{
				R: 205,
				G: 230,
				B: 255,
				A: alpha,
			},
		)

		return
	}

	alpha := atmosphereAlpha(
		particle.Alpha * fade,
	)

	size := particle.Size

	if particle.Depth > 0.9 {
		ebitenutil.DrawRect(
			screen,
			particle.X-1,
			particle.Y,
			size+2,
			1,
			color.RGBA{
				R: 170,
				G: 205,
				B: 240,
				A: alpha / 4,
			},
		)
	}

	ebitenutil.DrawRect(
		screen,
		particle.X,
		particle.Y,
		size,
		size,
		color.RGBA{
			R: 205,
			G: 225,
			B: 245,
			A: alpha,
		},
	)
}

func (a *Atmosphere) DrawBossParticle(screen *ebiten.Image, particle *AtmosphereParticle) {
	fade := a.ParticleFade(particle)

	pulse := (math.Sin(particle.Phase*2.2) + 1) / 2

	alpha := atmosphereAlpha(
		particle.Alpha * fade * (0.65 + pulse*0.35),
	)

	if particle.Kind == AtmosphereParticleGlow {
		size := particle.Size

		ebitenutil.DrawRect(
			screen,
			particle.X-size,
			particle.Y-size,
			size*2,
			size*2,
			color.RGBA{
				R: 255,
				G: 80,
				B: 35,
				A: alpha / 4,
			},
		)

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			1,
			1,
			color.RGBA{
				R: 255,
				G: 220,
				B: 100,
				A: alpha,
			},
		)

		return
	}

	size := particle.Size

	ebitenutil.DrawRect(
		screen,
		particle.X,
		particle.Y,
		size,
		size,
		color.RGBA{
			R: 255,
			G: 105,
			B: 40,
			A: alpha,
		},
	)

	if particle.Depth > 0.85 {
		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y+size,
			1,
			2,
			color.RGBA{
				R: 255,
				G: 55,
				B: 20,
				A: alpha / 2,
			},
		)
	}
}
