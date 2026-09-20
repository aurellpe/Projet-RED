package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	CombatTurnPlayer = iota + 1
	CombatTurnMonster
)

const (
	DungeonAttackSword = iota
	DungeonAttackMagic
	DungeonAttackFree
)

const (
	DungeonActionSword = iota
	DungeonActionRay
)

const (
	DungeonPendingNone = iota
	DungeonPendingSword
	DungeonPendingRay
	DungeonPendingMonster
)

const TrainingRayonManaCost = 40

type TrainingDungeon struct {
	Active bool

	Selecting bool

	SelectionRow         int
	SelectedEnemyType    int
	SelectedAttackSystem int

	SelectedCombatAction int

	Enemy *Enemy

	CombatTurn int
	TurnNumber int

	ActionInProgress bool
	ActionTimer      int
	DamageApplied    bool
	PendingAction    int

	Finished bool
	Won      bool

	Message      string
	MessageTimer int
}

func NewTrainingDungeon() *TrainingDungeon {
	return &TrainingDungeon{
		Active:               false,
		Selecting:            true,
		SelectionRow:         0,
		SelectedEnemyType:    EnemyTypeNormal,
		SelectedAttackSystem: DungeonAttackSword,
		SelectedCombatAction: DungeonActionSword,
		CombatTurn:           CombatTurnPlayer,
		TurnNumber:           1,
		PendingAction:        DungeonPendingNone,
		Message:              "CHOISIS TON ENTRAINEMENT",
		MessageTimer:         0,
	}
}

func (d *TrainingDungeon) Reset() {
	d.Active = false

	d.Selecting = true

	d.SelectionRow = 0
	d.SelectedEnemyType = EnemyTypeNormal
	d.SelectedAttackSystem = DungeonAttackSword

	d.SelectedCombatAction = DungeonActionSword

	d.Enemy = nil

	d.CombatTurn = CombatTurnPlayer
	d.TurnNumber = 1

	d.ActionInProgress = false
	d.ActionTimer = 0
	d.DamageApplied = false
	d.PendingAction = DungeonPendingNone

	d.Finished = false
	d.Won = false

	d.Message = "CHOISIS TON ENTRAINEMENT"
	d.MessageTimer = 0
}

func (g *Game) StartTrainingDungeon() {
	if g.Dungeon == nil {
		g.Dungeon = NewTrainingDungeon()
	}

	g.Dungeon.Reset()

	g.Dungeon.Active = true
	g.Dungeon.Selecting = true

	g.State = GameStatePlaying
	g.CurrentLevel = 0

	groundY = g.Levels[g.CurrentLevel].GroundY

	g.Player = NewPlayer()

	g.Player.X = 55
	g.Player.Y = groundY - g.Player.Height
	g.Player.VelocityY = 0
	g.Player.OnGround = true

	g.Boss = nil

	g.Enemies = []*Enemy{}

	g.Exit = nil
	g.Puzzle = nil

	g.Pickups = []*Pickup{}

	g.DamageTexts = []*DamageText{}
	g.Impacts = []*Impact{}

	if g.Combat != nil {
		g.Combat.Reset()
	}

	if g.Story != nil {
		g.Story.Stop()
	}

	if g.SkillTree != nil {
		g.SkillTree.Open = false
	}

	if g.Atmosphere != nil {
		g.Atmosphere.Reset(
			g.Levels[g.CurrentLevel].Number,
		)
	}

	if g.Music != nil {
		g.Music.PlayLevel()
	}
}

func (g *Game) StartSelectedTrainingCombat() {
	if g.Dungeon == nil {
		return
	}

	d := g.Dungeon

	d.Selecting = false

	d.Finished = false
	d.Won = false

	d.ActionInProgress = false
	d.ActionTimer = 0
	d.DamageApplied = false
	d.PendingAction = DungeonPendingNone

	d.TurnNumber = 1

	g.Player = NewPlayer()

	g.Player.X = 55
	g.Player.Y = groundY - g.Player.Height
	g.Player.VelocityY = 0
	g.Player.OnGround = true

	g.ApplyCurrentLevelProgression()

	switch d.SelectedEnemyType {
	case EnemyTypeFast:
		d.Enemy = NewFastEnemy(225)

	case EnemyTypeHeavy:
		d.Enemy = NewHeavyEnemy(225)

	case EnemyTypeRanged:
		d.Enemy = NewRangedEnemy(225)

	default:
		d.Enemy = NewEnemy(225)
	}

	d.Enemy.X = 225
	d.Enemy.Y = groundY - d.Enemy.Height

	d.Enemy.Moving = false
	d.Enemy.Attacking = false
	d.Enemy.CurrentFrame = 0
	d.Enemy.ProjectileActive = false

	g.Enemies = []*Enemy{
		d.Enemy,
	}

	if d.SelectedAttackSystem == DungeonAttackMagic {
		d.SelectedCombatAction = DungeonActionRay
	} else {
		d.SelectedCombatAction = DungeonActionSword
	}

	if d.Enemy.Initiative > g.Player.Initiative {
		d.CombatTurn = CombatTurnMonster
		d.Message = "LE MONSTRE EST PLUS RAPIDE"
	} else {
		d.CombatTurn = CombatTurnPlayer
		d.Message = "A TOI DE JOUER"
	}

	d.MessageTimer = 80
}

func (g *Game) LeaveTrainingDungeon() {
	if g.Dungeon != nil {
		g.Dungeon.Reset()
	}

	if g.Combat != nil {
		g.Combat.Reset()
	}

	g.Enemies = []*Enemy{}

	g.Boss = nil
	g.Exit = nil
	g.Puzzle = nil

	g.Pickups = []*Pickup{}

	g.DamageTexts = []*DamageText{}
	g.Impacts = []*Impact{}

	g.State = GameStateMenu

	g.Paused = false
	g.OptionsFromPause = false

	if g.Menu != nil {
		g.Menu.Selected = 2
	}

	if g.Music != nil {
		g.Music.Stop()
	}
}

func (g *Game) UpdateTrainingDungeon() error {
	if g.Dungeon == nil || !g.Dungeon.Active {
		return nil
	}

	g.UpdateAtmosphere()
	g.UpdateDamageTexts()
	g.UpdateImpacts()

	if g.ShakeTimer > 0 {
		g.ShakeTimer--
	}

	if g.Dungeon.MessageTimer > 0 {
		g.Dungeon.MessageTimer--
	}

	if g.Dungeon.Selecting {
		g.Dungeon.UpdateSelection(g)
		return nil
	}

	if g.Player == nil || g.Dungeon.Enemy == nil {
		g.LeaveTrainingDungeon()
		return nil
	}

	g.Dungeon.updateVisualTimers(
		g.Player,
	)

	if g.Dungeon.Enemy.FlashTimer > 0 {
		g.Dungeon.Enemy.FlashTimer--
	}

	if g.Dungeon.Enemy.Dying {
		g.Dungeon.Enemy.Update(
			g.Player,
		)
	}

	if g.Dungeon.Finished {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.Dungeon.Selecting = true

			g.Dungeon.Finished = false
			g.Dungeon.Won = false

			g.Dungeon.Enemy = nil

			g.Enemies = []*Enemy{}

			g.Player = NewPlayer()

			g.Player.X = 55
			g.Player.Y = groundY - g.Player.Height

			g.Dungeon.Message = "CHOISIS TON ENTRAINEMENT"

			return nil
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.LeaveTrainingDungeon()
			return nil
		}

		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Dungeon.Selecting = true

		g.Dungeon.Enemy = nil

		g.Enemies = []*Enemy{}

		g.Dungeon.ActionInProgress = false
		g.Dungeon.PendingAction = DungeonPendingNone

		g.Dungeon.Message = "CHOISIS TON ENTRAINEMENT"

		return nil
	}

	if !g.Player.Alive {
		g.Dungeon.Finished = true
		g.Dungeon.Won = false

		g.Dungeon.Message = "ENTRAINEMENT PERDU"

		return nil
	}

	if g.Dungeon.Enemy.Dying || !g.Dungeon.Enemy.Alive {
		g.Dungeon.Finished = true
		g.Dungeon.Won = true

		g.Dungeon.Message = "ENTRAINEMENT REUSSI"

		return nil
	}

	if g.Dungeon.ActionInProgress {
		g.Dungeon.UpdateAction(g)
		return nil
	}

	switch g.Dungeon.CombatTurn {
	case CombatTurnPlayer:
		g.Dungeon.UpdatePlayerTurn(g)

	case CombatTurnMonster:
		g.Dungeon.StartMonsterAction(g)
	}

	return nil
}

func (d *TrainingDungeon) UpdateSelection(game *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.LeaveTrainingDungeon()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		d.SelectionRow--

		if d.SelectionRow < 0 {
			d.SelectionRow = 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		d.SelectionRow++

		if d.SelectionRow > 1 {
			d.SelectionRow = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if d.SelectionRow == 0 {
			d.SelectedEnemyType--

			if d.SelectedEnemyType < EnemyTypeNormal {
				d.SelectedEnemyType = EnemyTypeRanged
			}
		} else {
			d.SelectedAttackSystem--

			if d.SelectedAttackSystem < DungeonAttackSword {
				d.SelectedAttackSystem = DungeonAttackFree
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if d.SelectionRow == 0 {
			d.SelectedEnemyType++

			if d.SelectedEnemyType > EnemyTypeRanged {
				d.SelectedEnemyType = EnemyTypeNormal
			}
		} else {
			d.SelectedAttackSystem++

			if d.SelectedAttackSystem > DungeonAttackFree {
				d.SelectedAttackSystem = DungeonAttackSword
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		game.StartSelectedTrainingCombat()
	}
}

func (d *TrainingDungeon) updateVisualTimers(player *Player) {
	if player.InvincibleTimer > 0 {
		player.InvincibleTimer--
	}

	if player.FlashTimer > 0 {
		player.FlashTimer--
	}

	if player.HurtTimer > 0 {
		player.HurtTimer--
	}

	player.UpdateDodgeDust()

	if player.Attacking {
		player.UpdateAttackAnimation()
	}
}

func (d *TrainingDungeon) UpdatePlayerTurn(game *Game) {
	if d.SelectedAttackSystem == DungeonAttackSword {
		d.SelectedCombatAction = DungeonActionSword
	}

	if d.SelectedAttackSystem == DungeonAttackMagic {
		d.SelectedCombatAction = DungeonActionRay
	}

	if d.SelectedAttackSystem == DungeonAttackFree {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			d.SelectedCombatAction--

			if d.SelectedCombatAction < DungeonActionSword {
				d.SelectedCombatAction = DungeonActionRay
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			d.SelectedCombatAction++

			if d.SelectedCombatAction > DungeonActionRay {
				d.SelectedCombatAction = DungeonActionSword
			}
		}
	}

	if d.MessageTimer <= 0 {
		switch d.SelectedCombatAction {
		case DungeonActionRay:
			d.Message = "RAYON D'ELDORIA"

		default:
			d.Message = "ATTAQUE A L'EPEE"
		}
	}

	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return
	}

	switch d.SelectedCombatAction {
	case DungeonActionRay:
		d.StartRayAction(game)

	default:
		d.StartSwordAction(game)
	}
}

func (d *TrainingDungeon) StartSwordAction(game *Game) {
	d.ActionInProgress = true
	d.ActionTimer = 0

	d.DamageApplied = false

	d.PendingAction = DungeonPendingSword

	d.Message = "TU ATTAQUES A L'EPEE"

	game.Player.FacingLeft = false

	game.Player.StartAttack()

	if game.Sounds != nil {
		game.Sounds.PlaySword()
	}
}

func (d *TrainingDungeon) StartRayAction(game *Game) {
	if game.Player.Mana < TrainingRayonManaCost {
		d.Message = fmt.Sprintf(
			"MANA INSUFFISANT : %d/%d",
			game.Player.Mana,
			TrainingRayonManaCost,
		)

		d.MessageTimer = 90

		return
	}

	if !game.Player.SpendMana(TrainingRayonManaCost) {
		d.Message = "MANA INSUFFISANT"

		d.MessageTimer = 90

		return
	}

	d.ActionInProgress = true
	d.ActionTimer = 0

	d.DamageApplied = false

	d.PendingAction = DungeonPendingRay

	d.Message = fmt.Sprintf(
		"RAYON D'ELDORIA  -%d MANA",
		TrainingRayonManaCost,
	)

	game.Player.FacingLeft = false

	game.Player.StartEnergyAnimation()
}

func (d *TrainingDungeon) StartMonsterAction(game *Game) {
	d.ActionInProgress = true
	d.ActionTimer = 0

	d.DamageApplied = false

	d.PendingAction = DungeonPendingMonster

	d.Message = "LE MONSTRE ATTAQUE"

	d.Enemy.Attacking = true
	d.Enemy.Moving = false
	d.Enemy.CurrentFrame = 3
	d.Enemy.FacingLeft = true
}

func (d *TrainingDungeon) UpdateAction(game *Game) {
	d.ActionTimer++

	switch d.PendingAction {
	case DungeonPendingSword:
		d.UpdateSwordAction(game)

	case DungeonPendingRay:
		d.UpdateRayAction(game)

	case DungeonPendingMonster:
		d.UpdateMonsterAction(game)
	}
}

func (d *TrainingDungeon) UpdateSwordAction(game *Game) {
	if d.ActionTimer == 12 && !d.DamageApplied {
		d.DamageApplied = true

		damage := game.Player.AttackDamage

		if game.SkillTree != nil {
			damage = game.SkillTree.NormalAttackDamage(
				damage,
			)
		}

		d.ApplyDamageToEnemy(
			game,
			damage,
		)

		if d.Finished {
			return
		}
	}

	if d.ActionTimer < 30 {
		return
	}

	game.Player.Attacking = false
	game.Player.AttackTimer = 0
	game.Player.CurrentAttackFrame = 0
	game.Player.AttackHit = false

	d.EndPlayerAction()
}

func (d *TrainingDungeon) UpdateRayAction(game *Game) {
	player := game.Player

	switch {
	case d.ActionTimer <= 7:
		player.SetEnergyFrame(0)
		d.Message = "CONCENTRATION..."

	case d.ActionTimer <= 14:
		player.SetEnergyFrame(1)
		d.Message = "L'ENERGIE APPARAIT..."

	case d.ActionTimer <= 21:
		player.SetEnergyFrame(2)
		d.Message = "CHARGE..."

	case d.ActionTimer <= 28:
		player.SetEnergyFrame(3)
		d.Message = "PUISSANCE MAXIMALE..."

	case d.ActionTimer <= 35:
		player.SetEnergyFrame(4)
		d.Message = "RAYON D'ELDORIA !"

	case d.ActionTimer <= 45:
		player.SetEnergyFrame(5)
		d.Message = "RAYON D'ELDORIA !"

	case d.ActionTimer <= 55:
		player.SetEnergyFrame(6)
		d.Message = "PLEINE PUISSANCE !"

	case d.ActionTimer <= 65:
		player.SetEnergyFrame(7)
		d.Message = "IMPACT !"
	}

	if d.ActionTimer == 48 && !d.DamageApplied {
		d.DamageApplied = true

		damage := game.Player.AttackDamage * 2

		if game.SkillTree != nil {
			damage = game.SkillTree.RayDamage(
				game.Player.AttackDamage,
			)
		}

		d.ApplyDamageToEnemy(
			game,
			damage,
		)

		game.StartShake(
			18,
			3,
		)

		if d.Finished {
			player.StopEnergyAnimation()
			return
		}
	}

	if d.ActionTimer < 72 {
		return
	}

	player.StopEnergyAnimation()

	d.EndPlayerAction()
}

func (d *TrainingDungeon) ApplyDamageToEnemy(game *Game, damage int) {
	if d.Enemy == nil {
		return
	}

	enemyBox := d.Enemy.HitBox()

	if d.Enemy.Hit(
		damage,
		game.Player.X+game.Player.Width/2,
	) {
		if game.Sounds != nil {
			game.Sounds.PlayImpact()
		}

		game.Impacts = append(
			game.Impacts,
			NewImpact(
				enemyBox.X+enemyBox.W/2,
				enemyBox.Y+enemyBox.H/2,
				false,
			),
		)

		game.DamageTexts = append(
			game.DamageTexts,
			NewDamageText(
				d.Enemy.X+d.Enemy.Width/2,
				d.Enemy.Y-10,
				damage,
			),
		)

		game.StartShake(
			8,
			1.5,
		)
	}

	if d.Enemy.Dying || !d.Enemy.Alive {
		d.Finished = true
		d.Won = true

		d.ActionInProgress = false

		d.Message = "ENTRAINEMENT REUSSI"

		game.Player.Attacking = false
		game.Player.StopEnergyAnimation()
	}
}

func (d *TrainingDungeon) EndPlayerAction() {
	d.ActionInProgress = false
	d.ActionTimer = 0

	d.DamageApplied = false

	d.PendingAction = DungeonPendingNone

	d.CombatTurn = CombatTurnMonster

	d.TurnNumber++

	d.Message = "TOUR DU MONSTRE"
}

func (d *TrainingDungeon) UpdateMonsterAction(game *Game) {
	d.Enemy.CurrentFrame = 3

	if d.ActionTimer == 18 && !d.DamageApplied {
		d.DamageApplied = true

		damage := d.Enemy.Damage

		sourceX := d.Enemy.X + d.Enemy.Width/2

		game.Player.InvincibleTimer = 0

		if game.Player.Hit(
			damage,
			sourceX,
		) {
			if game.Sounds != nil {
				game.Sounds.PlayHurt()
			}

			playerBox := game.Player.HitBox()

			game.Impacts = append(
				game.Impacts,
				NewImpact(
					playerBox.X+playerBox.W/2,
					playerBox.Y+playerBox.H/2,
					true,
				),
			)

			game.DamageTexts = append(
				game.DamageTexts,
				NewDamageText(
					game.Player.X,
					game.Player.Y-10,
					damage,
				),
			)

			game.StartShake(
				7,
				1.4,
			)
		}

		game.Player.KnockbackX = 0
		game.Player.VelocityY = 0

		game.Player.Y = groundY - game.Player.Height
		game.Player.OnGround = true

		if !game.Player.Alive {
			d.Finished = true
			d.Won = false

			d.ActionInProgress = false

			d.Message = "ENTRAINEMENT PERDU"

			d.Enemy.Attacking = false
			d.Enemy.CurrentFrame = 0

			return
		}
	}

	if d.ActionTimer < 38 {
		return
	}

	d.Enemy.Attacking = false
	d.Enemy.AttackTimer = 0
	d.Enemy.CurrentFrame = 0

	d.ActionInProgress = false
	d.ActionTimer = 0

	d.DamageApplied = false

	d.PendingAction = DungeonPendingNone

	d.CombatTurn = CombatTurnPlayer

	d.TurnNumber++

	d.Message = "A TOI DE JOUER"
}

func (d *TrainingDungeon) EnemyName() string {
	switch d.SelectedEnemyType {
	case EnemyTypeFast:
		return "OMBRE RAPIDE"

	case EnemyTypeHeavy:
		return "GARDE LOURD"

	case EnemyTypeRanged:
		return "MAGE DE L'OMBRE"
	}

	return "OMBRE NORMALE"
}

func (d *TrainingDungeon) AttackSystemName() string {
	switch d.SelectedAttackSystem {
	case DungeonAttackMagic:
		return "MAGIE"

	case DungeonAttackFree:
		return "LIBRE"
	}

	return "EPEE"
}

func (d *TrainingDungeon) CombatActionName() string {
	if d.SelectedCombatAction == DungeonActionRay {
		return "RAYON D'ELDORIA"
	}

	return "EPEE"
}

func (d *TrainingDungeon) DrawHUD(screen *ebiten.Image, player *Player) {
	if d == nil || !d.Active {
		return
	}

	if d.Selecting {
		d.DrawSelection(screen)
		return
	}

	if player == nil || d.Enemy == nil {
		return
	}

	d.DrawCombatHUD(
		screen,
		player,
	)
}

func (d *TrainingDungeon) DrawSelection(screen *ebiten.Image) {
	hudRect(
		screen,
		230,
		100,
		820,
		500,
		color.RGBA{
			R: 4,
			G: 8,
			B: 15,
			A: 235,
		},
	)

	hudRect(
		screen,
		230,
		100,
		820,
		4,
		color.RGBA{
			R: 80,
			G: 180,
			B: 255,
			A: 255,
		},
	)

	hudCenteredText(
		screen,
		"DONJON D'ENTRAINEMENT",
		float64(renderWidth)/2,
		135,
	)

	hudCenteredText(
		screen,
		"PREPARE TON COMBAT",
		float64(renderWidth)/2,
		175,
	)

	d.drawSelectionLine(
		screen,
		0,
		"ENNEMI",
		d.EnemyName(),
		235,
	)

	d.drawSelectionLine(
		screen,
		1,
		"SYSTEME D'ATTAQUE",
		d.AttackSystemName(),
		325,
	)

	hudCenteredText(
		screen,
		"ENNEMIS DISPONIBLES",
		float64(renderWidth)/2,
		410,
	)

	hudCenteredText(
		screen,
		"NORMALE  /  RAPIDE  /  LOURD  /  MAGE",
		float64(renderWidth)/2,
		440,
	)

	hudCenteredText(
		screen,
		"EPEE : ATTAQUES PHYSIQUES",
		float64(renderWidth)/2,
		485,
	)

	hudCenteredText(
		screen,
		"MAGIE : RAYON D'ELDORIA - 40 MANA",
		float64(renderWidth)/2,
		510,
	)

	hudCenteredText(
		screen,
		"LIBRE : CHOIX EPEE OU MAGIE A CHAQUE TOUR",
		float64(renderWidth)/2,
		535,
	)

	hudCenteredText(
		screen,
		"HAUT/BAS : LIGNE   GAUCHE/DROITE : CHOIX",
		float64(renderWidth)/2,
		565,
	)

	hudCenteredText(
		screen,
		"ENTREE : COMMENCER   |   ECHAP : MENU",
		float64(renderWidth)/2,
		590,
	)
}

func (d *TrainingDungeon) drawSelectionLine(screen *ebiten.Image, row int, title string, value string, y float64) {
	selected := d.SelectionRow == row

	backgroundColor := color.RGBA{
		R: 15,
		G: 20,
		B: 30,
		A: 230,
	}

	borderColor := color.RGBA{
		R: 65,
		G: 80,
		B: 100,
		A: 255,
	}

	if selected {
		backgroundColor = color.RGBA{
			R: 20,
			G: 55,
			B: 95,
			A: 235,
		}

		borderColor = color.RGBA{
			R: 80,
			G: 180,
			B: 255,
			A: 255,
		}
	}

	hudRect(
		screen,
		340,
		y,
		600,
		60,
		borderColor,
	)

	hudRect(
		screen,
		343,
		y+3,
		594,
		54,
		backgroundColor,
	)

	prefix := "  "

	if selected {
		prefix = "> "
	}

	drawHUDText(
		screen,
		prefix+title,
		380,
		y+17,
	)

	drawHUDText(
		screen,
		"< "+value+" >",
		700,
		y+17,
	)
}

func (d *TrainingDungeon) DrawCombatHUD(screen *ebiten.Image, player *Player) {
	hudRect(
		screen,
		360,
		18,
		560,
		82,
		color.RGBA{
			R: 4,
			G: 8,
			B: 15,
			A: 215,
		},
	)

	hudCenteredText(
		screen,
		"DONJON D'ENTRAINEMENT",
		float64(renderWidth)/2,
		31,
	)

	turnText := "TOUR DU JOUEUR"

	if d.CombatTurn == CombatTurnMonster {
		turnText = "TOUR DU MONSTRE"
	}

	hudCenteredText(
		screen,
		fmt.Sprintf(
			"TOUR %d  -  %s",
			d.TurnNumber,
			turnText,
		),
		float64(renderWidth)/2,
		57,
	)

	hudCenteredText(
		screen,
		"SYSTEME : "+d.AttackSystemName(),
		float64(renderWidth)/2,
		78,
	)

	playerRatio := 0.0

	if player.MaxHP > 0 {
		playerRatio = float64(player.HP) / float64(player.MaxHP)
	}

	enemyRatio := 0.0

	if d.Enemy.MaxHP > 0 {
		enemyRatio = float64(d.Enemy.HP) / float64(d.Enemy.MaxHP)
	}

	drawHUDText(
		screen,
		"JOUEUR",
		65,
		110,
	)

	drawMinimalBar(
		screen,
		65,
		138,
		300,
		11,
		playerRatio,
	)

	drawHUDText(
		screen,
		fmt.Sprintf(
			"%d/%d PV",
			player.HP,
			player.MaxHP,
		),
		65,
		158,
	)

	drawHUDText(
		screen,
		fmt.Sprintf(
			"MANA %d/%d",
			player.Mana,
			player.ManaMax,
		),
		65,
		184,
	)

	drawHUDText(
		screen,
		d.EnemyName(),
		865,
		110,
	)

	drawMinimalBar(
		screen,
		865,
		138,
		300,
		11,
		enemyRatio,
	)

	drawHUDText(
		screen,
		fmt.Sprintf(
			"%d/%d PV",
			d.Enemy.HP,
			d.Enemy.MaxHP,
		),
		865,
		158,
	)

	drawHUDText(
		screen,
		fmt.Sprintf(
			"INITIATIVE %d",
			d.Enemy.Initiative,
		),
		865,
		184,
	)

	hudRect(
		screen,
		280,
		560,
		720,
		120,
		color.RGBA{
			R: 4,
			G: 8,
			B: 15,
			A: 225,
		},
	)

	hudCenteredText(
		screen,
		d.Message,
		float64(renderWidth)/2,
		580,
	)

	if d.Finished {
		result := "VICTOIRE"

		if !d.Won {
			result = "DEFAITE"
		}

		hudCenteredText(
			screen,
			result,
			float64(renderWidth)/2,
			610,
		)

		hudCenteredText(
			screen,
			"ENTREE : NOUVEL ENTRAINEMENT",
			float64(renderWidth)/2,
			640,
		)

		hudCenteredText(
			screen,
			"ECHAP : MENU",
			float64(renderWidth)/2,
			663,
		)

		return
	}

	if d.CombatTurn == CombatTurnPlayer && !d.ActionInProgress {
		if d.SelectedAttackSystem == DungeonAttackFree {
			hudCenteredText(
				screen,
				"<  "+d.CombatActionName()+"  >",
				float64(renderWidth)/2,
				610,
			)

			hudCenteredText(
				screen,
				"FLECHES : CHOISIR   |   ENTREE : ATTAQUER",
				float64(renderWidth)/2,
				640,
			)
		} else {
			hudCenteredText(
				screen,
				d.CombatActionName(),
				float64(renderWidth)/2,
				610,
			)

			hudCenteredText(
				screen,
				"ENTREE : ATTAQUER",
				float64(renderWidth)/2,
				640,
			)
		}

		hudCenteredText(
			screen,
			"ECHAP : RETOUR AUX CHOIX",
			float64(renderWidth)/2,
			663,
		)

		return
	}

	hudCenteredText(
		screen,
		"PATIENTE...",
		float64(renderWidth)/2,
		630,
	)
}
