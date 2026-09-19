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

//go:embed assets/exit/portal.png
var portalImageData []byte

type PortalParticle struct {
	Angle  float64
	Radius float64
	Speed  float64
	Size   float64
	Phase  float64
}

type Exit struct {
	X float64
	Y float64

	Width  float64
	Height float64

	Active bool

	AnimationTimer int

	Image *ebiten.Image

	GroundOffset float64

	Particles []PortalParticle
}

func cropPortalImage(img image.Image) image.Image {
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
		return img
	}

	rect := image.Rect(minX, minY, maxX+1, maxY+1)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	if subImage, ok := img.(subImager); ok {
		return subImage.SubImage(rect)
	}

	return img
}

func loadPortalImage() *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(portalImageData))
	if err != nil {
		panic(err)
	}

	cropped := cropPortalImage(img)

	return ebiten.NewImageFromImage(cropped)
}

func NewExit(levelGroundY float64) *Exit {
	particles := []PortalParticle{
		{
			Angle:  0,
			Radius: 18,
			Speed:  0.020,
			Size:   1,
			Phase:  0,
		},
		{
			Angle:  1.2,
			Radius: 21,
			Speed:  0.017,
			Size:   1,
			Phase:  0.8,
		},
		{
			Angle:  2.4,
			Radius: 17,
			Speed:  0.022,
			Size:   1,
			Phase:  1.6,
		},
		{
			Angle:  3.6,
			Radius: 20,
			Speed:  0.018,
			Size:   1,
			Phase:  2.4,
		},
		{
			Angle:  4.8,
			Radius: 16,
			Speed:  0.024,
			Size:   1,
			Phase:  3.2,
		},
	}

	return &Exit{
		X: 246,
		Y: levelGroundY - 96,

		Width:  68,
		Height: 96,

		Active: false,

		AnimationTimer: 0,

		Image: loadPortalImage(),

		GroundOffset: 20,

		Particles: particles,
	}
}

func (e *Exit) Update() {
	e.AnimationTimer++

	for i := range e.Particles {
		e.Particles[i].Angle += e.Particles[i].Speed
	}
}

func (e *Exit) HitBox() Rect {
	return Rect{
		X: e.X + 15,
		Y: e.Y + 15,
		W: e.Width - 30,
		H: e.Height - 15,
	}
}

func (e *Exit) PlayerInside(player *Player) bool {
	if player == nil {
		return false
	}

	return Intersects(e.HitBox(), player.HitBox())
}

func (e *Exit) Draw(screen *ebiten.Image) {
	if e.Image == nil {
		return
	}

	e.drawPortalLight(screen)

	if e.Active {
		e.drawParticles(screen)
	}

	e.drawPortalImage(screen)

	if e.Active {
		e.drawGroundEnergy(screen)
	}
}

func (e *Exit) drawPortalLight(screen *ebiten.Image) {
	centerX := e.X + e.Width/2

	pulse := (math.Sin(float64(e.AnimationTimer)*0.055) + 1) / 2

	/*
		Le halo est volontairement beaucoup plus petit qu'avant.
		Il reste uniquement autour du centre bleu du portail.
	*/
	glowWidth := 32.0 + pulse*3
	glowHeight := 60.0 + pulse*3

	glowX := centerX - glowWidth/2
	glowY := e.Y + 24

	alpha := uint8(18)

	if e.Active {
		alpha = uint8(28 + pulse*12)
	}

	ebitenutil.DrawRect(
		screen,
		glowX,
		glowY,
		glowWidth,
		glowHeight,
		color.RGBA{
			R: 35,
			G: 85,
			B: 255,
			A: alpha,
		},
	)

	/*
		Petit noyau lumineux.
		Il ne dépasse presque plus de l'ouverture.
	*/
	coreWidth := 22.0 + pulse*2
	coreHeight := 48.0 + pulse*2

	coreX := centerX - coreWidth/2
	coreY := e.Y + 30

	coreAlpha := uint8(15)

	if e.Active {
		coreAlpha = uint8(30 + pulse*15)
	}

	ebitenutil.DrawRect(
		screen,
		coreX,
		coreY,
		coreWidth,
		coreHeight,
		color.RGBA{
			R: 40,
			G: 115,
			B: 255,
			A: coreAlpha,
		},
	)
}

func (e *Exit) drawPortalImage(screen *ebiten.Image) {
	imageWidth := float64(e.Image.Bounds().Dx())
	imageHeight := float64(e.Image.Bounds().Dy())

	if imageWidth <= 0 || imageHeight <= 0 {
		return
	}

	pulse := (math.Sin(float64(e.AnimationTimer)*0.045) + 1) / 2

	scaleX := e.Width / imageWidth
	scaleY := e.Height / imageHeight

	/*
		Très légère pulsation seulement.
		Le portail ne grossit presque plus.
	*/
	visualPulse := 1.0

	if e.Active {
		visualPulse = 1.0 + pulse*0.006
	}

	scaleX *= visualPulse
	scaleY *= visualPulse

	drawWidth := imageWidth * scaleX
	drawHeight := imageHeight * scaleY

	centerX := e.X + e.Width/2

	drawX := centerX - drawWidth/2
	drawY := groundY - drawHeight - e.GroundOffset

	e.Y = drawY

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(drawX, drawY)

	options.Filter = ebiten.FilterNearest

	if e.Active {
		brightness := float32(1.02 + pulse*0.05)

		options.ColorScale.Scale(
			brightness,
			brightness,
			1.08,
			1,
		)
	} else {
		options.ColorScale.Scale(
			0.55,
			0.55,
			0.65,
			1,
		)
	}

	screen.DrawImage(e.Image, options)
}

func (e *Exit) drawParticles(screen *ebiten.Image) {
	centerX := e.X + e.Width/2
	centerY := e.Y + e.Height/2

	timer := float64(e.AnimationTimer)

	for i := range e.Particles {
		particle := &e.Particles[i]

		/*
			Les particules restent très proches de la porte.
		*/
		angle := particle.Angle

		verticalWave := math.Sin(timer*0.035 + particle.Phase)

		x := centerX + math.Cos(angle)*particle.Radius
		y := centerY + math.Sin(angle)*particle.Radius*1.45 + verticalWave*3

		alphaPulse := (math.Sin(timer*0.08+particle.Phase) + 1) / 2

		alpha := uint8(90 + alphaPulse*100)

		ebitenutil.DrawRect(
			screen,
			x,
			y,
			particle.Size,
			particle.Size,
			color.RGBA{
				R: 95,
				G: 160,
				B: 255,
				A: alpha,
			},
		)
	}
}

func (e *Exit) drawGroundEnergy(screen *ebiten.Image) {
	centerX := e.X + e.Width/2
	ground := groundY - e.GroundOffset + 2

	timer := float64(e.AnimationTimer)

	for i := 0; i < 4; i++ {
		offset := math.Sin(timer*0.06+float64(i)*1.8) * 12

		x := centerX + offset

		alphaPulse := (math.Sin(timer*0.09+float64(i)) + 1) / 2

		alpha := uint8(70 + alphaPulse*90)

		ebitenutil.DrawRect(
			screen,
			x,
			ground,
			1,
			1,
			color.RGBA{
				R: 100,
				G: 175,
				B: 255,
				A: alpha,
			},
		)
	}
}
