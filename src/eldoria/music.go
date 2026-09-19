package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed assets/music/music_level.mp3
var musicLevelData []byte

//go:embed assets/music/music_boss.mp3
var musicBossData []byte

const (
	musicTrackNone = iota
	musicTrackLevel
	musicTrackBoss
)

type MusicManager struct {
	Context *audio.Context

	LevelPlayer *audio.Player
	BossPlayer  *audio.Player

	CurrentTrack int
	Volume       float64
}

func NewMusicManager(context *audio.Context) *MusicManager {
	levelPlayer := newLoopingMusicPlayer(context, musicLevelData)
	bossPlayer := newLoopingMusicPlayer(context, musicBossData)

	manager := &MusicManager{
		Context:      context,
		LevelPlayer:  levelPlayer,
		BossPlayer:   bossPlayer,
		CurrentTrack: musicTrackNone,
		Volume:       0.20,
	}

	manager.SetVolume(manager.Volume)

	return manager
}

func newLoopingMusicPlayer(context *audio.Context, data []byte) *audio.Player {
	stream, err := mp3.DecodeWithSampleRate(soundSampleRate, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())

	player, err := context.NewPlayer(loop)
	if err != nil {
		panic(err)
	}

	return player
}

func (m *MusicManager) SetVolume(volume float64) {
	if volume < 0 {
		volume = 0
	}

	if volume > 1 {
		volume = 1
	}

	m.Volume = volume

	if m.LevelPlayer != nil {
		m.LevelPlayer.SetVolume(m.Volume)
	}

	if m.BossPlayer != nil {
		m.BossPlayer.SetVolume(m.Volume)
	}
}

func (m *MusicManager) PlayLevel() {
	m.playTrack(m.LevelPlayer, musicTrackLevel)
}

func (m *MusicManager) PlayBoss() {
	m.playTrack(m.BossPlayer, musicTrackBoss)
}

func (m *MusicManager) playTrack(player *audio.Player, track int) {
	if m == nil || m.Context == nil || player == nil {
		return
	}

	if !m.Context.IsReady() {
		return
	}

	if m.CurrentTrack == track && player.IsPlaying() {
		return
	}

	if m.LevelPlayer != nil {
		m.LevelPlayer.Pause()
	}

	if m.BossPlayer != nil {
		m.BossPlayer.Pause()
	}

	if m.CurrentTrack != track {
		if err := player.Rewind(); err != nil {
			return
		}
	}

	player.SetVolume(m.Volume)
	player.Play()

	m.CurrentTrack = track
}

func (m *MusicManager) Stop() {
	if m == nil {
		return
	}

	if m.LevelPlayer != nil {
		m.LevelPlayer.Pause()
	}

	if m.BossPlayer != nil {
		m.BossPlayer.Pause()
	}

	m.CurrentTrack = musicTrackNone
}
