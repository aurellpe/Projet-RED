package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/exit/portal.png
var portalData []byte

type Exit struct {
	X float64
	Y float64

	Width  float64
	Height float64

	Active bool

	AnimationTimer int

	Image *ebiten.Image

	GroundOffset float64
}

func cropPortalImage(src image.Image) image.Image {
	bounds := src.Bounds()

	minX := bounds.Max.X
	minY := bounds.Max.Y
	maxX := bounds.Min.X
	maxY := bounds.Min.Y

	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, alpha := src.At(x, y).RGBA()

			if alpha > 1000 {
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
	}

	if !found {
		return src
	}

	cropRect := image.Rect(minX, minY, maxX+1, maxY+1)

	result := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))

	draw.Draw(result, result.Bounds(), src, cropRect.Min, draw.Src)

	return result
}

func loadPortalImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	cropped := cropPortalImage(img)

	return ebiten.NewImageFromImage(cropped)
}

func NewExit(levelGroundY float64) *Exit {
	width := 68.0
	height := 96.0
	groundOffset := 20.0

	return &Exit{
		X:              246,
		Y:              levelGroundY - height - groundOffset,
		Width:          width,
		Height:         height,
		Active:         false,
		AnimationTimer: 0,
		Image:          loadPortalImage(portalData),
		GroundOffset:   groundOffset,
	}
}

func (e *Exit) Update() {
	if e.Active {
		e.AnimationTimer++

		if e.AnimationTimer >= 120 {
			e.AnimationTimer = 0
		}
	} else {
		e.AnimationTimer = 0
	}

	e.Y = groundY - e.Height - e.GroundOffset
}

func (e *Exit) HitBox() Rect {
	interactionWidth := 36.0
	interactionHeight := 65.0

	centerX := e.X + e.Width/2
	visualGroundY := groundY - e.GroundOffset

	return Rect{
		X: centerX - interactionWidth/2,
		Y: visualGroundY - interactionHeight,
		W: interactionWidth,
		H: interactionHeight,
	}
}

func (e *Exit) PlayerInside(player *Player) bool {
	return Intersects(e.HitBox(), player.HitBox())
}

func (e *Exit) Draw(screen *ebiten.Image) {
	if e.Image == nil {
		return
	}

	imageWidth := float64(e.Image.Bounds().Dx())
	imageHeight := float64(e.Image.Bounds().Dy())

	scaleX := e.Width / imageWidth
	scaleY := e.Height / imageHeight

	scale := scaleX

	if scaleY < scale {
		scale = scaleY
	}

	drawWidth := imageWidth * scale
	drawHeight := imageHeight * scale

	visualGroundY := groundY - e.GroundOffset

	drawX := e.X + e.Width/2 - drawWidth/2
	drawY := visualGroundY - drawHeight

	options := &ebiten.DrawImageOptions{}

	if e.Active {
		pulse := 1.0

		if e.AnimationTimer < 30 {
			pulse = 1.00
		} else if e.AnimationTimer < 60 {
			pulse = 1.03
		} else if e.AnimationTimer < 90 {
			pulse = 1.05
		} else {
			pulse = 1.03
		}

		centerX := drawX + drawWidth/2
		centerY := drawY + drawHeight/2

		options.GeoM.Translate(-imageWidth/2, -imageHeight/2)
		options.GeoM.Scale(scale*pulse, scale*pulse)
		options.GeoM.Translate(centerX, centerY)

		options.ColorScale.Scale(1.15, 1.15, 1.3, 1)
	} else {
		options.GeoM.Scale(scale, scale)
		options.GeoM.Translate(drawX, drawY)

		options.ColorScale.Scale(0.45, 0.45, 0.55, 0.75)
	}

	screen.DrawImage(e.Image, options)
}