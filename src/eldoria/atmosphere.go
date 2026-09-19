package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ParticleLeaf = iota
	ParticleFirefly
	ParticleDust
	ParticleColdDust
	ParticleEmber
)

type AmbientParticle struct {
	X float64
	Y float64

	VX float64
	VY float64

	Width  float64
	Height float64

	Life    int
	MaxLife int

	Type int
}

type Atmosphere struct {
	Particles []AmbientParticle
	Random    *rand.Rand
	Level     int
}

func NewAtmosphere() *Atmosphere {
	return &Atmosphere{
		Particles: []AmbientParticle{},
		Random: rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
		Level: 0,
	}
}

func (a *Atmosphere) Reset(levelNumber int) {
	a.Particles = []AmbientParticle{}
	a.Level = levelNumber

	switch levelNumber {
	case 1:
		for i := 0; i < 16; i++ {
			a.spawnForestParticle(true)
		}

	case 2:
		for i := 0; i < 20; i++ {
			a.spawnDustParticle(true)
		}

	case 3:
		for i := 0; i < 18; i++ {
			a.spawnColdDustParticle(true)
		}

	case 4:
		for i := 0; i < 24; i++ {
			a.spawnEmberParticle(true)
		}
	}
}

func (a *Atmosphere) Update(levelNumber int) {
	if a.Level != levelNumber {
		a.Reset(levelNumber)
	}

	switch levelNumber {
	case 1:
		if a.Random.Intn(12) == 0 {
			a.spawnForestParticle(false)
		}

	case 2:
		if a.Random.Intn(10) == 0 {
			a.spawnDustParticle(false)
		}

	case 3:
		if a.Random.Intn(11) == 0 {
			a.spawnColdDustParticle(false)
		}

	case 4:
		if a.Random.Intn(5) == 0 {
			a.spawnEmberParticle(false)
		}
	}

	active := make([]AmbientParticle, 0, len(a.Particles))

	for i := range a.Particles {
		particle := a.Particles[i]

		particle.X += particle.VX
		particle.Y += particle.VY

		particle.Life--

		switch particle.Type {
		case ParticleLeaf:
			particle.VX += (a.Random.Float64() - 0.5) * 0.015

		case ParticleFirefly:
			particle.VX += (a.Random.Float64() - 0.5) * 0.01
			particle.VY += (a.Random.Float64() - 0.5) * 0.01

		case ParticleDust:
			particle.VX += (a.Random.Float64() - 0.5) * 0.003

		case ParticleColdDust:
			particle.VX += (a.Random.Float64() - 0.5) * 0.004

		case ParticleEmber:
			particle.VX += (a.Random.Float64() - 0.5) * 0.01
		}

		if particle.Life <= 0 {
			continue
		}

		if particle.X < -5 || particle.X > 325 {
			continue
		}

		if particle.Y < -5 || particle.Y > 185 {
			continue
		}

		active = append(active, particle)
	}

	a.Particles = active
}

func (a *Atmosphere) spawnForestParticle(initial bool) {
	if a.Random.Intn(4) == 0 {
		a.spawnFirefly(initial)
		return
	}

	x := -3.0

	if initial {
		x = a.Random.Float64() * 320
	}

	y := a.Random.Float64() * 135

	life := 500 + a.Random.Intn(500)

	a.Particles = append(
		a.Particles,
		AmbientParticle{
			X:       x,
			Y:       y,
			VX:      0.15 + a.Random.Float64()*0.25,
			VY:      0.08 + a.Random.Float64()*0.18,
			Width:   2,
			Height:  1,
			Life:    life,
			MaxLife: life,
			Type:    ParticleLeaf,
		},
	)
}

func (a *Atmosphere) spawnFirefly(initial bool) {
	x := -2.0

	if initial {
		x = a.Random.Float64() * 320
	}

	y := 25 + a.Random.Float64()*105

	life := 350 + a.Random.Intn(400)

	a.Particles = append(
		a.Particles,
		AmbientParticle{
			X:       x,
			Y:       y,
			VX:      0.05 + a.Random.Float64()*0.12,
			VY:      (a.Random.Float64() - 0.5) * 0.08,
			Width:   1,
			Height:  1,
			Life:    life,
			MaxLife: life,
			Type:    ParticleFirefly,
		},
	)
}

func (a *Atmosphere) spawnDustParticle(initial bool) {
	x := -2.0

	if initial {
		x = a.Random.Float64() * 320
	}

	y := 20 + a.Random.Float64()*125

	life := 500 + a.Random.Intn(500)

	a.Particles = append(
		a.Particles,
		AmbientParticle{
			X:       x,
			Y:       y,
			VX:      0.03 + a.Random.Float64()*0.08,
			VY:      0.01 + a.Random.Float64()*0.03,
			Width:   1,
			Height:  1,
			Life:    life,
			MaxLife: life,
			Type:    ParticleDust,
		},
	)
}

func (a *Atmosphere) spawnColdDustParticle(initial bool) {
	x := -2.0

	if initial {
		x = a.Random.Float64() * 320
	}

	y := 15 + a.Random.Float64()*130

	life := 500 + a.Random.Intn(500)

	a.Particles = append(
		a.Particles,
		AmbientParticle{
			X:       x,
			Y:       y,
			VX:      0.02 + a.Random.Float64()*0.06,
			VY:      0.015 + a.Random.Float64()*0.035,
			Width:   1,
			Height:  1,
			Life:    life,
			MaxLife: life,
			Type:    ParticleColdDust,
		},
	)
}

func (a *Atmosphere) spawnEmberParticle(initial bool) {
	x := a.Random.Float64() * 320
	y := 178.0

	if initial {
		y = 50 + a.Random.Float64()*128
	}

	life := 180 + a.Random.Intn(300)

	a.Particles = append(
		a.Particles,
		AmbientParticle{
			X:       x,
			Y:       y,
			VX:      (a.Random.Float64() - 0.5) * 0.12,
			VY:      -(0.15 + a.Random.Float64()*0.3),
			Width:   1,
			Height:  2,
			Life:    life,
			MaxLife: life,
			Type:    ParticleEmber,
		},
	)
}

func (a *Atmosphere) Draw(screen *ebiten.Image) {
	for _, particle := range a.Particles {
		alphaRatio := float64(particle.Life) / float64(particle.MaxLife)

		if alphaRatio > 1 {
			alphaRatio = 1
		}

		if alphaRatio < 0 {
			alphaRatio = 0
		}

		var particleColor color.RGBA

		switch particle.Type {
		case ParticleLeaf:
			particleColor = color.RGBA{
				R: 90,
				G: 130,
				B: 45,
				A: uint8(170 * alphaRatio),
			}

		case ParticleFirefly:
			particleColor = color.RGBA{
				R: 255,
				G: 225,
				B: 110,
				A: uint8(220 * alphaRatio),
			}

		case ParticleDust:
			particleColor = color.RGBA{
				R: 190,
				G: 175,
				B: 150,
				A: uint8(100 * alphaRatio),
			}

		case ParticleColdDust:
			particleColor = color.RGBA{
				R: 150,
				G: 175,
				B: 190,
				A: uint8(95 * alphaRatio),
			}

		case ParticleEmber:
			particleColor = color.RGBA{
				R: 255,
				G: 90,
				B: 35,
				A: uint8(220 * alphaRatio),
			}
		}

		ebitenutil.DrawRect(
			screen,
			particle.X,
			particle.Y,
			particle.Width,
			particle.Height,
			particleColor,
		)
	}
}
