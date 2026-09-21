package main

import (
	"fmt"
	"main/structure"
	"strconv"
)

// lireEntier lit une saisie et la convertit en nombre (-1 si ce n'est pas un nombre)
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

// retour affiche "0 - Retour" et attend que le joueur tape 0
func retour() {
	for {
		fmt.Println("\n0 - Retour")
		if lireEntier() == 0 {
			return
		}
	}
}

func menu(player *structure.Character) {
	for {
		fmt.Println("\n===== MENU =====")
		fmt.Println("1 - Afficher les informations du personnage")
		fmt.Println("2 - Accéder à l'inventaire")
		fmt.Println("3 - Accéder à la forge")
		fmt.Println("4 - Quitter")

		switch lireEntier() {
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
