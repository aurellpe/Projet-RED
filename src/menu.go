package main

import (
	"fmt"
	"main/structure"
)

func menu(p *structure.Character) bool {
	fmt.Println("\n===== MENU PRINCIPAL =====")
	fmt.Println("1. Afficher les informations du personnage")
	fmt.Println("2. Accéder à l'inventaire")
	fmt.Println("3. Marchand")
	fmt.Println("4. Forgeron")
	fmt.Println("5. Combat")
	fmt.Println("6. Qui sont-ils ?")
	fmt.Println("0. Quitter")

	switch lireEntier() {
	case 1:
		displayInfo(*p)
	case 2:
		inventoryMenu(p, nil)
	case 3:
		shop(p, lecteur)
	case 4:
		AccessForge(p)
	case 5:
		trainingFightCommand(p)
	case 6:
		whoAreThey()
	case 0:
		fmt.Println(" Ciao!")
		return false
	default:
		fmt.Println("Choix invalide.")
	}
	return true
}
