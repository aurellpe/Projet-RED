package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"time"
)

func isDead(c *structure.Character) {
	if c.PointsDeVieActuels <= 0 {
		fmt.Println()
		fmt.Println("==========================")
		fmt.Println("       VOUS ÊTES MORT")
		fmt.Println("==========================")

		time.Sleep(2 * time.Second)

		c.PointsDeVieActuels = c.PointsDeVieMaximum / 2

		fmt.Println()
		fmt.Println("Vous êtes ressuscité !")
		fmt.Println("PV :", c.PointsDeVieActuels, "/", c.PointsDeVieMaximum)
	}
}

func displayBoss(boss structure.Boss) {
	fmt.Println()
	fmt.Println("========== BOSS ==========")
	fmt.Println("Nom :", boss.Nom)
	fmt.Println("Niveau :", boss.Niveau)
	fmt.Println("PV :", boss.PointsDeVieActuels, "/", boss.PointsDeVieMaximum)
	fmt.Println("Dégâts :", boss.Degats)
	fmt.Println("==========================")
}

func combatBoss(
	perso *structure.Character,
	boss *structure.Boss,
	scanner *bufio.Scanner,
) {
	if boss.PointsDeVieActuels <= 0 {
		fmt.Println("Vous avez déjà vaincu", boss.Nom)
		return
	}

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("       UN BOSS APPARAÎT")
	fmt.Println("================================")
	fmt.Println()
	fmt.Println(boss.Nom)
	fmt.Println("Niveau :", boss.Niveau)
	fmt.Println("PV :", boss.PointsDeVieActuels, "/", boss.PointsDeVieMaximum)

	for boss.PointsDeVieActuels > 0 {
		fmt.Println()
		fmt.Println("========== COMBAT ==========")

		fmt.Println(
			perso.Prenom,
			":",
			perso.PointsDeVieActuels,
			"/",
			perso.PointsDeVieMaximum,
			"PV",
		)

		fmt.Println(
			boss.Nom,
			":",
			boss.PointsDeVieActuels,
			"/",
			boss.PointsDeVieMaximum,
			"PV",
		)

		fmt.Println()
		fmt.Println("1. Attaquer")
		fmt.Println("2. Utiliser une potion")
		fmt.Println("3. Informations du boss")
		fmt.Println("0. Fuir")
		fmt.Print("Choix : ")

		scanner.Scan()
		choix := scanner.Text()

		switch choix {
		case "1":
			tourJoueur(perso, boss)

			if boss.PointsDeVieActuels <= 0 {
				fmt.Println()
				fmt.Println("==========================")
				fmt.Println("       BOSS VAINCU")
				fmt.Println("==========================")

				fmt.Println("Vous avez vaincu", boss.Nom)

				fmt.Println(
					"Vous gagnez",
					boss.RecompenseArgent,
					"pièces d'or.",
				)

				perso.Argent += boss.RecompenseArgent

				gainExperience(perso, boss.RecompenseXP)

				return
			}

			fmt.Println()
			fmt.Println("Le boss prépare son attaque...")

			time.Sleep(2 * time.Second)

			fmt.Println()
			fmt.Println("===== TOUR DU BOSS =====")

			tourBoss(perso, boss)

			if perso.PointsDeVieActuels <= 0 {
				isDead(perso)
				fmt.Println()
				fmt.Println("Le combat est terminé.")
				return
			}

		case "2":
			takePot(perso)

			fmt.Println()
			fmt.Println("Le boss prépare son attaque...")

			time.Sleep(2 * time.Second)

			fmt.Println()
			fmt.Println("===== TOUR DU BOSS =====")

			tourBoss(perso, boss)

			if perso.PointsDeVieActuels <= 0 {
				isDead(perso)
				return
			}

		case "3":
			displayBoss(*boss)

		case "0":
			fmt.Println("Vous fuyez le combat.")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func tourJoueur(
	perso *structure.Character,
	boss *structure.Boss,
) {
	degatsJoueur := 50 + perso.Niveau*5

	fmt.Println()
	fmt.Println("===== VOTRE TOUR =====")
	fmt.Println()
	fmt.Println("Vous attaquez", boss.Nom, "!")

	boss.PointsDeVieActuels -= degatsJoueur

	if boss.PointsDeVieActuels < 0 {
		boss.PointsDeVieActuels = 0
	}

	fmt.Println(
		"Vous infligez",
		degatsJoueur,
		"dégâts.",
	)

	fmt.Println(
		"PV du boss :",
		boss.PointsDeVieActuels,
		"/",
		boss.PointsDeVieMaximum,
	)
}

func tourBoss(
	perso *structure.Character,
	boss *structure.Boss,
) {
	fmt.Println()
	fmt.Println(boss.Nom, "vous attaque !")

	perso.PointsDeVieActuels -= boss.Degats

	if perso.PointsDeVieActuels < 0 {
		perso.PointsDeVieActuels = 0
	}

	fmt.Println(
		"Le boss vous inflige",
		boss.Degats,
		"dégâts.",
	)

	fmt.Println(
		"Vos PV :",
		perso.PointsDeVieActuels,
		"/",
		perso.PointsDeVieMaximum,
	)
}
