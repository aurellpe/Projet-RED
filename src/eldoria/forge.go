package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ForgeItem struct {
	ID       string
	Name     string
	Cost     int
	HPBonus  int
	Category string
}

type ForgeProgress struct {
	Fragments    int             `json:"fragments"`
	Owned        map[string]bool `json:"owned"`
	BossRewarded bool            `json:"boss_rewarded"`
}

type ForgeMenu struct {
	Selected     int
	Progress     ForgeProgress
	Message      string
	MessageTimer int
}

var forgeItems = []ForgeItem{
	{
		ID:       "pull_brice",
		Name:     "PULL DE BRICE",
		Cost:     5,
		HPBonus:  80,
		Category: "ARMURE",
	},
	{
		ID:       "criniere_guigui",
		Name:     "CRINIERE DE GUIGUI",
		Cost:     3,
		HPBonus:  50,
		Category: "CASQUE",
	},
}

func newForgeProgress() ForgeProgress {
	return ForgeProgress{
		Fragments:    0,
		Owned:        map[string]bool{},
		BossRewarded: false,
	}
}

func getForgeSavePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "eldoria_forge.json"
	}

	saveDir := filepath.Join(
		configDir,
		"Eldoria",
	)

	if os.MkdirAll(saveDir, 0755) != nil {
		return "eldoria_forge.json"
	}

	return filepath.Join(
		saveDir,
		"forge.json",
	)
}

func LoadForgeProgress() ForgeProgress {
	progress := newForgeProgress()

	data, err := os.ReadFile(
		getForgeSavePath(),
	)

	if err != nil {
		return progress
	}

	if json.Unmarshal(data, &progress) != nil {
		return newForgeProgress()
	}

	if progress.Owned == nil {
		progress.Owned = map[string]bool{}
	}

	if progress.Fragments < 0 {
		progress.Fragments = 0
	}

	return progress
}

func SaveForgeProgress(progress ForgeProgress) error {
	if progress.Owned == nil {
		progress.Owned = map[string]bool{}
	}

	data, err := json.MarshalIndent(
		progress,
		"",
		"    ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		getForgeSavePath(),
		data,
		0644,
	)
}

func DeleteForgeProgress() error {
	err := os.Remove(
		getForgeSavePath(),
	)

	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func AddForgeFragments(amount int) {
	if amount <= 0 {
		return
	}

	progress := LoadForgeProgress()

	progress.Fragments += amount

	_ = SaveForgeProgress(progress)
}

func ForgeFragmentsForEnemy(enemyType int) int {
	switch enemyType {
	case EnemyTypeFast:
		return 2

	case EnemyTypeHeavy:
		return 3

	case EnemyTypeRanged:
		return 2

	default:
		return 1
	}
}

func RewardForgeBossOnce() {
	progress := LoadForgeProgress()

	if progress.BossRewarded {
		return
	}

	progress.BossRewarded = true
	progress.Fragments += 10

	_ = SaveForgeProgress(progress)
}

func ForgeHPBonus() int {
	progress := LoadForgeProgress()

	bonus := 0

	for _, item := range forgeItems {
		if progress.Owned[item.ID] {
			bonus += item.HPBonus
		}
	}

	return bonus
}

func NewForgeMenu() *ForgeMenu {
	menu := &ForgeMenu{
		Selected:     0,
		Message:      "",
		MessageTimer: 0,
	}

	menu.Reload()

	return menu
}

func (f *ForgeMenu) Reload() {
	f.Progress = LoadForgeProgress()

	if f.Selected < 0 || f.Selected >= len(forgeItems) {
		f.Selected = 0
	}

	f.Message = ""
	f.MessageTimer = 0
}

func (f *ForgeMenu) HasItem(item ForgeItem) bool {
	if f.Progress.Owned == nil {
		return false
	}

	return f.Progress.Owned[item.ID]
}

func (f *ForgeMenu) Update() bool {
	if f.MessageTimer > 0 {
		f.MessageTimer--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		return true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		f.Selected--

		if f.Selected < 0 {
			f.Selected = len(forgeItems) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		f.Selected++

		if f.Selected >= len(forgeItems) {
			f.Selected = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		f.BuySelected()
	}

	return false
}

func (f *ForgeMenu) BuySelected() {
	if f.Selected < 0 || f.Selected >= len(forgeItems) {
		return
	}

	item := forgeItems[f.Selected]

	if f.HasItem(item) {
		f.Message = "OBJET DEJA FORGE"
		f.MessageTimer = 100

		return
	}

	if f.Progress.Fragments < item.Cost {
		missing := item.Cost - f.Progress.Fragments

		f.Message = fmt.Sprintf(
			"IL MANQUE %d FRAGMENT(S)",
			missing,
		)

		f.MessageTimer = 100

		return
	}

	f.Progress.Fragments -= item.Cost

	if f.Progress.Owned == nil {
		f.Progress.Owned = map[string]bool{}
	}

	f.Progress.Owned[item.ID] = true

	_ = SaveForgeProgress(
		f.Progress,
	)

	f.Message = fmt.Sprintf(
		"%s FORGE ! +%d PV MAX",
		item.Name,
		item.HPBonus,
	)

	f.MessageTimer = 140
}

func (f *ForgeMenu) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		0,
		0,
		320,
		180,
		color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 150,
		},
	)

	panelColor := color.RGBA{
		R: 7,
		G: 8,
		B: 12,
		A: 245,
	}

	borderColor := color.RGBA{
		R: 205,
		G: 120,
		B: 45,
		A: 255,
	}

	ebitenutil.DrawRect(
		screen,
		32,
		9,
		256,
		162,
		borderColor,
	)

	ebitenutil.DrawRect(
		screen,
		33,
		10,
		254,
		160,
		panelColor,
	)

	drawCenteredMenuText(
		screen,
		"F O R G E",
		17,
	)

	drawCenteredMenuText(
		screen,
		fmt.Sprintf(
			"FRAGMENTS : %d",
			f.Progress.Fragments,
		),
		31,
	)

	for i, item := range forgeItems {
		y := 50 + i*46

		selected := f.Selected == i

		boxColor := color.RGBA{
			R: 20,
			G: 20,
			B: 24,
			A: 235,
		}

		itemBorder := color.RGBA{
			R: 80,
			G: 70,
			B: 60,
			A: 255,
		}

		if selected {
			boxColor = color.RGBA{
				R: 65,
				G: 35,
				B: 15,
				A: 240,
			}

			itemBorder = color.RGBA{
				R: 255,
				G: 150,
				B: 55,
				A: 255,
			}
		}

		ebitenutil.DrawRect(
			screen,
			48,
			float64(y),
			224,
			39,
			itemBorder,
		)

		ebitenutil.DrawRect(
			screen,
			49,
			float64(y+1),
			222,
			37,
			boxColor,
		)

		prefix := "  "

		if selected {
			prefix = "> "
		}

		ebitenutil.DebugPrintAt(
			screen,
			prefix+item.Name,
			57,
			y+5,
		)

		status := fmt.Sprintf(
			"%s  +%d PV  |  %d FRAG.",
			item.Category,
			item.HPBonus,
			item.Cost,
		)

		if f.HasItem(item) {
			status = item.Category + "  |  DEJA FORGE"
		}

		ebitenutil.DebugPrintAt(
			screen,
			status,
			57,
			y+20,
		)
	}

	message := f.Message

	if message == "" {
		message = "ENTREE : FORGER"
	}

	drawCenteredMenuText(
		screen,
		message,
		144,
	)

	drawCenteredMenuText(
		screen,
		"F OU ECHAP : FERMER",
		158,
	)
}
