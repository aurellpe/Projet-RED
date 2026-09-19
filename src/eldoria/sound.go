package main

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const soundSampleRate = 44100

type SoundManager struct {
	Context *audio.Context
	Volume  float64

	SwordData   []byte
	ImpactData  []byte
	HurtData    []byte
	BossData    []byte
	VictoryData []byte

	Players []*audio.Player
}

func NewSoundManager() *SoundManager {
	context := audio.NewContext(soundSampleRate)

	return &SoundManager{
		Context:     context,
		Volume:      0.70,
		SwordData:   generateSoundEffect(850, 320, 0.11, 0.65),
		ImpactData:  generateSoundEffect(180, 70, 0.13, 0.85),
		HurtData:    generateSoundEffect(300, 120, 0.20, 0.65),
		BossData:    generateSoundEffect(130, 55, 0.32, 0.80),
		VictoryData: generateSoundEffect(420, 900, 0.65, 0.65),
		Players:     []*audio.Player{},
	}
}

func generateSoundEffect(startFrequency float64, endFrequency float64, duration float64, volume float64) []byte {
	sampleCount := int(float64(soundSampleRate) * duration)

	if sampleCount <= 0 {
		return []byte{}
	}

	data := make([]byte, sampleCount*4)

	phase := 0.0

	for i := 0; i < sampleCount; i++ {
		progress := float64(i) / float64(sampleCount)

		frequency := startFrequency + (endFrequency-startFrequency)*progress

		phase += 2 * math.Pi * frequency / float64(soundSampleRate)

		envelope := 1 - progress
		envelope *= envelope

		sample := math.Sin(phase)
		sample += math.Sin(phase*2.03) * 0.20

		sample *= volume * envelope

		if sample > 1 {
			sample = 1
		}

		if sample < -1 {
			sample = -1
		}

		value := int16(sample * 32767)

		position := i * 4

		binary.LittleEndian.PutUint16(data[position:position+2], uint16(value))
		binary.LittleEndian.PutUint16(data[position+2:position+4], uint16(value))
	}

	return data
}

func (s *SoundManager) SetVolume(volume float64) {
	if volume < 0 {
		volume = 0
	}

	if volume > 1 {
		volume = 1
	}

	s.Volume = volume

	for _, player := range s.Players {
		if player != nil {
			player.SetVolume(s.Volume)
		}
	}
}

func (s *SoundManager) cleanPlayers() {
	activePlayers := []*audio.Player{}

	for _, player := range s.Players {
		if player != nil && player.IsPlaying() {
			activePlayers = append(activePlayers, player)
		}
	}

	s.Players = activePlayers
}

func (s *SoundManager) play(data []byte) {
	if s == nil || s.Context == nil {
		return
	}

	if len(data) == 0 {
		return
	}

	s.cleanPlayers()

	player := s.Context.NewPlayerFromBytes(data)

	if player == nil {
		return
	}

	player.SetVolume(s.Volume)
	player.Play()

	s.Players = append(s.Players, player)
}

func (s *SoundManager) PlaySword() {
	s.play(s.SwordData)
}

func (s *SoundManager) PlayImpact() {
	s.play(s.ImpactData)
}

func (s *SoundManager) PlayHurt() {
	s.play(s.HurtData)
}

func (s *SoundManager) PlayBossAttack() {
	s.play(s.BossData)
}

func (s *SoundManager) PlayVictory() {
	s.play(s.VictoryData)
}
