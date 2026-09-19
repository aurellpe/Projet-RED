package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	State int
	Menu  *MainMenu

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

	Sounds *SoundManager
	Music  *MusicManager
	World  *ebiten.Image

	ShakeTimer    int
	ShakeStrength float64

	HitStopTimer int

	Transitioning   bool
	TransitionTimer int
	NextLevel       int

	HealMessageTimer int
	LastHealAmount   int

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

	g := &Game{
		State:              GameStateMenu,
		Menu:               NewMainMenu(),
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
		VictorySoundPlayed: false,
	}

	g.SetupCurrentLevel()

	return g
}

func (g *Game) UpdateMainMenu() error {
	if g.Music != nil {
		g.Music.Stop()
	}

	action := g.Menu.Update()

	switch action {
	case MenuActionPlay:
		g.RestartGame()
		g.State = GameStatePlaying

	case MenuActionOptions:
		g.State = GameStateOptions

	case MenuActionQuit:
		if g.Music != nil {
			g.Music.Stop()
		}

		os.Exit(0)
	}

	return nil
}

func (g *Game) UpdateOptionsMenu() error {
	if g.Music != nil {
		g.Music.Stop()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.State = GameStateMenu
	}

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

func (g *Game) SetupCurrentLevel() {
	g.SyncGround()

	g.Player.X = 15
	g.Player.Y = groundY - g.Player.Height
	g.Player.VelocityY = 0
	g.Player.OnGround = true

	g.Boss = nil
	g.Enemies = []*Enemy{}
	g.Exit = nil
	g.Pickups = []*Pickup{}
	g.DamageTexts = []*DamageText{}
	g.Impacts = []*Impact{}

	g.ShakeTimer = 0
	g.ShakeStrength = 0
	g.HitStopTimer = 0

	g.HealMessageTimer = 0
	g.LastHealAmount = 0
	g.VictorySoundPlayed = false

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		g.Boss = NewBoss()
		return
	}

	g.SpawnEnemiesForLevel()
	g.Exit = NewExit(groundY)
}

func (g *Game) SpawnEnemiesForLevel() {
	g.Enemies = []*Enemy{}

	level := g.Levels[g.CurrentLevel]

	if level.HasBoss {
		return
	}

	switch level.Number {
	case 1:
		g.Enemies = append(g.Enemies, NewEnemy(220))

	case 2:
		g.Enemies = append(g.Enemies, NewEnemy(180), NewEnemy(255))

	case 3:
		g.Enemies = append(g.Enemies, NewEnemy(140), NewEnemy(210), NewEnemy(270))
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

func (g *Game) PlayerAttackCanHit() bool {
	if !g.Player.Alive {
		return false
	}

	if !g.Player.Attacking {
		return false
	}

	if g.Player.AttackHit {
		return false
	}

	return g.Player.CurrentAttackFrame == 2
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
	g.CurrentLevel = 0
	g.Player = NewPlayer()

	g.Transitioning = false
	g.TransitionTimer = 0
	g.NextLevel = 0

	g.SetupCurrentLevel()
}

func (g *Game) RestartCurrentLevel() {
	g.Player = NewPlayer()

	g.Transitioning = false
	g.TransitionTimer = 0

	g.SetupCurrentLevel()
}

func (g *Game) StartLevelTransition() {
	if g.CurrentLevel >= len(g.Levels)-1 {
		return
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
	}

	if g.TransitionTimer >= 40 {
		g.Transitioning = false
		g.TransitionTimer = 0
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
		}
	}

	active := []*Pickup{}

	for _, pickup := range g.Pickups {
		if !pickup.Collected {
			active = append(active, pickup)
		}
	}

	g.Pickups = active
}

func (g *Game) UpdateImpacts() {
	active := []*Impact{}

	for _, impact := range g.Impacts {
		impact.Update()

		if impact.Alive() {
			active = append(active, impact)
		}
	}

	g.Impacts = active
}

func (g *Game) UpdateDamageTexts() {
	active := []*DamageText{}

	for _, damageText := range g.DamageTexts {
		damageText.Update()

		if damageText.Alive() {
			active = append(active, damageText)
		}
	}

	g.DamageTexts = active
}

func (g *Game) Update() error {
	if g.State == GameStateMenu {
		return g.UpdateMainMenu()
	}

	if g.State == GameStateOptions {
		return g.UpdateOptionsMenu()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.State = GameStateMenu

		if g.Music != nil {
			g.Music.Stop()
		}

		return nil
	}

	g.SyncGround()
	g.UpdateMusic()

	if g.HealMessageTimer > 0 {
		g.HealMessageTimer--
	}

	if g.Transitioning {
		g.UpdateTransition()
		return nil
	}

	if g.HitStopTimer > 0 {
		g.HitStopTimer--
		return nil
	}

	if !g.Player.Alive {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.RestartCurrentLevel()
		}

		g.UpdateDamageTexts()
		g.UpdateImpacts()

		return nil
	}

	if g.Boss != nil && !g.Boss.Alive {
		if !g.VictorySoundPlayed {
			g.Sounds.PlayVictory()
			g.VictorySoundPlayed = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.State = GameStateMenu
			g.Menu.Selected = 0

			if g.Music != nil {
				g.Music.Stop()
			}
		}

		g.UpdateDamageTexts()
		g.UpdateImpacts()

		return nil
	}

	wasPlayerAttacking := g.Player.Attacking

	g.Player.Update()

	if !wasPlayerAttacking && g.Player.Attacking {
		g.Sounds.PlaySword()
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
		return nil
	}

	g.UpdateNormalEnemies()
	g.UpdatePlayerAttackAgainstEnemies()
	g.UpdatePickups()
	g.UpdateExit()
	g.UpdateDamageTexts()
	g.UpdateImpacts()

	return nil
}

func (g *Game) UpdateBossLevel() {
	if g.Boss == nil {
		g.Boss = NewBoss()
	}

	if !g.Boss.Alive {
		return
	}

	wasBossAttacking := g.Boss.Attacking

	damageToPlayer := g.Boss.Update(g.Player)

	if !wasBossAttacking && g.Boss.Attacking {
		g.Sounds.PlayBossAttack()
	}

	if damageToPlayer > 0 {
		g.Sounds.PlayHurt()

		playerBox := g.Player.HitBox()

		g.Impacts = append(g.Impacts, NewImpact(playerBox.X+playerBox.W/2, playerBox.Y+playerBox.H/2, true))
		g.DamageTexts = append(g.DamageTexts, NewDamageText(g.Player.X, g.Player.Y-10, damageToPlayer))

		g.StartShake(10, 2.5)
	}

	if !g.Boss.Alive || g.Boss.Dying || !g.PlayerAttackCanHit() {
		return
	}

	playerAttack := g.Player.AttackBox()
	bossBox := g.Boss.HitBox()

	if !Intersects(playerAttack, bossBox) {
		return
	}

	playerCenter := g.Player.X + g.Player.Width/2

	if g.Boss.Hit(50, playerCenter) {
		g.Sounds.PlayImpact()

		g.Impacts = append(g.Impacts, NewImpact(bossBox.X+bossBox.W/2, bossBox.Y+bossBox.H/2, false))
		g.DamageTexts = append(g.DamageTexts, NewDamageText(g.Boss.X+g.Boss.HitWidth/2, groundY-g.Boss.HitHeight-10, 50))

		g.StartShake(7, 1.5)
		g.StartHitStop(4)
	}

	g.Player.AttackHit = true
}

func (g *Game) UpdateNormalEnemies() {
	for _, enemy := range g.Enemies {
		if !enemy.Alive {
			if enemy.Dying && !enemy.DeathRewarded {
				enemy.DeathRewarded = true
				g.Pickups = append(g.Pickups, NewHealthPotion(enemy.X+enemy.Width/2-6))
			}

			continue
		}

		damageToPlayer := enemy.Update(g.Player)

		if !enemy.Alive && enemy.Dying && !enemy.DeathRewarded {
			enemy.DeathRewarded = true
			g.Pickups = append(g.Pickups, NewHealthPotion(enemy.X+enemy.Width/2-6))
			continue
		}

		if damageToPlayer > 0 {
			g.Sounds.PlayHurt()

			playerBox := g.Player.HitBox()

			g.Impacts = append(g.Impacts, NewImpact(playerBox.X+playerBox.W/2, playerBox.Y+playerBox.H/2, true))
			g.DamageTexts = append(g.DamageTexts, NewDamageText(g.Player.X, g.Player.Y-10, damageToPlayer))

			g.StartShake(6, 1.2)
		}
	}
}

func (g *Game) UpdatePlayerAttackAgainstEnemies() {
	if !g.PlayerAttackCanHit() {
		return
	}

	playerAttack := g.Player.AttackBox()

	for _, enemy := range g.Enemies {
		if !enemy.Alive || enemy.Dying {
			continue
		}

		enemyBox := enemy.HitBox()

		if !Intersects(playerAttack, enemyBox) {
			continue
		}

		playerCenter := g.Player.X + g.Player.Width/2

		if enemy.Hit(50, playerCenter) {
			g.Sounds.PlayImpact()

			g.Impacts = append(g.Impacts, NewImpact(enemyBox.X+enemyBox.W/2, enemyBox.Y+enemyBox.H/2, false))
			g.DamageTexts = append(g.DamageTexts, NewDamageText(enemy.X+enemy.Width/2, enemy.Y-10, 50))

			g.StartShake(5, 1)
			g.StartHitStop(4)
		}

		g.Player.AttackHit = true
		break
	}
}

func (g *Game) UpdateExit() {
	if g.Exit == nil {
		return
	}

	g.Exit.Active = g.AllEnemiesDead()
	g.Exit.Update()

	if g.Exit.Active && g.Exit.PlayerInside(g.Player) && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.StartLevelTransition()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.State == GameStateMenu {
		g.Menu.Draw(screen)
		return
	}

	if g.State == GameStateOptions {
		DrawOptionsScreen(screen, g.Menu.Image)
		return
	}

	g.World.Clear()

	background := g.Backgrounds[g.CurrentLevel]

	bgOptions := &ebiten.DrawImageOptions{}

	width := float64(background.Bounds().Dx())
	height := float64(background.Bounds().Dy())

	bgOptions.GeoM.Scale(320/width, 180/height)

	g.World.DrawImage(background, bgOptions)

	if g.Exit != nil {
		g.Exit.Draw(g.World)
	}

	for _, pickup := range g.Pickups {
		pickup.Draw(g.World)
	}

	g.Player.Draw(g.World)

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

	worldOptions := &ebiten.DrawImageOptions{}

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

		worldOptions.GeoM.Translate(shakeX, shakeY)
	}

	screen.Fill(color.RGBA{A: 255})
	screen.DrawImage(g.World, worldOptions)

	g.DrawHUD(screen)
}

func (g *Game) DrawHUD(screen *ebiten.Image) {
	level := g.Levels[g.CurrentLevel]

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Niveau %d - %s", level.Number, level.Name))

	ebitenutil.DrawRect(screen, 5, 15, 80, 7, color.RGBA{R: 30, G: 30, B: 30, A: 255})

	playerRatio := float64(g.Player.HP) / float64(g.Player.MaxHP)

	if playerRatio < 0 {
		playerRatio = 0
	}

	if playerRatio > 1 {
		playerRatio = 1
	}

	ebitenutil.DrawRect(screen, 5, 15, 80*playerRatio, 7, color.RGBA{R: 30, G: 200, B: 60, A: 255})

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PV %d/%d", g.Player.HP, g.Player.MaxHP), 5, 25)

	if g.HealMessageTimer > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("+%d PV", g.LastHealAmount), 10, 40)
	}

	if !level.HasBoss {
		remaining := g.EnemiesRemaining()

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Ennemis : %d", remaining), 235, 5)

		if remaining == 0 {
			ebitenutil.DebugPrintAt(screen, "PORTE OUVERTE", 210, 18)
		}
	}

	if g.Exit != nil && g.Exit.PlayerInside(g.Player) {
		if g.Exit.Active {
			ebitenutil.DebugPrintAt(screen, "E = ENTRER", 125, 150)
		} else {
			ebitenutil.DebugPrintAt(screen, "PORTE VERROUILLEE", 100, 150)
		}
	}

	if g.Boss != nil && g.Boss.Alive {
		ebitenutil.DrawRect(screen, 100, 15, 150, 8, color.RGBA{R: 30, G: 30, B: 30, A: 255})

		bossRatio := float64(g.Boss.HP) / float64(g.Boss.MaxHP)

		if bossRatio < 0 {
			bossRatio = 0
		}

		if bossRatio > 1 {
			bossRatio = 1
		}

		ebitenutil.DrawRect(screen, 100, 15, 150*bossRatio, 8, color.RGBA{R: 220, G: 30, B: 30, A: 255})
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BOSS %d/%d", g.Boss.HP, g.Boss.MaxHP), 135, 27)
	}

	if !g.Player.Alive {
		ebitenutil.DrawRect(screen, 0, 0, 320, 180, color.RGBA{A: 180})

		ebitenutil.DebugPrintAt(screen, "GAME OVER", 125, 65)
		ebitenutil.DebugPrintAt(screen, "ENTREE = RECOMMENCER", 85, 85)
		ebitenutil.DebugPrintAt(screen, "ECHAP = MENU", 110, 100)

		return
	}

	if g.Boss != nil && !g.Boss.Alive {
		ebitenutil.DrawRect(screen, 0, 0, 320, 180, color.RGBA{A: 170})

		ebitenutil.DebugPrintAt(screen, "VICTOIRE !", 125, 60)
		ebitenutil.DebugPrintAt(screen, "BOSS VAINCU", 115, 75)
		ebitenutil.DebugPrintAt(screen, "ENTREE = MENU", 105, 95)
	}

	if g.Transitioning {
		alpha := 0.0

		if g.TransitionTimer <= 20 {
			alpha = float64(g.TransitionTimer) / 20
		} else {
			alpha = float64(40-g.TransitionTimer) / 20
		}

		if alpha < 0 {
			alpha = 0
		}

		if alpha > 1 {
			alpha = 1
		}

		ebitenutil.DrawRect(screen, 0, 0, 320, 180, color.RGBA{A: uint8(alpha * 255)})
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return 320, 180
}
