package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"os"
	"strconv"
	"strings"
)

func initCharacter(nom string, prenom string, age int, classe string, niveau int, pointsDeVieMaximum int, pointsDeVieActuels int, inventaire []structure.Item, argent int) structure.Character {
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

	// ===== BOSS =====
	boss := initBoss()

	var run bool = true

	for run {
		fmt.Println()
		fmt.Println("===== MENU =====")
		fmt.Println("1. Afficher personnage")
		fmt.Println("2. Inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Combattre le boss")

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

		// ===== BOSS =====
		if choix == "4" {
			combatBoss(&perso, &boss, scanner)
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

// ==================================================
// ====================== BOSS =======================
// ==================================================

type Boss struct {
	Nom                string
	Niveau             int
	PointsDeVieMaximum int
	PointsDeVieActuels int
	Degats             int
	Initiative         int
	RecompenseXP       int
	RecompenseArgent   int
}

func initBoss() Boss {
	return Boss{
		Nom:                "Azrakar, Seigneur des Ombres",
		Niveau:             5,
		PointsDeVieMaximum: 500,
		PointsDeVieActuels: 500,
		Degats:             35,
		Initiative:         15,
		RecompenseXP:       200,
		RecompenseArgent:   100,
	}
}

func displayBoss(boss Boss) {
	fmt.Println()
	fmt.Println("========== BOSS ==========")
	fmt.Println("Nom :", boss.Nom)
	fmt.Println("Niveau :", boss.Niveau)
	fmt.Println("PV :", boss.PointsDeVieActuels, "/", boss.PointsDeVieMaximum)
	fmt.Println("Dégâts :", boss.Degats)
	fmt.Println("==========================")
}

func combatBoss(perso *structure.Character, boss *Boss, scanner *bufio.Scanner) {
	if boss.PointsDeVieActuels <= 0 {
		fmt.Println("Vous avez déjà vaincu", boss.Nom)
		return
	}

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("       UN BOSS APPARAÎT !")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println(boss.Nom)
	fmt.Println("Niveau :", boss.Niveau)
	fmt.Println("PV :", boss.PointsDeVieActuels)
	fmt.Println()

	for perso.PointsDeVieActuels > 0 && boss.PointsDeVieActuels > 0 {
		fmt.Println()
		fmt.Println("---------- COMBAT ----------")
		fmt.Println(perso.Prenom)
		fmt.Println("PV :", perso.PointsDeVieActuels, "/", perso.PointsDeVieMaximum)
		fmt.Println()
		fmt.Println(boss.Nom)
		fmt.Println("PV :", boss.PointsDeVieActuels, "/", boss.PointsDeVieMaximum)
		fmt.Println()
		fmt.Println("1. Attaquer")
		fmt.Println("2. Informations du boss")
		fmt.Println("0. Fuir")
		fmt.Print("Choix : ")

		scanner.Scan()
		choix := scanner.Text()

		switch choix {

		case "1":
			degatsJoueur := 50 + perso.Niveau*5

			boss.PointsDeVieActuels -= degatsJoueur

			if boss.PointsDeVieActuels < 0 {
				boss.PointsDeVieActuels = 0
			}

			fmt.Println()
			fmt.Println("Vous attaquez", boss.Nom)
			fmt.Println("Vous infligez", degatsJoueur, "dégâts.")

			if boss.PointsDeVieActuels <= 0 {
				fmt.Println()
				fmt.Println("======================================")
				fmt.Println("           BOSS VAINCU !")
				fmt.Println("======================================")
				fmt.Println()
				fmt.Println("Vous avez vaincu :", boss.Nom)
				fmt.Println("XP gagné :", boss.RecompenseXP)
				fmt.Println("Or gagné :", boss.RecompenseArgent)

				perso.XPactu += boss.RecompenseXP
				perso.Argent += boss.RecompenseArgent

				return
			}

			fmt.Println()
			fmt.Println(boss.Nom, "vous attaque !")

			perso.PointsDeVieActuels -= boss.Degats

			if perso.PointsDeVieActuels < 0 {
				perso.PointsDeVieActuels = 0
			}

			fmt.Println("Le boss vous inflige", boss.Degats, "dégâts.")

			if perso.PointsDeVieActuels <= 0 {
				fmt.Println()
				fmt.Println("======================================")
				fmt.Println("          VOUS ÊTES MORT")
				fmt.Println("======================================")
				return
			}

		case "2":
			displayBoss(*boss)

		case "0":
			fmt.Println("Vous fuyez le combat.")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}