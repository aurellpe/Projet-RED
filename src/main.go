package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"os"
	"strings"
	"unicode"
)

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
		XPmax:              100,
		XPactu:             0,
		Argent:             argent,
	}
}

func characterCreation(scanner *bufio.Scanner) structure.Character {
	var nom string

	fmt.Println()
	fmt.Println("===== CRÉATION DU PERSONNAGE =====")

	for {
		fmt.Print("Choisissez le nom de votre personnage : ")

		scanner.Scan()
		nom = strings.TrimSpace(scanner.Text())

		nomValide := true

		if nom == "" {
			nomValide = false
		}

		for _, lettre := range nom {
			if !unicode.IsLetter(lettre) {
				nomValide = false
				break
			}
		}

		if nomValide {
			break
		}

		fmt.Println("Nom invalide.")
		fmt.Println("Le nom doit contenir uniquement des lettres.")
	}

	nom = strings.ToLower(nom)

	runes := []rune(nom)
	runes[0] = unicode.ToUpper(runes[0])

	nom = string(runes)

	fmt.Println()
	fmt.Println("Votre personnage s'appelle :", nom)

	return initCharacter(
		nom,
		"",
		26,
		"Légionnaire",
		1,
		10000,
		150,
		[]structure.Item{
			{Nom: "Potion de vie", Quantity: 1},
			{Nom: "Épée", Quantity: 1},
		},
		100,
	)
}

func displayInfo(perso structure.Character) {
	fmt.Println()
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", perso.Nom)
	fmt.Println("Âge :", perso.Age)
	fmt.Println("Classe :", perso.Classe)
	fmt.Println("Niveau :", perso.Niveau)
	fmt.Println("PV :", perso.PointsDeVieActuels, "/", perso.PointsDeVieMaximum)
	fmt.Println("XP :", perso.XPactu, "/", perso.XPmax)
	fmt.Println("Argent :", perso.Argent)
}

func accessInventory(perso structure.Character) {
	fmt.Println()
	fmt.Println("===== INVENTAIRE =====")

	if len(perso.Inventaire) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	for i, objet := range perso.Inventaire {
		fmt.Printf("%d - %s x%d\n", i+1, objet.Nom, objet.Quantity)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Création du personnage par le joueur
	perso := characterCreation(scanner)

	// Création du boss
	boss := structure.NewBoss()

	run := true

	for run {
		fmt.Println()
		fmt.Println("========== MENU ==========")
		fmt.Println("1. Informations du personnage")
		fmt.Println("2. Inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Combattre le boss")
		fmt.Println("0. Quitter")
		fmt.Print("Choix : ")

		scanner.Scan()
		choix := scanner.Text()

		switch choix {
		case "1":
			displayInfo(perso)

		case "2":
			accessInventory(perso)

		case "3":
			shopMenu(&perso, scanner)

		case "4":
			combatBoss(&perso, &boss, scanner)

		case "0":
			fmt.Println("Au revoir", perso.Nom, "!")
			run = false

		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func takePot(p *structure.Character) {
	for i, item := range p.Inventaire {
		if item.Nom == "Potion de vie" && item.Quantity > 0 {
			p.PointsDeVieActuels += 50

			if p.PointsDeVieActuels > p.PointsDeVieMaximum {
				p.PointsDeVieActuels = p.PointsDeVieMaximum
			}

			p.Inventaire[i].Quantity--

			fmt.Println()
			fmt.Println("Vous utilisez une potion de vie.")
			fmt.Println("PV :", p.PointsDeVieActuels, "/", p.PointsDeVieMaximum)

			if p.Inventaire[i].Quantity == 0 {
				p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
			}

			return
		}
	}

	fmt.Println("Vous n'avez pas de potion de vie.")
}

type shopItem struct {
	Nom  string
	Prix int
}

var shopItems = []shopItem{
	{Nom: "Potion de vie", Prix: 3},
	{Nom: "Potion de poison", Prix: 6},
	{Nom: "Livre de Sort : Boule de Feu", Prix: 25},
	{Nom: "Fourrure de Loup", Prix: 4},
	{Nom: "Peau de Troll", Prix: 7},
	{Nom: "Cuir de vache", Prix: 1},
}

func shopMenu(c *structure.Character, scanner *bufio.Scanner) {
	for {
		fmt.Println()
		fmt.Println("===== MARCHAND =====")

		for i, item := range shopItems {
			fmt.Printf("%d. %s (%d pièces d'or)\n", i+1, item.Nom, item.Prix)
		}

		fmt.Println("0. Retour")
		fmt.Println("Votre argent :", c.Argent)
		fmt.Print("Choix : ")

		scanner.Scan()
		choix := scanner.Text()

		switch choix {
		case "1":
			acheterItemPayant(c, shopItems[0].Nom, shopItems[0].Prix)

		case "2":
			acheterItemPayant(c, shopItems[1].Nom, shopItems[1].Prix)

		case "3":
			acheterItemPayant(c, shopItems[2].Nom, shopItems[2].Prix)

		case "4":
			acheterItemPayant(c, shopItems[3].Nom, shopItems[3].Prix)

		case "5":
			acheterItemPayant(c, shopItems[4].Nom, shopItems[4].Prix)

		case "6":
			acheterItemPayant(c, shopItems[5].Nom, shopItems[5].Prix)

		case "0":
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func acheterItemPayant(c *structure.Character, nom string, prix int) {
	if c.Argent < prix {
		fmt.Println("Vous n'avez pas assez d'argent.")
		return
	}

	c.Argent -= prix

	for i := range c.Inventaire {
		if c.Inventaire[i].Nom == nom {
			c.Inventaire[i].Quantity++

			fmt.Println("Vous avez acheté :", nom)
			return
		}
	}

	c.AddInventory(structure.Item{
		Nom:      nom,
		Quantity: 1,
	})

	fmt.Println("Vous avez acheté :", nom)
}

func gainExperience(c *structure.Character, xpGagne int) {
	fmt.Println("Vous gagnez", xpGagne, "XP.")

	c.XPactu += xpGagne

	for c.XPactu >= c.XPmax {
		c.XPactu -= c.XPmax
		levelUp(c)
	}
}

func levelUp(c *structure.Character) {
	c.Niveau++

	bonusPV := 10

	c.PointsDeVieMaximum += bonusPV
	c.PointsDeVieActuels += bonusPV

	c.XPmax = int(float64(c.XPmax) * 1.2)

	fmt.Println()
	fmt.Println("===== NIVEAU SUPÉRIEUR =====")
	fmt.Println("Vous êtes maintenant niveau", c.Niveau)
	fmt.Println("Bonus : +", bonusPV, "PV maximum")
	fmt.Println("PV :", c.PointsDeVieActuels, "/", c.PointsDeVieMaximum)
}
