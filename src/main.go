package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"os"
	"strconv"
	"strings"
)

func initCharacter(nom string, prenom string, age int, classe string, niveau int, pointsDeVieMaximum int, pointsDeVieActuels int, inventaire []structure.Item, argent int ) structure.Character {
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
	}
}

func displayInfo(perso structure.Character) {
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", perso.Nom)
	fmt.Println("Prénom :", perso.Prenom)
	fmt.Println("Âge :", perso.Age)
	fmt.Println("Classe :", perso.Classe)
	fmt.Println("Niveau :", perso.Niveau)
	fmt.Println("Points de vie :", perso.PointsDeVieActuels, "/", perso.PointsDeVieMaximum)
}

func accessInventory(perso structure.Character) {
	fmt.Println("===== INVENTAIRE =====")

	if len(perso.Inventaire) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	for i, objet := range perso.Inventaire {
		fmt.Println(i+1, "-", objet)
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	perso := initCharacter(
		"DE Bergerac",
		"Dartagnan",
		26,
		"Légionnaire",
		1,
		10000,
		150,
		[]structure.Item{{Nom: "potion", Quantity: 0}, {Nom: "epee", Quantity: 1}},
		100,
	)

	var run bool = true 
	for run{
		scanner.Scan()
		choix := scanner.Text()
		if choix == "1" {
			displayInfo(perso)
			fmt.Println()
		}
		if choix == "2" {
			accessInventory(perso)
			fmt.Println()
		}
		if choix == "3" {
			shopMenu(&perso, scanner)
			fmt.Println()
}
	}


	shopMenu(&perso, scanner)
}


func takePot(p structure.Character) {
	for i, it := range p.Inventaire {
		if it.Nom == "potion" && it.Quantity > 0 {
			p.PointsDeVieActuels += 50
			if p.PointsDeVieActuels > p.PointsDeVieMaximum {
				p.PointsDeVieActuels = p.PointsDeVieMaximum
			}
			p.Inventaire[i].Quantity--

			if p.Inventaire[i].Quantity == 0 {
				p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
			}
			return

		}
	}
}

func shopMenu(c *structure.Character, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n----- Marchand -----")
		fmt.Println("1. Potion de vie")
		fmt.Println("0. Retour")

		scanner.Scan()
		choix := scanner.Text()

		switch choix {
		case "1":
			acheterItem(c, "Potion de vie")
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
	{Nom: "Potion de vie", Prix: 3},
	{Nom: "Potion de poison", Prix: 6},
	{Nom: "Livre de Sort : Boule de Feu", Prix: 25},
	{Nom: "Fourrure de Loup", Prix: 4},
	{Nom: "Peau de Troll", Prix: 7},
	{Nom: "Cuir de vache", Prix: 1},
}

func shop(c *structure.Character, reader *bufio.Reader) {
	for {
		fmt.Println("\n----- Marchand -----")
		for i, si := range shopItems {
			fmt.Printf("%d. %s (%d pièces d'or)\n", i+1, si.Nom, si.Prix)
		}
		fmt.Println("0. Retour")
		fmt.Println("Votre argent :", c.Argent, "pièces d'or")

		choix := lireChoix(reader)

		if choix == 0 {
			return
		} else if choix >= 1 && choix <= len(shopItems) {
			si := shopItems[choix-1]
			acheterItemPayant(c, si.Nom, si.Prix)
		} else {
			fmt.Println("Choix invalide.")
		}
	}
}

func acheterItemPayant(c *structure.Character, nom string, prix int) {
	if c.Argent < prix {
		fmt.Println("Vous n'avez pas assez d'argent pour acheter :", nom)
		return
	}

	c.Argent -= prix
	c.AddInventory(structure.Item{Nom: nom, Quantity: 1})
	fmt.Println("Vous avez acheté :", nom)

}

func lireChoix(reader *bufio.Reader) int {
	fmt.Print("Choix : ")
	ligne, _ := reader.ReadString('\n')
	ligne = strings.TrimSpace(ligne)
	choix, err := strconv.Atoi(ligne)
	if err != nil {
		return -1
	}
	return choix
}

func gainExperience(c *structure.Character, xpGagne int) {
	fmt.Printf("Vous gagnez %d points d'expérience.\n", xpGagne)
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

	c.XPmax = int(float64(c.XPactu) * 1.2)

	fmt.Printf("Niveau supérieur ! Vous êtes maintenant niveau %d.\n", c.Niveau)
	fmt.Printf("Bonus : +%d points de vie maximum.\n", bonusPV)
	fmt.Printf("PV : %d / %d\n", c.PointsDeVieActuels, c.PointsDeVieMaximum)
}