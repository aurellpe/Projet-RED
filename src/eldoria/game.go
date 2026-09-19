package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/background/level1.png
var level1BackgroundData []byte

//go:embed assets/background/level2.png
var level2BackgroundData []byte

//go:embed assets/background/level3.png
var level3BackgroundData []byte

//go:embed assets/background/boss.png
var bossBackgroundData []byte

type Game struct {
	State   int
	Menu    *MainMenu
	Options *OptionsMenu

	Paused           bool
	PauseMenu        *PauseMenu
	OptionsFromPause bool

	Player *Player

	Levels       []Level
	CurrentLevel int

	Backgrounds []*ebiten.Image

	Boss    *Boss
	Enemies []*Enemy
	Exit    *Exit
	Pickups []*Pickup

	DamageTexts []*DamageText
	Impacts     []*Impact

	Sounds     *SoundManager
	Music      *MusicManager
	World      *ebiten.Image
	Atmosphere *Atmosphere
	Story      *StoryManager
	Combat     *TurnCombat
	Puzzle     *DoorPuzzle
	SkillTree  *SkillTree

	ShakeTimer    int
	ShakeStrength float64

	HitStopTimer int

	Transitioning   bool
	TransitionTimer int
	NextLevel       int

	HealMessageTimer int
	LastHealAmount   int

	AutosaveTimer int

	UpgradeMessage      string
	UpgradeMessageTimer int

	VictorySoundPlayed bool
}

func loadBackground(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func NewGame() *Game {
	levels := CreateLevels()

	groundY = levels[0].GroundY

	backgrounds := []*ebiten.Image{
		loadBackground(level1BackgroundData),
		loadBackground(level2BackgroundData),
		loadBackground(level3BackgroundData),
		loadBackground(bossBackgroundData),
	}

	sounds := NewSoundManager()
	music := NewMusicManager(sounds.Context)
	options := NewOptionsMenu(music, sounds)

	game := &Game{
		State:              GameStateMenu,
		Menu:               NewMainMenu(),
		Options:            options,
		Paused:             false,
		PauseMenu:          NewPauseMenu(),
		OptionsFromPause:   false,
		Player:             NewPlayer(),
		Levels:             levels,
		CurrentLevel:       0,
		Backgrounds:        backgrounds,
		Enemies:            []*Enemy{},
		Pickups:            []*Pickup{},
		DamageTexts:        []*DamageText{},
		Impacts:            []*Impact{},
		Sounds:             sounds,
		Music:              music,
		World:              ebiten.NewImage(320, 180),
		Atmosphere:         NewAtmosphere(),
		Story:              NewStoryManager(),
		Combat:             NewTurnCombat(),
		Puzzle:             nil,
		SkillTree:          LoadSkillTreeProgress(),
		AutosaveTimer:      0,
		VictorySoundPlayed: false,
	}

	game.SetupCurrentLevel()

	return game
}

func (g *Game) UpdateMainMenu() error {
	if g.Music != nil {
		g.Music.Stop()
	}

	g.Paused = false
	g.OptionsFromPause = false

	if g.SkillTree != nil {
		g.SkillTree.Open = false
	}

	action := g.Menu.Update()

	switch action {
	case MenuActionContinue:
		g.ContinueGame()

	case MenuActionPlay:
		g.RestartGame()
		g.State = GameStatePlaying

	case MenuActionOptions:
		g.OptionsFromPause = false
		g.State = GameStateOptions

	case MenuActionQuit:
		if g.Music != nil {
			g.Music.Stop()
		}

		os.Exit(0)
	}

	return nil
}

func (g *Game) ContinueGame() {
	save, ok := LoadSave()

	if !ok {
		return
	}

	if save.CurrentLevel < 0 || save.CurrentLevel >= len(g.Levels) {
		return
	}

	g.CurrentLevel = save.CurrentLevel
	g.Player = NewPlayer()

	g.SkillTree = LoadSkillTreeProgress()
	g.SkillTree.GrantCatchUpRewards(g.CurrentLevel)

	_ = SaveSkillTreeProgress(g.SkillTree)

	g.Paused = false
	g.OptionsFromPause = false

	g.Transitioning = false
	g.TransitionTimer = 0
	g.NextLevel = g.CurrentLevel

	if g.Combat != nil {
		g.Combat.Reset()
	}

	if g.Story != nil {
		g.Story.Stop()
		g.Story.IntroSeen = true

		if g.Levels[g.CurrentLevel].HasBoss {
			g.Story.BossSeen = true
		}
	}

	g.SetupCurrentLevel()

	g.Player.HP = save.PlayerHP

	if g.Player.HP <= 0 {
		g.Player.HP = 1
	}

	if g.Player.HP > g.Player.MaxHP {
		g.Player.HP = g.Player.MaxHP
	}

	g.State = GameStatePlaying
}

func (g *Game) UpdatePauseMenu() error {
	action := g.PauseMenu.Update()

	switch action {
	case PauseActionResume:
		g.Paused = false

	case PauseActionOptions:
		g.OptionsFromPause = true
		g.Paused = false
		g.State = GameStateOptions

	case PauseActionMenu:
		if g.Player != nil && g.Player.Alive {
			_ = SaveGame(g)
		}

		if g.SkillTree != nil {
			g.SkillTree.Open = false
			_ = SaveSkillTreeProgress(g.SkillTree)
		}

		if g.Combat != nil {
			g.Combat.Reset()
		}

		if g.Puzzle != nil {
			g.Puzzle.Active = false
		}

		g.Paused = false
		g.OptionsFromPause = false
		g.State = GameStateMenu

		if HasSave() {
			g.Menu.Selected = 0
		} else {
			g.Menu.Selected = 1
		}

		if g.Music != nil {
			g.Music.Stop()
		}
	}

	return nil
}

func (g *Game) UpdateOptionsMenu() error {
	action := g.Options.Update(g.Music, g.Sounds)

	if action != OptionsActionBack {
		return nil
	}

	if g.OptionsFromPause {
		g.State = GameStatePlaying
		g.Paused = true
		g.OptionsFromPause = false
		return nil
	}

	g.State = GameStateMenu

	return nil
}

func (g *Game) SyncGround() {
	groundY = g.Levels[g.CurrentLevel].GroundY
}

func (g *Game) UpdateMusic() {
	if g.Music == nil {
		return
	}

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		if g.Boss != nil && !g.Boss.Alive {
			g.Music.Stop()
			return
		}

		g.Music.PlayBoss()
		return
	}

	g.Music.PlayLevel()
}

func (g *Game) ApplyCurrentLevelProgression() {
	if g.Player == nil {
		return
	}

	levelIndex := g.CurrentLevel

	targetMaxHP := 100 + levelIndex*10
	baseAttack := 50 + levelIndex*5
	targetSpeed := 1.5 + float64(levelIndex)*0.08
	targetDodgeSpeed := 4.0 + float64(levelIndex)*0.15
	targetDodgeCooldown := 38 - levelIndex*3

	if g.SkillTree != nil {
		targetMaxHP += g.SkillTree.MaxHPBonus()
		baseAttack = g.SkillTree.ApplyAttackBonus(baseAttack)
	}

	if targetDodgeCooldown < 28 {
		targetDodgeCooldown = 28
	}

	oldMaxHP := g.Player.MaxHP

	g.Player.MaxHP = targetMaxHP

	if targetMaxHP > oldMaxHP {
		g.Player.HP += targetMaxHP - oldMaxHP
	}

	if g.Player.HP > g.Player.MaxHP {
		g.Player.HP = g.Player.MaxHP
	}

	g.Player.AttackDamage = baseAttack
	g.Player.Speed = targetSpeed
	g.Player.DodgeSpeed = targetDodgeSpeed
	g.Player.DodgeCooldown = targetDodgeCooldown
}

func (g *Game) ShowUpgradeMessage() {
	if g.CurrentLevel <= 0 {
		return
	}

	g.UpgradeMessage = "AMELIORATION : PV + ATTAQUE + VITESSE   |   C : COMPETENCES"
	g.UpgradeMessageTimer = 210
}

func (g *Game) UpdateAutoSave() {
	if g.State != GameStatePlaying {
		return
	}

	if g.Paused || g.Transitioning {
		return
	}

	if g.Combat != nil && g.Combat.Active {
		return
	}

	if g.Puzzle != nil && g.Puzzle.Active {
		return
	}

	if g.SkillTree != nil && g.SkillTree.Open {
		return
	}

	if g.Player == nil || !g.Player.Alive {
		return
	}

	if g.Story != nil && g.Story.Active {
		return
	}

	g.AutosaveTimer++

	if g.AutosaveTimer >= 120 {
		_ = SaveGame(g)
		_ = SaveSkillTreeProgress(g.SkillTree)
		g.AutosaveTimer = 0
	}
}

func (g *Game) UpdateAtmosphere() {
	if g.Atmosphere == nil {
		return
	}

	level := g.Levels[g.CurrentLevel]

	g.Atmosphere.Update(level.Number)
}

func (g *Game) SetupCurrentLevel() {
	g.SyncGround()

	g.Player.X = 15
	g.Player.Y = groundY - g.Player.Height
	g.Player.VelocityY = 0
	g.Player.OnGround = true

	g.Player.Dodging = false
	g.Player.DodgeTimer = 0
	g.Player.DodgeCooldownTimer = 0

	g.ApplyCurrentLevelProgression()

	if g.Combat != nil {
		g.Combat.Reset()
	}

	if g.SkillTree != nil {
		g.SkillTree.Open = false
	}

	g.Boss = nil
	g.Enemies = []*Enemy{}
	g.Exit = nil
	g.Pickups = []*Pickup{}
	g.Puzzle = nil

	g.DamageTexts = []*DamageText{}
	g.Impacts = []*Impact{}

	g.ShakeTimer = 0
	g.ShakeStrength = 0
	g.HitStopTimer = 0

	g.HealMessageTimer = 0
	g.LastHealAmount = 0

	g.AutosaveTimer = 0

	g.UpgradeMessage = ""
	g.UpgradeMessageTimer = 0

	g.VictorySoundPlayed = false

	level := g.Levels[g.CurrentLevel]

	if g.Atmosphere != nil {
		g.Atmosphere.Reset(level.Number)
	}

	if level.HasBoss {
		g.Boss = NewBoss()
		return
	}

	g.SpawnEnemiesForLevel()

	g.Exit = NewExit(groundY)
	g.Puzzle = NewDoorPuzzle(level.Number)
}

func (g *Game) SpawnEnemiesForLevel() {
	g.Enemies = []*Enemy{}

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		return
	}

	switch level.Number {
	case 1:
		g.Enemies = append(
			g.Enemies,
			NewEnemy(220),
		)

	case 2:
		g.Enemies = append(
			g.Enemies,
			NewFastEnemy(180),
			NewEnemy(255),
		)

	case 3:
		g.Enemies = append(
			g.Enemies,
			NewHeavyEnemy(145),
			NewRangedEnemy(220),
			NewFastEnemy(275),
		)
	}
}

func (g *Game) EnemiesRemaining() int {
	count := 0

	for _, enemy := range g.Enemies {
		if enemy.Alive {
			count++
		}
	}

	return count
}

func (g *Game) AllEnemiesDead() bool {
	return g.EnemiesRemaining() == 0
}

func (g *Game) StartHitStop(frames int) {
	if frames > g.HitStopTimer {
		g.HitStopTimer = frames
	}
}

func (g *Game) StartShake(duration int, strength float64) {
	g.ShakeTimer = duration
	g.ShakeStrength = strength
}

func (g *Game) RestartGame() {
	_ = DeleteSave()
	_ = DeleteSkillTreeProgress()

	g.CurrentLevel = 0
	g.Player = NewPlayer()
	g.SkillTree = NewSkillTree()

	g.Paused = false
	g.OptionsFromPause = false

	g.Transitioning = false
	g.TransitionTimer = 0
	g.NextLevel = 0

	g.AutosaveTimer = 0

	if g.Combat != nil {
		g.Combat.Reset()
	}

	if g.Story != nil {
		g.Story.Reset()
	}

	g.SetupCurrentLevel()

	if g.Story != nil {
		g.Story.StartIntro()
	}
}

func (g *Game) RestartCurrentLevel() {
	g.Player = NewPlayer()

	g.Paused = false
	g.OptionsFromPause = false

	g.Transitioning = false
	g.TransitionTimer = 0

	g.AutosaveTimer = 0

	if g.Combat != nil {
		g.Combat.Reset()
	}

	g.SetupCurrentLevel()
}

func (g *Game) StartLevelTransition() {
	if g.CurrentLevel >= len(g.Levels)-1 {
		return
	}

	if g.SkillTree != nil {
		g.SkillTree.AwardLevelCompletion(g.CurrentLevel)
	}

	g.NextLevel = g.CurrentLevel + 1
	g.Transitioning = true
	g.TransitionTimer = 0
}

func (g *Game) UpdateTransition() {
	g.TransitionTimer++

	if g.TransitionTimer == 20 {
		g.CurrentLevel = g.NextLevel

		g.SetupCurrentLevel()
		g.ShowUpgradeMessage()

		_ = SaveGame(g)
		_ = SaveSkillTreeProgress(g.SkillTree)
	}

	if g.TransitionTimer >= 40 {
		g.Transitioning = false
		g.TransitionTimer = 0

		if g.Story != nil && g.Levels[g.CurrentLevel].HasBoss {
			g.Story.StartBossIntro()
		}
	}
}

func (g *Game) UpdatePickups() {
	for _, pickup := range g.Pickups {
		if pickup.Collected {
			continue
		}

		healed := pickup.Update(g.Player)

		if healed > 0 {
			g.LastHealAmount = healed
			g.HealMessageTimer = 60

			_ = SaveGame(g)
			_ = SaveSkillTreeProgress(g.SkillTree)

			g.AutosaveTimer = 0
		}
	}

	activePickups := []*Pickup{}

	for _, pickup := range g.Pickups {
		if !pickup.Collected {
			activePickups = append(activePickups, pickup)
		}
	}

	g.Pickups = activePickups
}

func (g *Game) UpdateImpacts() {
	activeImpacts := []*Impact{}

	for _, impact := range g.Impacts {
		impact.Update()

		if impact.Alive() {
			activeImpacts = append(activeImpacts, impact)
		}
	}

	g.Impacts = activeImpacts
}

func (g *Game) UpdateDamageTexts() {
	activeTexts := []*DamageText{}

	for _, damageText := range g.DamageTexts {
		damageText.Update()

		if damageText.Alive() {
			activeTexts = append(activeTexts, damageText)
		}
	}

	g.DamageTexts = activeTexts
}

func (g *Game) UpdateInterfaceTimers() {
	if g.HealMessageTimer > 0 {
		g.HealMessageTimer--
	}

	if g.UpgradeMessageTimer > 0 {
		g.UpgradeMessageTimer--
	}
}

func (g *Game) Update() error {
	if g.State == GameStateMenu {
		return g.UpdateMainMenu()
	}

	if g.State == GameStateOptions {
		return g.UpdateOptionsMenu()
	}

	g.SyncGround()
	g.UpdateMusic()

	if g.Paused {
		return g.UpdatePauseMenu()
	}

	if g.Player == nil {
		return nil
	}

	g.UpdateInterfaceTimers()

	if g.Story != nil && g.Story.Active {
		g.UpdateAtmosphere()
		g.Story.Update()
		return nil
	}

	if g.Puzzle != nil && g.Puzzle.Active {
		g.UpdateAtmosphere()
		g.Puzzle.Update()

		if g.Exit != nil && g.Puzzle.Solved {
			g.Exit.Active = true
		}

		return nil
	}

	if !g.Player.Alive {
		g.UpdateAtmosphere()

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.RestartCurrentLevel()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = GameStateMenu

			if HasSave() {
				g.Menu.Selected = 0
			} else {
				g.Menu.Selected = 1
			}

			if g.Music != nil {
				g.Music.Stop()
			}
		}

		g.UpdateDamageTexts()
		g.UpdateImpacts()

		return nil
	}

	if g.Boss != nil && !g.Boss.Alive {
		g.UpdateAtmosphere()

		if !g.VictorySoundPlayed {
			g.Sounds.PlayVictory()
			g.VictorySoundPlayed = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = GameStateMenu

			if HasSave() {
				g.Menu.Selected = 0
			} else {
				g.Menu.Selected = 1
			}

			if g.Music != nil {
				g.Music.Stop()
			}
		}

		g.UpdateDamageTexts()
		g.UpdateImpacts()

		return nil
	}

	if g.SkillTree != nil && g.SkillTree.Open {
		g.UpdateAtmosphere()
		g.SkillTree.Update(g)
		return nil
	}

	if !g.Transitioning && (g.Combat == nil || !g.Combat.Active) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if g.SkillTree != nil {
			g.SkillTree.OpenMenu()
		}

		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Paused = true
		g.PauseMenu.Selected = 0
		return nil
	}

	if g.Combat != nil && g.Combat.Active {
		g.UpdateAtmosphere()

		g.Combat.Update(g)

		g.UpdateDamageTexts()
		g.UpdateImpacts()

		if g.ShakeTimer > 0 {
			g.ShakeTimer--
		}

		return nil
	}

	if g.Transitioning {
		g.UpdateAtmosphere()
		g.UpdateTransition()
		return nil
	}

	g.UpdateAtmosphere()

	wasPlayerAttacking := g.Player.Attacking

	g.Player.Update()

	if !wasPlayerAttacking && g.Player.Attacking {
		g.Player.Attacking = false
		g.Player.AttackTimer = 0
		g.Player.CurrentAttackFrame = 0
		g.Player.AttackHit = false
	}

	if g.ShakeTimer > 0 {
		g.ShakeTimer--
	}

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		g.UpdateBossLevel()
		g.UpdatePickups()
		g.UpdateDamageTexts()
		g.UpdateImpacts()
		g.UpdateAutoSave()

		return nil
	}

	g.UpdateNormalEnemies()

	if g.Combat != nil && g.Combat.Active {
		return nil
	}

	g.UpdatePickups()
	g.UpdateExit()
	g.UpdateDamageTexts()
	g.UpdateImpacts()
	g.UpdateAutoSave()

	return nil
}

func (g *Game) UpdateBossLevel() {
	if g.Boss == nil {
		g.Boss = NewBoss()
	}

	if g.Boss.Dying {
		g.Boss.Update(g.Player)
		return
	}

	if !g.Boss.Alive {
		return
	}

	if g.Combat != nil && !g.Combat.Active {
		g.Combat.StartBoss(
			g.Boss,
			g.Player,
			g.SkillTree,
		)
	}
}

func (g *Game) UpdateNormalEnemies() {
	for _, enemy := range g.Enemies {
		if enemy.Dying {
			enemy.Update(g.Player)

			if !enemy.Alive && !enemy.DeathRewarded {
				enemy.DeathRewarded = true

				g.Pickups = append(
					g.Pickups,
					NewHealthPotion(enemy.X+enemy.Width/2-6),
				)
			}

			continue
		}

		if !enemy.Alive {
			if !enemy.DeathRewarded {
				enemy.DeathRewarded = true

				g.Pickups = append(
					g.Pickups,
					NewHealthPotion(enemy.X+enemy.Width/2-6),
				)
			}

			continue
		}

		enemy.Moving = false
		enemy.Attacking = false
		enemy.CurrentFrame = 0
		enemy.ProjectileActive = false

		playerCenter := g.Player.X + g.Player.Width/2
		enemyCenter := enemy.X + enemy.Width/2

		distance := math.Abs(playerCenter - enemyCenter)

		if distance <= 46 {
			if g.Combat != nil {
				g.Combat.StartEnemy(
					enemy,
					g.Player,
					g.SkillTree,
				)
			}

			return
		}
	}
}

func (g *Game) UpdateExit() {
	if g.Exit == nil {
		return
	}

	allEnemiesDead := g.AllEnemiesDead()

	if !allEnemiesDead {
		g.Exit.Active = false
		g.Exit.Update()
		return
	}

	if g.Puzzle == nil {
		g.Exit.Active = true
	} else {
		g.Exit.Active = g.Puzzle.Solved
	}

	g.Exit.Update()

	if !g.Exit.PlayerInside(g.Player) {
		return
	}

	if g.Puzzle != nil && !g.Puzzle.Solved {
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			g.Puzzle.Start()
		}

		return
	}

	if g.Exit.Active && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.StartLevelTransition()
	}
}

func (g *Game) drawWorldToScreen(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}

	scaleX := float64(renderWidth) / 320
	scaleY := float64(renderHeight) / 180

	options.GeoM.Scale(scaleX, scaleY)
	options.Filter = ebiten.FilterNearest

	if g.ShakeTimer > 0 {
		var shakeX float64
		var shakeY float64

		switch g.ShakeTimer % 4 {
		case 0:
			shakeX = g.ShakeStrength

		case 1:
			shakeX = -g.ShakeStrength
			shakeY = g.ShakeStrength

		case 2:
			shakeY = -g.ShakeStrength

		case 3:
			shakeX = g.ShakeStrength
			shakeY = -g.ShakeStrength
		}

		options.GeoM.Translate(shakeX*scaleX, shakeY*scaleY)
	}

	screen.Fill(
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		},
	)

	screen.DrawImage(g.World, options)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Clear()

	if g.State == GameStateMenu {
		g.Menu.Draw(g.World)
		g.drawWorldToScreen(screen)
		return
	}

	if g.State == GameStateOptions {
		g.Options.Draw(
			g.World,
			g.Menu.Image,
			g.OptionsFromPause,
		)

		g.drawWorldToScreen(screen)

		return
	}

	background := g.Backgrounds[g.CurrentLevel]

	bgOptions := &ebiten.DrawImageOptions{}

	backgroundWidth := float64(background.Bounds().Dx())
	backgroundHeight := float64(background.Bounds().Dy())

	bgOptions.GeoM.Scale(
		320/backgroundWidth,
		180/backgroundHeight,
	)

	g.World.DrawImage(
		background,
		bgOptions,
	)

	if g.Atmosphere != nil {
		g.Atmosphere.Draw(g.World)
	}

	if g.Exit != nil {
		g.Exit.Draw(g.World)
	}

	for _, pickup := range g.Pickups {
		pickup.Draw(g.World)
	}

	if g.Player != nil {
		g.Player.Draw(g.World)
	}

	for _, enemy := range g.Enemies {
		enemy.Draw(g.World)
	}

	if g.Boss != nil {
		g.Boss.Draw(g.World)
	}

	for _, impact := range g.Impacts {
		impact.Draw(g.World)
	}

	for _, damageText := range g.DamageTexts {
		damageText.Draw(g.World)
	}

	if g.Paused {
		g.PauseMenu.Draw(g.World)
	}

	g.drawWorldToScreen(screen)

	if g.Story != nil && g.Story.Active {
		g.Story.Draw(screen)
		return
	}

	if !g.Paused {
		g.DrawHUD(screen)
	}

	if g.Combat != nil && g.Combat.Active && !g.Paused {
		g.Combat.Draw(screen)
	}

	if g.Puzzle != nil && g.Puzzle.Active && !g.Paused {
		g.Puzzle.Draw(screen)
	}

	if g.SkillTree != nil && g.SkillTree.Open && !g.Paused {
		g.SkillTree.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return renderWidth, renderHeight
}
