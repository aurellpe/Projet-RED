package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	potionVie       = "Potion de vie"
	potionPoison    = "Potion de poison"
	potionMana      = "Potion de mana"
	livreBouleDeFeu = "Livre de Sort : Boule de Feu"
	augmentationInv = "Augmentation d'inventaire"
	fourrureYael    = "fourrure de Yael"
	peauAurelien    = "Peau d'Aurelien"
	cuirSanglier    = "Cuir de Sanglier"
	plumeCorbeau    = "Plume de Corbeau"
	chapeau         = "Chapeau de l'aventurier"
	tunique         = "Tunique de l'aventurier"
	bottes          = "Bottes de l'aventurier"

	pichenette = "pichenette"
	sortBouleDeFeu  = "Boule de Feu"

	coutForge = 5
)

type Spell struct {
	Nom    string
	Degats int
	Mana   int
}

var grimoire = map[string]Spell{
	pichenette: {Nom: pichenette, Degats: 8, Mana: 5},
	sortBouleDeFeu:  {Nom: sortBouleDeFeu, Degats: 18, Mana: 15},
}

var equipBonus = map[string]int{
	chapeau: 10,
	tunique: 25,
	bottes:  15,
}

var lecteur = bufio.NewReader(os.Stdin)

func initCharacter(
	nom string,
	prenom string,
	age int,
	classe string,
	niveau int,
	pointsDeVieMaximum int,
	pointsDeVieActuels int,
	inventaire []structure.Item,
	argent int,
	xpmax int,
	xpactu int,
	degats int,
) structure.Character {

	return structure.Character{
		Nom:                nom,
		Prenom:             prenom,
		Age:                age,
		Classe:             classe,
		Niveau:             niveau,
		PointsDeVieMaximum: pointsDeVieMaximum,
		PointsDeVieActuels: pointsDeVieActuels,
		Inventaire:         inventaire,
		Argent:             argent,
		XPmax:              xpmax,
		XPactu:             xpactu,
		Degats:             degats,
		InventaireMax: 10,
		Skills:        []string{pichenette},
		Mana:          50,
		ManaMax:       50,
		Initiative:    10,
	}
}

func displayInfo(player structure.Character) {

	fmt.Println()
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", player.Nom)
	fmt.Println("Prénom :", player.Prenom)
	fmt.Println("Âge :", player.Age)
	fmt.Println("Classe :", player.Classe)
	fmt.Println("Niveau :", player.Niveau)
	fmt.Println("Points de vie :", player.PointsDeVieActuels, "/", player.PointsDeVieMaximum)
	fmt.Println("Arme :", player.Arme)
	fmt.Println("Dégâts :", player.Degats)
	fmt.Println("Argent :", player.Argent)
	fmt.Println("Débris :", player.Debris)
	fmt.Println("Mana :", player.Mana, "/", player.ManaMax)
	fmt.Println("XP :", player.XPactu, "/", player.XPmax)
	fmt.Println("Initiative :", player.Initiative)
	fmt.Println("Sorts :", strings.Join(player.Skills, ", "))
	fmt.Printf("Équipement : tête [%s] | torse [%s] | pieds [%s]\n",
		orEmpty(player.Equipement.Tete), orEmpty(player.Equipement.Torse), orEmpty(player.Equipement.Pieds))
}

func pause() {

	fmt.Println("\nAppuyez sur Entrée pour continuer...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func accessInventory(player structure.Character) {

	fmt.Println()
	fmt.Println("===== INVENTAIRE =====")

	fmt.Printf("Emplacements : %d / %d\n", len(player.Inventaire), player.InventaireMax)
	fmt.Println()

	if len(player.Inventaire) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	for i, objet := range player.Inventaire {

		fmt.Println(
			i+1,
			"-",
			objet.Nom,
			"(x",
			objet.Quantity,
			")",
		)
	}
}

func main() {

	player := createCharacter()

	var run bool = true

	for run {

		run = menu(&player)

	}
}

func takePot(p *structure.Character) bool {

	if countItem(p, potionVie) == 0 {
		fmt.Println("Vous n'avez pas de potion de vie.")
		return false
	}

	if p.PointsDeVieActuels >= p.PointsDeVieMaximum {
		fmt.Println("Vos PV sont déjà au maximum.")
		return false
	}

	removeInventory(p, potionVie)

	p.PointsDeVieActuels += 50

	if p.PointsDeVieActuels > p.PointsDeVieMaximum {
		p.PointsDeVieActuels = p.PointsDeVieMaximum
	}

	fmt.Printf("Vous utilisez %s. PV : %d / %d\n", potionVie, p.PointsDeVieActuels, p.PointsDeVieMaximum)
	return true
}

func shopMenu(c *structure.Character, scanner *bufio.Scanner) {

	for {

		fmt.Println("\n===== Marchand =====")
		fmt.Println("1. Potion de vie")
		fmt.Println("2. Fourrure de Brice")
		fmt.Println("3. Casque de Guigui")
		fmt.Println("0. Retour")

		scanner.Scan()

		choix := scanner.Text()

		switch choix {

		case "1":
			acheterItem(c, "Potion de vie")

		case "2":
			acheterItem(c, "Fourrure de Brice")

		case "3":
			acheterItem(c, "Casque de Guigui")

		case "0":
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func acheterItem(c *structure.Character, nom string) {

	fmt.Println("Vous avez acheté :", nom)
}

type shopItem struct {
	Nom  string
	Prix int
}

var shopItems = []shopItem{

	{Nom: potionVie, Prix: 3},
	{Nom: potionPoison, Prix: 6},
	{Nom: potionMana, Prix: 5},
	{Nom: livreBouleDeFeu, Prix: 25},
	{Nom: fourrureYael, Prix: 4},
	{Nom: peauAurelien, Prix: 7},
	{Nom: cuirSanglier, Prix: 3},
	{Nom: plumeCorbeau, Prix: 1},
	{Nom: augmentationInv, Prix: 30},
}

func shop(c *structure.Character, reader *bufio.Reader) {

	for {

		fmt.Println("\n===== Marchand =====")

		for i, si := range shopItems {

			fmt.Printf(
				"%d. %s (%d pièces d'or)\n",
				i+1,
				si.Nom,
				si.Prix,
			)
		}

		fmt.Println("0. Retour")
		fmt.Println("Votre argent :", c.Argent, "pièces d'or")

		choix := lireChoix(reader)

		if choix == 0 {
			return
		}

		if choix >= 1 && choix <= len(shopItems) {

			si := shopItems[choix-1]

			acheterItemPayant(
				c,
				si.Nom,
				si.Prix,
			)

		} else {

			fmt.Println("Choix invalide.")
		}
	}
}

func acheterItemPayant(
	c *structure.Character,
	nom string,
	prix int,
) {

	if c.Argent < prix {

		fmt.Println(
			"Vous n'avez pas assez d'argent pour acheter :",
			nom,
		)

		return
	}

	if !c.AddInventory(structure.Item{
		Nom:      nom,
		Quantity: 1,
	}) {

		fmt.Println("Votre inventaire est plein !")
		fmt.Printf("Maximum : %d emplacements.\n", c.InventaireMax)

		return
	}

	c.Argent -= prix

	fmt.Println("Vous avez acheté :", nom)
	fmt.Println("Argent restant :", c.Argent)
}

func gainExperience(c *structure.Character, xpGagne int) {

	fmt.Printf(
		"Vous gagnez %d points d'expérience.\n",
		xpGagne,
	)

	c.XPactu += xpGagne

	for c.XPactu >= c.XPmax {

		c.XPactu -= c.XPmax

		levelUp(c)
	}

	fmt.Printf("XP : %d / %d\n", c.XPactu, c.XPmax)
}

func levelUp(c *structure.Character) {

	c.Niveau++

	bonusPV := 10
	bonusMana := 5

	c.PointsDeVieMaximum += bonusPV
	c.PointsDeVieActuels += bonusPV
	c.ManaMax += bonusMana
	c.Mana += bonusMana

	c.XPmax = int(float64(c.XPmax) * 1.2)

	fmt.Printf(
		"Niveau supérieur ! Vous êtes maintenant niveau %d.\n",
		c.Niveau,
	)

	fmt.Printf(
		"Bonus : +%d points de vie maximum, +%d mana maximum.\n",
		bonusPV,
		bonusMana,
	)

	fmt.Printf(
		"PV : %d / %d\n",
		c.PointsDeVieActuels,
		c.PointsDeVieMaximum,
	)
}

func createCharacter() structure.Character {

	fmt.Println("===== CRÉATION DU PERSONNAGE =====")

	nom := lireNomValide("Nom")
	prenom := lireNomValide("Prénom")

	age := 0
	for age <= 0 {
		fmt.Print("Âge : ")
		ligne, _ := lecteur.ReadString('\n')
		a, err := strconv.Atoi(strings.TrimSpace(ligne))
		if err == nil && a > 0 {
			age = a
		} else {
			fmt.Println("Âge invalide.")
		}
	}

	classe := ""
	pvMax := 0
	initiative := 0

	for classe == "" {
		fmt.Println("Classe :")
		fmt.Println("1. Humain (100 PV max)")
		fmt.Println("2. ogre   (80 PV max)")
		fmt.Println("3. chevalier   (120 PV max)")

		switch lireEntier() {
		case 1:
			classe, pvMax, initiative = "Humain", 100, 10
		case 2:
			classe, pvMax, initiative = "ogre", 80, 12
		case 3:
			classe, pvMax, initiative = "chevalier", 120, 8
		default:
			fmt.Println("Choix invalide.")
		}
	}

	inventaire := []structure.Item{{Nom: potionVie, Quantity: 3}}

	p := initCharacter(
		nom,
		prenom,
		age,
		classe,
		1,
		pvMax,
		pvMax/2,
		inventaire,
		100, 
		100, 
		0,
		15, 
	)

	p.Initiative = initiative

	return p
}

func formatName(s string) string {
	r := []rune(strings.ToLower(s))
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func validName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func lireNomValide(label string) string {
	for {
		fmt.Printf("%s (lettres uniquement) : ", label)
		ligne, _ := lecteur.ReadString('\n')
		ligne = strings.TrimSpace(ligne)
		if validName(ligne) {
			return formatName(ligne)
		}
		fmt.Println("Invalide : uniquement des lettres (sans espace ni chiffre).")
	}
}

func orEmpty(s string) string {
	if s == "" {
		return "vide"
	}
	return s
}

func countItem(p *structure.Character, nom string) int {
	for _, it := range p.Inventaire {
		if it.Nom == nom {
			return it.Quantity
		}
	}
	return 0
}

func removeInventory(p *structure.Character, nom string) bool {
	for i := range p.Inventaire {
		if p.Inventaire[i].Nom == nom && p.Inventaire[i].Quantity > 0 {
			p.Inventaire[i].Quantity--
			if p.Inventaire[i].Quantity == 0 {
				p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
			}
			return true
		}
	}
	return false
}

func poisonPot(nom string, pv *int, pvMax int) {
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		*pv -= 10
		if *pv < 0 {
			*pv = 0
		}
		fmt.Printf("Poison : %s subit 10 dégâts. PV : %d / %d\n", nom, *pv, pvMax)
		if *pv == 0 {
			break
		}
	}
}

func takeManaPot(p *structure.Character) bool {
	if p.Mana >= p.ManaMax {
		fmt.Println("Votre mana est déjà au maximum.")
		return false
	}
	removeInventory(p, potionMana)
	p.Mana += 30
	if p.Mana > p.ManaMax {
		p.Mana = p.ManaMax
	}
	fmt.Printf("Vous utilisez %s. Mana : %d / %d\n", potionMana, p.Mana, p.ManaMax)
	return true
}

func knowsSpell(p *structure.Character, sort string) bool {
	for _, s := range p.Skills {
		if s == sort {
			return true
		}
	}
	return false
}

func spellBook(p *structure.Character) bool {
	if knowsSpell(p, sortBouleDeFeu) {
		fmt.Println("Vous connaissez déjà le sort Boule de Feu.")
		return false
	}
	removeInventory(p, livreBouleDeFeu)
	p.Skills = append(p.Skills, sortBouleDeFeu)
	fmt.Println("Vous apprenez le sort : Boule de Feu !")
	return true
}

func upgradeInventorySlot(p *structure.Character) bool {
	if p.Ameliorations >= 3 {
		fmt.Println("Vous avez déjà utilisé les 3 augmentations d'inventaire possibles.")
		return false
	}
	removeInventory(p, augmentationInv)
	p.Ameliorations++
	p.InventaireMax += 10
	fmt.Printf("Capacité de l'inventaire : %d (%d / 3 améliorations)\n", p.InventaireMax, p.Ameliorations)
	return true
}

func equip(p *structure.Character, nom string) bool {
	var slot *string
	switch nom {
	case chapeau:
		slot = &p.Equipement.Tete
	case tunique:
		slot = &p.Equipement.Torse
	case bottes:
		slot = &p.Equipement.Pieds
	default:
		return false
	}

	ancien := *slot
	removeInventory(p, nom)

	if ancien != "" {
		if !p.AddInventory(structure.Item{Nom: ancien, Quantity: 1}) {
			p.AddInventory(structure.Item{Nom: nom, Quantity: 1})
			fmt.Println("Inventaire plein, impossible de retirer l'ancien équipement.")
			return false
		}
		p.PointsDeVieMaximum -= equipBonus[ancien]
		fmt.Printf("Vous retirez %s (remis dans l'inventaire).\n", ancien)
	}

	*slot = nom
	p.PointsDeVieMaximum += equipBonus[nom]
	if p.PointsDeVieActuels > p.PointsDeVieMaximum {
		p.PointsDeVieActuels = p.PointsDeVieMaximum
	}
	fmt.Printf("Vous équipez %s (+%d PV max). PV : %d / %d\n", nom, equipBonus[nom],
		p.PointsDeVieActuels, p.PointsDeVieMaximum)
	return true
}

func useItem(p *structure.Character, nom string, m *Monster) bool {
	switch nom {
	case potionVie:
		return takePot(p)
	case potionMana:
		return takeManaPot(p)
	case potionPoison:
		removeInventory(p, potionPoison)
		if m != nil {
			fmt.Printf("Vous lancez %s sur %s !\n", potionPoison, m.Nom)
			poisonPot(m.Nom, &m.PV, m.PVMax)
		} else {
			fmt.Println("Vous buvez la potion de poison")
			poisonPot(p.Nom, &p.PointsDeVieActuels, p.PointsDeVieMaximum)
			isDead(p)
		}
		return true
	case livreBouleDeFeu:
		return spellBook(p)
	case augmentationInv:
		return upgradeInventorySlot(p)
	case chapeau, tunique, bottes:
		return equip(p, nom)
	default:
		fmt.Printf("%s n'est pas utilisable (matériau de fabrication).\n", nom)
		return false
	}
}

func inventoryMenu(p *structure.Character, m *Monster) bool {
	for {
		accessInventory(*p)

		if len(p.Inventaire) == 0 {
			return false
		}

		fmt.Println("\nTapez le numéro d'un objet pour l'utiliser, 0 pour revenir.")
		choix := lireEntier()

		if choix == 0 {
			return false
		}
		if choix < 1 || choix > len(p.Inventaire) {
			fmt.Println("Choix invalide.")
			continue
		}

		nom := p.Inventaire[choix-1].Nom
		if useItem(p, nom, m) && m != nil {
			return true
		}
	}
}

// ===================== AJOUT : menu principal =====================

func whoAreThey() {
	fmt.Println("\n===== Qui sont-ils ? =====")
	fmt.Println("Partie 2 (économie) : ABBA")
	fmt.Println("Partie 3 (combat)   : Steven Spielberg")
}
