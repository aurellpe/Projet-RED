package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TurnStatePlayer = iota
	TurnStatePlayerAction
	TurnStateEnemy
)

const (
	TurnActionAttack = iota
	TurnActionSkill
	TurnActionDefend
	TurnActionItem
)

type TurnCombat struct {
	Active         bool
	State          int
	Selected       int
	Timer          int
	Player         *Player
	Enemy          *Enemy
	Boss           *Boss
	Skills         *SkillTree
	Message        string
	PendingAction  int
	PendingDamage  int
	DamageApplied  bool
	Defending      bool
	Potions        int
	SkillCooldown  int
	TargetDefeated bool
}

func NewTurnCombat() *TurnCombat {
	return &TurnCombat{
		State:   TurnStatePlayer,
		Potions: 2,
	}
}

func (c *TurnCombat) Reset() {
	if c.Player != nil {
		c.Player.StopEnergyAnimation()
	}

	c.Active = false
	c.State = TurnStatePlayer
	c.Selected = 0
	c.Timer = 0
	c.Player = nil
	c.Enemy = nil
	c.Boss = nil
	c.Skills = nil
	c.Message = ""
	c.PendingAction = TurnActionAttack
	c.PendingDamage = 0
	c.DamageApplied = false
	c.Defending = false
	c.Potions = 2
	c.SkillCooldown = 0
	c.TargetDefeated = false
}

func (c *TurnCombat) StartEnemy(enemy *Enemy, player *Player, skills *SkillTree) {
	if enemy == nil || player == nil {
		return
	}

	c.Reset()

	c.Active = true
	c.Player = player
	c.Enemy = enemy
	c.Skills = skills
	c.Potions = 2
	c.Message = "TON TOUR"

	enemy.Attacking = false
	enemy.Moving = false
	enemy.CurrentFrame = 0

	player.Attacking = false
	player.Moving = false
	player.CurrentWalkFrame = 0
	player.StopEnergyAnimation()

	c.FaceCharacters(player)
}

func (c *TurnCombat) StartBoss(boss *Boss, player *Player, skills *SkillTree) {
	if boss == nil || player == nil {
		return
	}

	c.Reset()

	c.Active = true
	c.Player = player
	c.Boss = boss
	c.Skills = skills
	c.Potions = 3
	c.Message = "TON TOUR"

	boss.Attacking = false
	boss.Moving = false
	boss.CurrentAttackFrame = 0
	boss.CurrentWalkFrame = 0

	player.Attacking = false
	player.Moving = false
	player.CurrentWalkFrame = 0
	player.StopEnergyAnimation()

	c.FaceCharacters(player)
}

func (c *TurnCombat) FaceCharacters(player *Player) {
	if player == nil {
		return
	}

	playerCenter := player.X + player.Width/2
	targetCenter := c.TargetCenterX()

	player.FacingLeft = playerCenter >= targetCenter

	if c.Enemy != nil {
		c.Enemy.FacingLeft = targetCenter > playerCenter
	}

	if c.Boss != nil {
		c.Boss.FacingLeft = targetCenter > playerCenter
	}
}

func (c *TurnCombat) TargetCenterX() float64 {
	if c.Enemy != nil {
		return c.Enemy.X + c.Enemy.Width/2
	}

	if c.Boss != nil {
		return c.Boss.X + c.Boss.HitWidth/2
	}

	return 0
}

func (c *TurnCombat) TargetAlive() bool {
	if c.Enemy != nil {
		return c.Enemy.Alive && !c.Enemy.Dying
	}

	if c.Boss != nil {
		return c.Boss.Alive && !c.Boss.Dying
	}

	return false
}

func (c *TurnCombat) Update(game *Game) {
	if !c.Active || game == nil || game.Player == nil {
		return
	}

	c.UpdateVisualTimers(game.Player)

	if c.Boss != nil {
		c.Boss.AuraTimer++
		c.Boss.CombatTimer++

		if c.Boss.FlashTimer > 0 {
			c.Boss.FlashTimer--
		}

		if c.Boss.PhaseTransition {
			game.Player.StopEnergyAnimation()

			if c.State == TurnStatePlayerAction && c.PendingAction == TurnActionSkill {
				c.Timer = 71
			}

			c.Message = "LE GARDIEN LIBERE SA PUISSANCE..."

			c.Boss.UpdatePhaseTransition()

			return
		}
	}

	if c.Enemy != nil && c.Enemy.FlashTimer > 0 {
		c.Enemy.FlashTimer--
	}

	if !game.Player.Alive {
		c.EndCombat(game.Player)
		return
	}

	switch c.State {
	case TurnStatePlayer:
		c.UpdatePlayerTurn(game)

	case TurnStatePlayerAction:
		c.UpdatePlayerAction(game)

	case TurnStateEnemy:
		c.UpdateEnemyTurn(game)
	}
}

func (c *TurnCombat) UpdateVisualTimers(player *Player) {
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

func (c *TurnCombat) UpdatePlayerTurn(game *Game) {
	if !c.TargetAlive() {
		c.EndCombat(game.Player)
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		c.Selected--

		if c.Selected < 0 {
			c.Selected = 3
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		c.Selected++

		if c.Selected > 3 {
			c.Selected = 0
		}
	}

	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}

	switch c.Selected {
	case TurnActionAttack:
		c.StartPlayerAttack(game)

	case TurnActionSkill:
		c.StartEnergyAttack(game)

	case TurnActionDefend:
		c.Defending = true
		c.PendingAction = TurnActionDefend
		c.State = TurnStatePlayerAction
		c.Timer = 0
		c.Message = "TU TE METS EN GARDE"

	case TurnActionItem:
		c.UsePotion(game)
	}
}

func (c *TurnCombat) StartPlayerAttack(game *Game) {
	player := game.Player

	c.FaceCharacters(player)

	player.StopEnergyAnimation()
	player.StartAttack()

	damage := player.AttackDamage

	if c.Skills != nil {
		damage = c.Skills.NormalAttackDamage(damage)
	}

	c.PendingAction = TurnActionAttack
	c.PendingDamage = damage
	c.State = TurnStatePlayerAction
	c.Timer = 0
	c.DamageApplied = false
	c.TargetDefeated = false
	c.Message = "ATTAQUE"
}

func (c *TurnCombat) StartEnergyAttack(game *Game) {
	if c.SkillCooldown > 0 {
		c.Message = fmt.Sprintf(
			"RAYON DISPONIBLE DANS %d TOUR(S)",
			c.SkillCooldown,
		)

		return
	}

	player := game.Player

	c.FaceCharacters(player)

	player.StartEnergyAnimation()

	damage := player.AttackDamage * 2
	cooldown := 3

	if c.Skills != nil {
		damage = c.Skills.RayDamage(player.AttackDamage)
		cooldown = c.Skills.RayCooldown()
	}

	c.PendingAction = TurnActionSkill
	c.PendingDamage = damage
	c.SkillCooldown = cooldown
	c.State = TurnStatePlayerAction
	c.Timer = 0
	c.DamageApplied = false
	c.TargetDefeated = false
	c.Message = "RAYON D'ELDORIA"
}

func (c *TurnCombat) UsePotion(game *Game) {
	player := game.Player

	if c.Potions <= 0 {
		c.Message = "PLUS DE POTION"
		return
	}

	if player.HP >= player.MaxHP {
		c.Message = "TES PV SONT DEJA AU MAXIMUM"
		return
	}

	heal := 30

	if c.Skills != nil {
		heal = c.Skills.PotionHeal()
	}

	player.HP += heal

	if player.HP > player.MaxHP {
		player.HP = player.MaxHP
	}

	c.Potions--
	c.PendingAction = TurnActionItem
	c.State = TurnStatePlayerAction
	c.Timer = 0

	c.Message = fmt.Sprintf("+%d PV", heal)
}

func (c *TurnCombat) UpdatePlayerAction(game *Game) {
	c.Timer++

	switch c.PendingAction {
	case TurnActionAttack:
		if c.Timer == 12 && !c.DamageApplied {
			c.ApplyPlayerDamage(game)
		}

		if c.Timer >= 30 {
			game.Player.Attacking = false
			game.Player.CurrentAttackFrame = 0
			game.Player.AttackTimer = 0

			if c.TargetDefeated {
				c.EndCombat(game.Player)
				return
			}

			c.StartEnemyTurn(game)
		}

	case TurnActionSkill:
		c.UpdateEnergyAnimation(game)

	case TurnActionDefend, TurnActionItem:
		if c.Timer >= 18 {
			c.StartEnemyTurn(game)
		}
	}
}

func (c *TurnCombat) UpdateEnergyAnimation(game *Game) {
	player := game.Player

	switch {
	case c.Timer <= 7:
		player.SetEnergyFrame(0)
		c.Message = "CONCENTRATION..."

	case c.Timer <= 14:
		player.SetEnergyFrame(1)
		c.Message = "L'ENERGIE APPARAIT..."

	case c.Timer <= 21:
		player.SetEnergyFrame(2)
		c.Message = "CHARGE..."

	case c.Timer <= 28:
		player.SetEnergyFrame(3)
		c.Message = "PUISSANCE MAXIMALE..."

	case c.Timer <= 35:
		player.SetEnergyFrame(4)
		c.Message = "RAYON D'ELDORIA !"

	case c.Timer <= 45:
		player.SetEnergyFrame(5)
		c.Message = "RAYON D'ELDORIA !"

	case c.Timer <= 55:
		player.SetEnergyFrame(6)
		c.Message = "PLEINE PUISSANCE !"

	case c.Timer <= 65:
		player.SetEnergyFrame(7)
		c.Message = "IMPACT !"
	}

	if c.Timer == 48 && !c.DamageApplied {
		c.ApplyPlayerDamage(game)
		game.StartShake(20, 3.2)
	}

	if c.Timer < 72 {
		return
	}

	player.StopEnergyAnimation()

	if c.TargetDefeated {
		c.EndCombat(player)
		return
	}

	c.StartEnemyTurn(game)
}

func (c *TurnCombat) ApplyPlayerDamage(game *Game) {
	c.DamageApplied = true

	player := game.Player
	damage := c.PendingDamage

	player.AttackHit = true

	if c.Enemy != nil {
		box := c.Enemy.HitBox()

		if c.Enemy.Hit(
			damage,
			player.X+player.Width/2,
		) {
			game.Sounds.PlayImpact()

			game.Impacts = append(
				game.Impacts,
				NewImpact(
					box.X+box.W/2,
					box.Y+box.H/2,
					false,
				),
			)

			game.DamageTexts = append(
				game.DamageTexts,
				NewDamageText(
					c.Enemy.X+c.Enemy.Width/2,
					c.Enemy.Y-10,
					damage,
				),
			)

			game.StartShake(10, 2)
		}

		if c.Enemy.Dying || !c.Enemy.Alive {
			c.TargetDefeated = true
			c.Message = "ENNEMI VAINCU"
		}

		return
	}

	if c.Boss != nil {
		box := c.Boss.HitBox()

		if c.Boss.Hit(
			damage,
			player.X+player.Width/2,
		) {
			game.Sounds.PlayImpact()

			game.Impacts = append(
				game.Impacts,
				NewImpact(
					box.X+box.W/2,
					box.Y+box.H/2,
					false,
				),
			)

			game.DamageTexts = append(
				game.DamageTexts,
				NewDamageText(
					c.Boss.X+c.Boss.HitWidth/2,
					groundY-c.Boss.HitHeight-10,
					damage,
				),
			)

			game.StartShake(14, 2.8)
		}

		if !c.Boss.Dying && c.Boss.Phase == 1 && c.Boss.HP <= c.Boss.MaxHP/2 {
			c.Boss.StartPhaseTwo()
		}

		if c.Boss.Dying || !c.Boss.Alive {
			c.TargetDefeated = true
			c.Message = "LE GARDIEN EST VAINCU"
		}
	}
}

func (c *TurnCombat) StartEnemyTurn(game *Game) {
	game.Player.StopEnergyAnimation()

	c.State = TurnStateEnemy
	c.Timer = 0
	c.DamageApplied = false
	c.Message = "TOUR DE L'ENNEMI"

	c.FaceCharacters(game.Player)

	if c.Enemy != nil {
		c.Enemy.Attacking = true
		c.Enemy.Moving = false
		c.Enemy.CurrentFrame = 3
	}

	if c.Boss != nil {
		c.Boss.Attacking = true
		c.Boss.Moving = false
		c.Boss.CurrentAttackFrame = 1
	}
}

func (c *TurnCombat) UpdateEnemyTurn(game *Game) {
	c.Timer++

	if c.Enemy != nil {
		c.Enemy.CurrentFrame = 3
	}

	if c.Boss != nil {
		if c.Timer <= 9 {
			c.Boss.CurrentAttackFrame = 1
		} else if c.Timer <= 18 {
			c.Boss.CurrentAttackFrame = 2
		} else {
			c.Boss.CurrentAttackFrame = 3
		}
	}

	if c.Timer == 14 && !c.DamageApplied {
		c.ApplyEnemyDamage(game)

		if !game.Player.Alive {
			c.EndCombat(game.Player)
			return
		}
	}

	if c.Timer >= 32 {
		c.StopEnemyAnimation()

		c.Defending = false

		c.StartPlayerTurn()
	}
}

func (c *TurnCombat) ApplyEnemyDamage(game *Game) {
	c.DamageApplied = true

	player := game.Player

	damage := 10
	sourceX := c.TargetCenterX()

	if c.Enemy != nil {
		damage = c.Enemy.Damage
	}

	if c.Boss != nil {
		damage = c.Boss.GetAttackDamage()
	}

	if c.Defending {
		if c.Skills != nil {
			damage = c.Skills.GuardDamage(damage)
		} else {
			damage = (damage + 1) / 2
		}

		c.Message = "GARDE REUSSIE : DEGATS REDUITS"
	} else {
		c.Message = fmt.Sprintf(
			"TU SUBIS %d DEGATS",
			damage,
		)
	}

	player.InvincibleTimer = 0

	if player.Hit(
		damage,
		sourceX,
	) {
		game.Sounds.PlayHurt()

		box := player.HitBox()

		game.Impacts = append(
			game.Impacts,
			NewImpact(
				box.X+box.W/2,
				box.Y+box.H/2,
				true,
			),
		)

		game.DamageTexts = append(
			game.DamageTexts,
			NewDamageText(
				player.X,
				player.Y-10,
				damage,
			),
		)

		game.StartShake(8, 2)
	}

	player.KnockbackX = 0
	player.VelocityY = 0
	player.Y = groundY - player.Height
	player.OnGround = true
}

func (c *TurnCombat) StartPlayerTurn() {
	c.State = TurnStatePlayer
	c.Timer = 0
	c.DamageApplied = false
	c.PendingDamage = 0
	c.PendingAction = TurnActionAttack

	if c.SkillCooldown > 0 {
		c.SkillCooldown--
	}

	c.Message = "TON TOUR"
}

func (c *TurnCombat) StopEnemyAnimation() {
	if c.Enemy != nil {
		c.Enemy.Attacking = false
		c.Enemy.CurrentFrame = 0
		c.Enemy.AttackTimer = 0
	}

	if c.Boss != nil {
		c.Boss.Attacking = false
		c.Boss.CurrentAttackFrame = 0
		c.Boss.AttackTimer = 0
	}
}

func (c *TurnCombat) EndCombat(player *Player) {
	c.StopEnemyAnimation()

	if player != nil {
		player.StopEnergyAnimation()

		player.Attacking = false
		player.AttackTimer = 0
		player.CurrentAttackFrame = 0
		player.AttackHit = false

		player.HurtTimer = 0
		player.InvincibleTimer = 0
		player.KnockbackX = 0
		player.VelocityY = 0
		player.Y = groundY - player.Height
		player.OnGround = true
	}

	c.Active = false
	c.State = TurnStatePlayer
	c.Timer = 0
	c.Selected = 0
	c.DamageApplied = false
	c.Defending = false
}

func (c *TurnCombat) TargetName() string {
	if c.Boss != nil {
		return "GARDIEN D'ELDORIA"
	}

	if c.Enemy == nil {
		return "ENNEMI"
	}

	switch c.Enemy.Type {
	case EnemyTypeFast:
		return "OMBRE RAPIDE"

	case EnemyTypeHeavy:
		return "GARDE LOURD"

	case EnemyTypeRanged:
		return "MAGE DE L'OMBRE"
	}

	return "OMBRE"
}

func (c *TurnCombat) TargetHP() (int, int) {
	if c.Enemy != nil {
		return c.Enemy.HP, c.Enemy.MaxHP
	}

	if c.Boss != nil {
		return c.Boss.HP, c.Boss.MaxHP
	}

	return 0, 1
}

func (c *TurnCombat) Draw(screen *ebiten.Image) {
	if !c.Active {
		return
	}

	c.DrawEnergyScreenEffect(screen)

	panelY := 500.0

	hudRect(
		screen,
		40,
		panelY,
		1200,
		195,
		color.RGBA{
			R: 5,
			G: 7,
			B: 12,
			A: 230,
		},
	)

	hudRect(
		screen,
		40,
		panelY,
		1200,
		3,
		color.RGBA{
			R: 180,
			G: 140,
			B: 65,
			A: 255,
		},
	)

	targetName := c.TargetName()

	hp, maxHP := c.TargetHP()

	ratio := 0.0

	if maxHP > 0 {
		ratio = float64(hp) / float64(maxHP)
	}

	drawHUDText(
		screen,
		targetName,
		75,
		520,
	)

	drawMinimalBar(
		screen,
		75,
		548,
		330,
		10,
		ratio,
	)

	drawHUDText(
		screen,
		fmt.Sprintf("%d/%d PV", hp, maxHP),
		75,
		568,
	)

	if c.State == TurnStatePlayer {
		c.DrawPlayerMenu(screen)
	} else {
		hudCenteredText(
			screen,
			c.Message,
			float64(renderWidth)/2,
			610,
		)
	}

	hudCenteredText(
		screen,
		c.Message,
		float64(renderWidth)/2,
		675,
	)
}

func (c *TurnCombat) DrawPlayerMenu(screen *ebiten.Image) {
	options := []string{
		"ATTAQUER",
		"RAYON D'ELDORIA",
		"DEFENDRE",
		fmt.Sprintf("OBJET  x%d", c.Potions),
	}

	startX := 495.0
	startY := 522.0

	for i, option := range options {
		y := startY + float64(i)*34

		if i == c.Selected {
			hudRect(
				screen,
				startX-24,
				y-5,
				310,
				27,
				color.RGBA{
					R: 80,
					G: 55,
					B: 20,
					A: 190,
				},
			)

			drawHUDText(
				screen,
				">",
				startX-15,
				y,
			)
		}

		drawHUDText(
			screen,
			option,
			startX+10,
			y,
		)
	}

	if c.Selected == TurnActionAttack {
		percent := 100

		if c.Skills != nil && c.Skills.Has(SkillPower3) {
			percent = 125
		}

		drawHUDText(
			screen,
			"ATTAQUE NORMALE",
			825,
			555,
		)

		drawHUDText(
			screen,
			fmt.Sprintf("%d%% DES DEGATS", percent),
			825,
			580,
		)
	}

	if c.Selected == TurnActionSkill {
		percent := 200

		if c.Skills != nil {
			percent = c.Skills.RayPercent()
		}

		drawHUDText(
			screen,
			"RAYON D'ENERGIE",
			825,
			555,
		)

		drawHUDText(
			screen,
			fmt.Sprintf("%d%% DES DEGATS", percent),
			825,
			580,
		)

		if c.SkillCooldown > 0 {
			drawHUDText(
				screen,
				fmt.Sprintf("RECHARGE : %d TOUR(S)", c.SkillCooldown),
				825,
				605,
			)
		}
	}

	if c.Selected == TurnActionDefend {
		if c.Skills != nil && c.Skills.Has(SkillDefense2) {
			drawHUDText(
				screen,
				"REDUIT 70% DES DEGATS",
				825,
				555,
			)
		} else {
			drawHUDText(
				screen,
				"REDUIT DE MOITIE",
				825,
				555,
			)
		}

		drawHUDText(
			screen,
			"LE PROCHAIN COUP",
			825,
			580,
		)
	}

	if c.Selected == TurnActionItem {
		heal := 30

		if c.Skills != nil {
			heal = c.Skills.PotionHeal()
		}

		drawHUDText(
			screen,
			fmt.Sprintf("RESTAURE %d PV", heal),
			825,
			555,
		)
	}

	drawHUDText(
		screen,
		"FLECHES : CHOISIR",
		825,
		630,
	)

	drawHUDText(
		screen,
		"ENTREE : VALIDER",
		825,
		652,
	)
}

func (c *TurnCombat) DrawEnergyScreenEffect(screen *ebiten.Image) {
	if !c.Active {
		return
	}

	if c.State != TurnStatePlayerAction {
		return
	}

	if c.PendingAction != TurnActionSkill {
		return
	}

	if c.Player == nil {
		return
	}

	if !c.Player.EnergyAttacking {
		return
	}

	if c.Player.CurrentEnergyFrame <= 3 {
		c.DrawEnergyChargeGlow(screen)
		return
	}

	c.DrawEnergyBeam(screen)

	if c.Player.CurrentEnergyFrame >= 6 {
		c.DrawEnergyImpact(screen)
	}
}

func (c *TurnCombat) DrawEnergyChargeGlow(screen *ebiten.Image) {
	timer := float64(c.Timer)

	pulse := (math.Sin(timer*0.7) + 1) / 2

	alpha := uint8(12 + pulse*18)

	hudRect(
		screen,
		0,
		0,
		renderWidth,
		renderHeight,
		color.RGBA{
			R: 20,
			G: 80,
			B: 180,
			A: alpha,
		},
	)
}

func (c *TurnCombat) DrawEnergyBeam(screen *ebiten.Image) {
	scaleX := float64(renderWidth) / 320
	scaleY := float64(renderHeight) / 180

	playerCenterX := (c.Player.X + c.Player.Width/2) * scaleX
	playerCenterY := (c.Player.VisualGroundY() - 31) * scaleY

	direction := 1.0

	if c.Player.FacingLeft {
		direction = -1
	}

	pulse := (math.Sin(float64(c.Timer)*0.8) + 1) / 2

	startX := playerCenterX + direction*70

	beamLength := 900.0

	if direction < 0 {
		startX -= beamLength
	}

	outerHeight := 70 + pulse*20
	middleHeight := 42 + pulse*12
	coreHeight := 18 + pulse*6

	ebitenutil.DrawRect(
		screen,
		startX,
		playerCenterY-outerHeight/2,
		beamLength,
		outerHeight,
		color.RGBA{
			R: 20,
			G: 85,
			B: 255,
			A: 45,
		},
	)

	ebitenutil.DrawRect(
		screen,
		startX,
		playerCenterY-middleHeight/2,
		beamLength,
		middleHeight,
		color.RGBA{
			R: 70,
			G: 175,
			B: 255,
			A: 110,
		},
	)

	ebitenutil.DrawRect(
		screen,
		startX,
		playerCenterY-coreHeight/2,
		beamLength,
		coreHeight,
		color.RGBA{
			R: 235,
			G: 250,
			B: 255,
			A: 210,
		},
	)
}

func (c *TurnCombat) DrawEnergyImpact(screen *ebiten.Image) {
	targetX := c.TargetCenterX() * float64(renderWidth) / 320

	targetY := (groundY - 42) * float64(renderHeight) / 180

	if c.Boss != nil {
		targetY = (groundY - c.Boss.HitHeight/2) * float64(renderHeight) / 180
	}

	pulse := (math.Sin(float64(c.Timer)*1.1) + 1) / 2

	size := 35 + pulse*35

	ebitenutil.DrawRect(
		screen,
		targetX-size/2,
		targetY-size/2,
		size,
		size,
		color.RGBA{
			R: 100,
			G: 200,
			B: 255,
			A: 70,
		},
	)

	innerSize := size * 0.40

	ebitenutil.DrawRect(
		screen,
		targetX-innerSize/2,
		targetY-innerSize/2,
		innerSize,
		innerSize,
		color.RGBA{
			R: 245,
			G: 255,
			B: 255,
			A: 220,
		},
	)
}
