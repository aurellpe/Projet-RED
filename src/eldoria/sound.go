package main

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const soundSampleRate = 44100

type SoundManager struct {
	Context *audio.Context

	SwordSound   []byte
	ImpactSound  []byte
	HurtSound    []byte
	BossSound    []byte
	VictorySound []byte
}

func NewSoundManager() *SoundManager {
	context := audio.NewContext(soundSampleRate)

	return &SoundManager{
		Context: context,

		SwordSound: makeSwordSound(),

		ImpactSound: makeImpactSound(),

		HurtSound: makeHurtSound(),

		BossSound: makeBossSound(),

		VictorySound: makeVictorySound(),
	}
}

func (s *SoundManager) play(data []byte) {
	if s == nil {
		return
	}

	if s.Context == nil {
		return
	}

	if len(data) == 0 {
		return
	}

	player := s.Context.NewPlayerFromBytes(data)

	player.Play()
}

func (s *SoundManager) PlaySword() {
	s.play(s.SwordSound)
}

func (s *SoundManager) PlayImpact() {
	s.play(s.ImpactSound)
}

func (s *SoundManager) PlayHurt() {
	s.play(s.HurtSound)
}

func (s *SoundManager) PlayBossAttack() {
	s.play(s.BossSound)
}

func (s *SoundManager) PlayVictory() {
	s.play(s.VictorySound)
}

// --------------------------------------------------
// AJOUT D'UN SAMPLE STEREO
// --------------------------------------------------

func appendStereoSample(
	buffer []byte,
	value float64,
) []byte {

	if value > 1 {
		value = 1
	}

	if value < -1 {
		value = -1
	}

	sample := int16(
		value * 32767,
	)

	temp := make([]byte, 4)

	binary.LittleEndian.PutUint16(
		temp[0:2],
		uint16(sample),
	)

	binary.LittleEndian.PutUint16(
		temp[2:4],
		uint16(sample),
	)

	return append(
		buffer,
		temp...,
	)
}

// --------------------------------------------------
// PETIT GENERATEUR DE BRUIT
// --------------------------------------------------

func nextNoise(seed *uint32) float64 {
	*seed = *seed*1664525 + 1013904223

	value := float64(
		(*seed>>16)&0xFFFF,
	) / 65535.0

	return value*2 - 1
}

// --------------------------------------------------
// SON D'EPEE
// --------------------------------------------------

func makeSwordSound() []byte {
	duration := 0.14

	totalSamples := int(
		float64(soundSampleRate) * duration,
	)

	buffer := make(
		[]byte,
		0,
		totalSamples*4,
	)

	seed := uint32(12345)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) /
			float64(soundSampleRate)

		progress := float64(i) /
			float64(totalSamples)

		frequency := 1500 -
			progress*1000

		envelope := 1 - progress

		tone := math.Sin(
			2 *
				math.Pi *
				frequency *
				t,
		)

		noise := nextNoise(&seed)

		value := tone*0.22 +
			noise*0.18

		value *= envelope

		buffer = appendStereoSample(
			buffer,
			value,
		)
	}

	return buffer
}

// --------------------------------------------------
// IMPACT D'EPEE
// --------------------------------------------------

func makeImpactSound() []byte {
	duration := 0.11

	totalSamples := int(
		float64(soundSampleRate) * duration,
	)

	buffer := make(
		[]byte,
		0,
		totalSamples*4,
	)

	seed := uint32(99991)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) /
			float64(soundSampleRate)

		progress := float64(i) /
			float64(totalSamples)

		envelope := math.Pow(
			1-progress,
			2,
		)

		lowTone := math.Sin(
			2 *
				math.Pi *
				110 *
				t,
		)

		metalTone := math.Sin(
			2 *
				math.Pi *
				750 *
				t,
		)

		noise := nextNoise(&seed)

		value := lowTone*0.40 +
			metalTone*0.18 +
			noise*0.25

		value *= envelope

		buffer = appendStereoSample(
			buffer,
			value,
		)
	}

	return buffer
}

// --------------------------------------------------
// JOUEUR BLESSE
// --------------------------------------------------

func makeHurtSound() []byte {
	duration := 0.18

	totalSamples := int(
		float64(soundSampleRate) * duration,
	)

	buffer := make(
		[]byte,
		0,
		totalSamples*4,
	)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) /
			float64(soundSampleRate)

		progress := float64(i) /
			float64(totalSamples)

		frequency := 240 -
			progress*130

		envelope := 1 - progress

		tone := math.Sin(
			2 *
				math.Pi *
				frequency *
				t,
		)

		secondTone := math.Sin(
			2 *
				math.Pi *
				frequency *
				2 *
				t,
		)

		value := tone*0.35 +
			secondTone*0.12

		value *= envelope

		buffer = appendStereoSample(
			buffer,
			value,
		)
	}

	return buffer
}

// --------------------------------------------------
// SON D'ATTAQUE DU BOSS
// --------------------------------------------------

func makeBossSound() []byte {
	duration := 0.32

	totalSamples := int(
		float64(soundSampleRate) * duration,
	)

	buffer := make(
		[]byte,
		0,
		totalSamples*4,
	)

	seed := uint32(7777)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) /
			float64(soundSampleRate)

		progress := float64(i) /
			float64(totalSamples)

		frequency := 85 +
			progress*70

		envelope := math.Sin(
			progress * math.Pi,
		)

		low := math.Sin(
			2 *
				math.Pi *
				frequency *
				t,
		)

		second := math.Sin(
			2 *
				math.Pi *
				frequency *
				0.5 *
				t,
		)

		noise := nextNoise(&seed)

		value := low*0.38 +
			second*0.22 +
			noise*0.08

		value *= envelope

		buffer = appendStereoSample(
			buffer,
			value,
		)
	}

	return buffer
}

// --------------------------------------------------
// VICTOIRE
// --------------------------------------------------

func makeVictorySound() []byte {
	notes := []float64{
		523.25,
		659.25,
		783.99,
		1046.50,
	}

	noteDuration := 0.14

	samplesPerNote := int(
		float64(soundSampleRate) *
			noteDuration,
	)

	totalSamples := samplesPerNote *
		len(notes)

	buffer := make(
		[]byte,
		0,
		totalSamples*4,
	)

	for _, frequency := range notes {
		for i := 0; i < samplesPerNote; i++ {
			t := float64(i) /
				float64(soundSampleRate)

			progress := float64(i) /
				float64(samplesPerNote)

			envelope := 1 - progress*0.65

			mainTone := math.Sin(
				2 *
					math.Pi *
					frequency *
					t,
			)

			harmonic := math.Sin(
				2 *
					math.Pi *
					frequency *
					2 *
					t,
			)

			value := mainTone*0.30 +
				harmonic*0.10

			value *= envelope

			buffer = appendStereoSample(
				buffer,
				value,
			)
		}
	}

	return buffer
}
