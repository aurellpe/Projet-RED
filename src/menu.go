package main

import (
	"fmt"
	"main/structure"
)

func menu(player *structure.Character) {
	for {
		fmt.Println("\n===== MENU =====")
		fmt.Println("1 - Afficher les informations du personnage")
		fmt.Println("2 - Accéder à l'inventaire")
		fmt.Println("3 - Quitter")
		fmt.Print("Votre choix : ")

		var choice int
		fmt.Scan(&choice)

		switch choice {

		case 1:
			displayInfo(*player)

			fmt.Println("\n0 - Retour")

			var back int
			fmt.Scan(&back)

			if back == 0 {
				continue
			}

		case 2:
			accessInventory(*player)

			fmt.Println("\n0 - Retour")

			var back int
			fmt.Scan(&back)

			if back == 0 {
				continue
			}

		case 3:
			fmt.Println("Fermeture du jeu...")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}
