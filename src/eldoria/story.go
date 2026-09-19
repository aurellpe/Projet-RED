package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	StoryKindNone = iota
	StoryKindIntro
	StoryKindBoss
)

type StoryPage struct {
	Speaker string
	Text    string
}

type StoryManager struct {
	Active bool
	Kind   int

	Pages     []StoryPage
	PageIndex int

	VisibleChars int
	TextTimer    int

	IntroSeen bool
	BossSeen  bool
}

var storyTextBuffer *ebiten.Image

func NewStoryManager() *StoryManager {
	return &StoryManager{
		Active:       false,
		Kind:         StoryKindNone,
		Pages:        []StoryPage{},
		PageIndex:    0,
		VisibleChars: 0,
		TextTimer:    0,
		IntroSeen:    false,
		BossSeen:     false,
	}
}

func (s *StoryManager) Reset() {
	s.Active = false
	s.Kind = StoryKindNone
	s.Pages = []StoryPage{}
	s.PageIndex = 0
	s.VisibleChars = 0
	s.TextTimer = 0
	s.IntroSeen = false
	s.BossSeen = false
}

func (s *StoryManager) Stop() {
	s.Active = false
	s.Kind = StoryKindNone
	s.Pages = []StoryPage{}
	s.PageIndex = 0
	s.VisibleChars = 0
	s.TextTimer = 0
}

func (s *StoryManager) StartIntro() {
	if s.IntroSeen {
		return
	}

	s.IntroSeen = true

	pages := []StoryPage{
		{
			Speaker: "NARRATEUR",
			Text:    "Autrefois, Eldoria etait un royaume de lumiere. Puis une ombre s'est abattue sur ses terres.",
		},
		{
			Speaker: "NARRATEUR",
			Text:    "Les routes sont tombees, les anciens sanctuaires se sont tus, et le Gardien a ferme le coeur du royaume.",
		},
		{
			Speaker: "NARRATEUR",
			Text:    "Tu es le dernier guerrier encore capable de traverser les ruines. Avance. Retrouve le Gardien. Libere Eldoria.",
		},
	}

	s.startStory(StoryKindIntro, pages)
}

func (s *StoryManager) StartBossIntro() {
	if s.BossSeen {
		return
	}

	s.BossSeen = true

	pages := []StoryPage{
		{
			Speaker: "GARDIEN",
			Text:    "Tu as traverse la foret, les ruines et le chateau. Peu d'humains sont arrives jusqu'ici.",
		},
		{
			Speaker: "HEROS",
			Text:    "Je ne suis pas venu pour survivre. Je suis venu pour rendre Eldoria a ceux que tu as condamnes.",
		},
		{
			Speaker: "GARDIEN",
			Text:    "Alors approche. Si tu veux la lumiere d'Eldoria, viens la reprendre de mes mains.",
		},
	}

	s.startStory(StoryKindBoss, pages)
}

func (s *StoryManager) startStory(kind int, pages []StoryPage) {
	s.Active = true
	s.Kind = kind
	s.Pages = pages
	s.PageIndex = 0
	s.VisibleChars = 0
	s.TextTimer = 0
}

func (s *StoryManager) currentPage() StoryPage {
	if s.PageIndex < 0 || s.PageIndex >= len(s.Pages) {
		return StoryPage{}
	}

	return s.Pages[s.PageIndex]
}

func (s *StoryManager) Update() {
	if !s.Active {
		return
	}

	page := s.currentPage()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.finish()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if s.VisibleChars < len(page.Text) {
			s.VisibleChars = len(page.Text)
			return
		}

		s.nextPage()
		return
	}

	s.TextTimer++

	if s.TextTimer >= 2 {
		s.TextTimer = 0

		if s.VisibleChars < len(page.Text) {
			s.VisibleChars++
		}
	}
}

func (s *StoryManager) nextPage() {
	s.PageIndex++

	if s.PageIndex >= len(s.Pages) {
		s.finish()
		return
	}

	s.VisibleChars = 0
	s.TextTimer = 0
}

func (s *StoryManager) finish() {
	s.Active = false
	s.Kind = StoryKindNone
	s.Pages = []StoryPage{}
	s.PageIndex = 0
	s.VisibleChars = 0
	s.TextTimer = 0
}

func (s *StoryManager) visibleText() string {
	page := s.currentPage()

	if s.VisibleChars >= len(page.Text) {
		return page.Text
	}

	if s.VisibleChars <= 0 {
		return ""
	}

	return page.Text[:s.VisibleChars]
}

func storyRect(screen *ebiten.Image, x float64, y float64, width float64, height float64, c color.RGBA) {
	ebitenutil.DrawRect(screen, x, y, width, height, c)
}

func drawStoryText(screen *ebiten.Image, text string, x float64, y float64, scale float64) {
	if storyTextBuffer == nil {
		storyTextBuffer = ebiten.NewImage(1000, 32)
	}

	storyTextBuffer.Clear()

	ebitenutil.DebugPrintAt(storyTextBuffer, text, 0, 0)

	options := &ebiten.DrawImageOptions{}

	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(x, y)
	options.Filter = ebiten.FilterNearest

	screen.DrawImage(storyTextBuffer, options)
}

func storyCenteredText(screen *ebiten.Image, text string, centerX float64, y float64, scale float64) {
	textWidth := float64(len(text)*6) * scale
	x := centerX - textWidth/2

	drawStoryText(screen, text, x, y, scale)
}

func wrapStoryText(text string, maxCharacters int) []string {
	words := strings.Fields(text)

	if len(words) == 0 {
		return []string{}
	}

	lines := []string{}
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
			continue
		}

		testLine := currentLine + " " + word

		if len(testLine) > maxCharacters {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (s *StoryManager) Draw(screen *ebiten.Image) {
	if !s.Active {
		return
	}

	storyRect(screen, 0, 0, renderWidth, renderHeight, color.RGBA{R: 0, G: 0, B: 0, A: 115})

	storyRect(screen, 0, 0, renderWidth, 86, color.RGBA{R: 0, G: 0, B: 0, A: 210})
	storyRect(screen, 0, 425, renderWidth, 295, color.RGBA{R: 0, G: 0, B: 0, A: 225})

	title := "PROLOGUE"

	if s.Kind == StoryKindBoss {
		title = "LE GARDIEN D'ELDORIA"
	}

	storyCenteredText(screen, title, float64(renderWidth)/2, 30, 2.0)

	storyRect(screen, 150, 78, 980, 2, color.RGBA{R: 145, G: 115, B: 65, A: 180})

	page := s.currentPage()

	speakerWidth := float64(len(page.Speaker)*6)*1.7 + 38

	storyRect(screen, 110, 448, speakerWidth, 38, color.RGBA{R: 25, G: 19, B: 12, A: 240})
	storyRect(screen, 110, 448, speakerWidth, 2, color.RGBA{R: 200, G: 155, B: 75, A: 220})

	drawStoryText(screen, page.Speaker, 128, 458, 1.7)

	visibleText := s.visibleText()
	lines := wrapStoryText(visibleText, 78)

	startY := 515.0

	for i, line := range lines {
		drawStoryText(screen, line, 120, startY+float64(i)*31, 1.65)
	}

	hint := "ENTREE / ESPACE : SUIVANT"

	if s.VisibleChars < len(page.Text) {
		hint = "ENTREE / ESPACE : AFFICHER LE TEXTE"
	}

	storyCenteredText(screen, hint, float64(renderWidth)/2, 661, 1.05)
	storyCenteredText(screen, "ECHAP : PASSER", float64(renderWidth)/2, 685, 0.95)

	s.drawPageIndicators(screen)
}

func (s *StoryManager) drawPageIndicators(screen *ebiten.Image) {
	if len(s.Pages) <= 1 {
		return
	}

	totalWidth := float64(len(s.Pages)*18 - 8)
	startX := float64(renderWidth)/2 - totalWidth/2

	for i := 0; i < len(s.Pages); i++ {
		size := 6.0
		alpha := uint8(90)

		if i == s.PageIndex {
			size = 10
			alpha = 230
		}

		x := startX + float64(i)*18
		y := 625.0

		storyRect(screen, x, y, size, size, color.RGBA{R: 220, G: 175, B: 90, A: alpha})
	}
}
