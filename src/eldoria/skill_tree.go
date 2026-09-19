package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SkillBranchPower = iota
	SkillBranchDefense
	SkillBranchEnergy
)

const (
	SkillPower1   = "power_1"
	SkillPower2   = "power_2"
	SkillPower3   = "power_3"
	SkillDefense1 = "defense_1"
	SkillDefense2 = "defense_2"
	SkillDefense3 = "defense_3"
	SkillEnergy1  = "energy_1"
	SkillEnergy2  = "energy_2"
	SkillEnergy3  = "energy_3"
)

type SkillDefinition struct {
	ID          string
	Name        string
	Description string
}

type SkillTreeSave struct {
	Points               int             `json:"points"`
	Unlocked             map[string]bool `json:"unlocked"`
	HighestRewardedLevel int             `json:"highest_rewarded_level"`
}

type SkillTree struct {
	Points               int
	Unlocked             map[string]bool
	HighestRewardedLevel int
	Open                 bool
	SelectedBranch       int
	SelectedRow          int
	Message              string
	MessageTimer         int
}

var skillBranches = [][]SkillDefinition{
	{
		{ID: SkillPower1, Name: "FORCE I", Description: "ATTAQUE +10%"},
		{ID: SkillPower2, Name: "FORCE II", Description: "ATTAQUE +10%"},
		{ID: SkillPower3, Name: "FRAPPE BRUTALE", Description: "ATTAQUE NORMALE +25%"},
	},
	{
		{ID: SkillDefense1, Name: "VITALITE I", Description: "+20 PV MAX"},
		{ID: SkillDefense2, Name: "GARDE RENFORCEE", Description: "GARDE REDUIT 70%"},
		{ID: SkillDefense3, Name: "POTION +", Description: "POTION REND 50 PV"},
	},
	{
		{ID: SkillEnergy1, Name: "RAYON I", Description: "RAYON = 225%"},
		{ID: SkillEnergy2, Name: "MAITRISE", Description: "RECHARGE PLUS COURTE"},
		{ID: SkillEnergy3, Name: "SURCHARGE", Description: "RAYON = 300%"},
	},
}

func NewSkillTree() *SkillTree {
	return &SkillTree{
		Points:               0,
		Unlocked:             map[string]bool{},
		HighestRewardedLevel: -1,
		Open:                 false,
		SelectedBranch:       0,
		SelectedRow:          0,
		Message:              "",
		MessageTimer:         0,
	}
}

func skillTreeSavePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "eldoria_skills.json"
	}

	dir := filepath.Join(configDir, "Eldoria")
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, "skills.json")
}

func LoadSkillTreeProgress() *SkillTree {
	tree := NewSkillTree()

	data, err := os.ReadFile(skillTreeSavePath())
	if err != nil {
		return tree
	}

	var save SkillTreeSave

	if json.Unmarshal(data, &save) != nil {
		return tree
	}

	tree.Points = save.Points
	tree.HighestRewardedLevel = save.HighestRewardedLevel

	if save.Unlocked != nil {
		tree.Unlocked = save.Unlocked
	}

	return tree
}

func SaveSkillTreeProgress(tree *SkillTree) error {
	if tree == nil {
		return nil
	}

	save := SkillTreeSave{
		Points:               tree.Points,
		Unlocked:             tree.Unlocked,
		HighestRewardedLevel: tree.HighestRewardedLevel,
	}

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(skillTreeSavePath(), data, 0644)
}

func DeleteSkillTreeProgress() error {
	err := os.Remove(skillTreeSavePath())
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func (s *SkillTree) Has(skillID string) bool {
	if s == nil || s.Unlocked == nil {
		return false
	}

	return s.Unlocked[skillID]
}

func (s *SkillTree) OpenMenu() {
	if s == nil {
		return
	}

	s.Open = true
	s.Message = "CHOISIS UNE COMPETENCE"
	s.MessageTimer = 0
}

func (s *SkillTree) CloseMenu() {
	if s == nil {
		return
	}

	s.Open = false
	_ = SaveSkillTreeProgress(s)
}

func (s *SkillTree) Update(game *Game) {
	if s == nil || !s.Open {
		return
	}

	if s.MessageTimer > 0 {
		s.MessageTimer--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.CloseMenu()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		s.SelectedBranch--
		if s.SelectedBranch < 0 {
			s.SelectedBranch = 2
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		s.SelectedBranch++
		if s.SelectedBranch > 2 {
			s.SelectedBranch = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.SelectedRow--
		if s.SelectedRow < 0 {
			s.SelectedRow = 2
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.SelectedRow++
		if s.SelectedRow > 2 {
			s.SelectedRow = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if s.UnlockSelected() {
			if game != nil {
				game.ApplyCurrentLevelProgression()
			}

			_ = SaveSkillTreeProgress(s)
		}
	}
}

func (s *SkillTree) UnlockSelected() bool {
	if s == nil {
		return false
	}

	if s.SelectedBranch < 0 || s.SelectedBranch >= len(skillBranches) {
		return false
	}

	branch := skillBranches[s.SelectedBranch]

	if s.SelectedRow < 0 || s.SelectedRow >= len(branch) {
		return false
	}

	skill := branch[s.SelectedRow]

	if s.Has(skill.ID) {
		s.SetMessage("COMPETENCE DEJA DEBLOQUEE")
		return false
	}

	if s.SelectedRow > 0 {
		previousSkill := branch[s.SelectedRow-1]

		if !s.Has(previousSkill.ID) {
			s.SetMessage("DEBLOQUE D'ABORD LA COMPETENCE PRECEDENTE")
			return false
		}
	}

	if s.Points <= 0 {
		s.SetMessage("PAS ASSEZ DE POINTS")
		return false
	}

	s.Points--
	s.Unlocked[skill.ID] = true
	s.SetMessage(skill.Name + " DEBLOQUEE")

	return true
}

func (s *SkillTree) SetMessage(message string) {
	s.Message = message
	s.MessageTimer = 120
}

func (s *SkillTree) AwardLevelCompletion(levelIndex int) bool {
	if s == nil {
		return false
	}

	if levelIndex < 0 || levelIndex > 2 {
		return false
	}

	if levelIndex <= s.HighestRewardedLevel {
		return false
	}

	s.Points += 2
	s.HighestRewardedLevel = levelIndex
	s.SetMessage("+2 POINTS DE COMPETENCE")

	_ = SaveSkillTreeProgress(s)

	return true
}

func (s *SkillTree) GrantCatchUpRewards(currentLevel int) {
	if s == nil {
		return
	}

	for levelIndex := 0; levelIndex < currentLevel && levelIndex <= 2; levelIndex++ {
		s.AwardLevelCompletion(levelIndex)
	}
}

func (s *SkillTree) MaxHPBonus() int {
	if s != nil && s.Has(SkillDefense1) {
		return 20
	}

	return 0
}

func (s *SkillTree) ApplyAttackBonus(baseDamage int) int {
	if s == nil {
		return baseDamage
	}

	damage := baseDamage

	if s.Has(SkillPower1) {
		damage += baseDamage / 10
	}

	if s.Has(SkillPower2) {
		damage += baseDamage / 10
	}

	return damage
}

func (s *SkillTree) NormalAttackDamage(baseDamage int) int {
	if s != nil && s.Has(SkillPower3) {
		return baseDamage + baseDamage/4
	}

	return baseDamage
}

func (s *SkillTree) GuardDamage(damage int) int {
	if s != nil && s.Has(SkillDefense2) {
		reduced := (damage*3 + 9) / 10

		if reduced < 1 {
			return 1
		}

		return reduced
	}

	reduced := (damage + 1) / 2

	if reduced < 1 {
		return 1
	}

	return reduced
}

func (s *SkillTree) PotionHeal() int {
	if s != nil && s.Has(SkillDefense3) {
		return 50
	}

	return 30
}

func (s *SkillTree) RayDamage(baseDamage int) int {
	if s != nil && s.Has(SkillEnergy3) {
		return baseDamage * 3
	}

	if s != nil && s.Has(SkillEnergy1) {
		return baseDamage * 9 / 4
	}

	return baseDamage * 2
}

func (s *SkillTree) RayPercent() int {
	if s != nil && s.Has(SkillEnergy3) {
		return 300
	}

	if s != nil && s.Has(SkillEnergy1) {
		return 225
	}

	return 200
}

func (s *SkillTree) RayCooldown() int {
	if s != nil && s.Has(SkillEnergy2) {
		return 2
	}

	return 3
}

func (s *SkillTree) Draw(screen *ebiten.Image) {
	if s == nil || !s.Open {
		return
	}

	hudRect(screen, 0, 0, renderWidth, renderHeight, color.RGBA{R: 0, G: 0, B: 0, A: 205})

	hudRect(
		screen,
		80,
		45,
		1120,
		630,
		color.RGBA{R: 7, G: 9, B: 15, A: 248},
	)

	hudRect(
		screen,
		80,
		45,
		1120,
		4,
		color.RGBA{R: 205, G: 155, B: 65, A: 255},
	)

	hudCenteredText(
		screen,
		"ARBRE DE COMPETENCES",
		float64(renderWidth)/2,
		72,
	)

	hudCenteredText(
		screen,
		fmt.Sprintf("POINTS DISPONIBLES : %d", s.Points),
		float64(renderWidth)/2,
		105,
	)

	branchNames := []string{
		"PUISSANCE",
		"DEFENSE",
		"ENERGIE",
	}

	branchX := []float64{
		150,
		500,
		850,
	}

	for branchIndex := 0; branchIndex < 3; branchIndex++ {
		hudCenteredText(
			screen,
			branchNames[branchIndex],
			branchX[branchIndex]+140,
			145,
		)

		for row := 0; row < 3; row++ {
			s.drawSkillNode(
				screen,
				branchIndex,
				row,
				branchX[branchIndex],
				185+float64(row)*125,
			)
		}
	}

	message := s.Message

	if message == "" {
		message = "ENTREE : DEBLOQUER"
	}

	hudCenteredText(
		screen,
		message,
		float64(renderWidth)/2,
		584,
	)

	hudCenteredText(
		screen,
		"FLECHES : NAVIGUER   |   ENTREE : ACHETER   |   C / ECHAP : FERMER",
		float64(renderWidth)/2,
		630,
	)
}

func (s *SkillTree) drawSkillNode(screen *ebiten.Image, branch int, row int, x float64, y float64) {
	skill := skillBranches[branch][row]

	selected := s.SelectedBranch == branch && s.SelectedRow == row
	unlocked := s.Has(skill.ID)
	available := row == 0 || s.Has(skillBranches[branch][row-1].ID)

	outerColor := color.RGBA{R: 45, G: 48, B: 58, A: 255}
	innerColor := color.RGBA{R: 18, G: 21, B: 30, A: 255}
	status := "VERROUILLE"

	if available {
		outerColor = color.RGBA{R: 65, G: 95, B: 140, A: 255}
		status = "1 POINT"
	}

	if unlocked {
		outerColor = color.RGBA{R: 55, G: 135, B: 80, A: 255}
		innerColor = color.RGBA{R: 14, G: 40, B: 25, A: 255}
		status = "DEBLOQUEE"
	}

	if selected {
		outerColor = color.RGBA{R: 220, G: 170, B: 70, A: 255}
	}

	hudRect(
		screen,
		x,
		y,
		280,
		95,
		outerColor,
	)

	hudRect(
		screen,
		x+3,
		y+3,
		274,
		89,
		innerColor,
	)

	hudCenteredText(
		screen,
		skill.Name,
		x+140,
		y+17,
	)

	hudCenteredText(
		screen,
		skill.Description,
		x+140,
		y+43,
	)

	hudCenteredText(
		screen,
		status,
		x+140,
		y+68,
	)

	if row < 2 {
		hudRect(
			screen,
			x+138,
			y+95,
			4,
			30,
			color.RGBA{R: 80, G: 85, B: 100, A: 255},
		)
	}
}
