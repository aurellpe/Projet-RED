package main

import (
	"fmt"
	"main/structure"
	"strconv"
)

func lireEntier() int {
	var saisie string

	fmt.Print("Votre choix : ")
	fmt.Scan(&saisie)

	n, err := strconv.Atoi(saisie)

	if err != nil {
		return -1
	}

	return n
}

func retour() {
	for {
		fmt.Println("\n0 - Retour")

		if lireEntier() == 0 {
			return
		}

		fmt.Println("Veuillez taper 0 pour revenir.")
	}
}


func menu(player *structure.Character) {

	for {

		fmt.Println()
		fmt.Println("==============================")
		fmt.Println("           MENU")
		fmt.Println("==============================")

		fmt.Println("1 - Afficher les informations du personnage")
		fmt.Println("2 - Accéder à l'inventaire")
		fmt.Println("3 - Accéder à la forge")
		fmt.Println("4 - Quitter")

		fmt.Println()

		choix := lireEntier()

		switch choix {

		case 1:

			displayInfo(*player)

			retour()

		case 2:

			accessInventory(*player)

			retour()

		case 3:

			accessForge(player)

			
		case 4:
			fmt.Println("Fermeture du jeu...")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}
