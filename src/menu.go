package main

import (
	"fmt"
	"main/structure"
)

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
            return

        case 2:
            accessInventory(*player)
            return

        case 3:
            AccessForge(player)

        case 4:
            fmt.Println("Fermeture du jeu...")
            return

        default:
            fmt.Println("Choix invalide.")
        }
    }
}
