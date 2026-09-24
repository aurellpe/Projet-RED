package main

import (
    "fmt"
    "main/structure"
)

func AccessForge(player *structure.Character) {
    for {
        fmt.Println("===============================")
        fmt.Println("            FORGE")
        fmt.Println("===============================")
        fmt.Println()

        fmt.Printf("Vous avez : %d pièces d'or | %d débris\n", player.Argent, player.Debris)
        fmt.Printf("Arme équipée : %s (%d dégâts)\n\n", player.Arme, player.Degats)

        fmt.Println("--------- ARMES ---------")
        fmt.Println("1 - Épée en bois (35 dégâts) - GRATUIT")
        fmt.Println("2 - Épée en métal (65 dégâts) - 4 débris + 5 or")
        fmt.Println("3 - Hache (50 dégâts) - 2 débris + 4 or")
        fmt.Println("0 - Retour")
        fmt.Println()

        choix := lireEntier()

        switch choix {
        case 1:
            player.Arme = "Épée en bois"
            player.Degats = 35
            fmt.Println("Vous équipez l'Épée en bois.")
            pause()

        case 2:
            if player.Debris >= 4 && player.Argent >= 5 {
                player.Arme = "Épée en métal"
                player.Degats = 65
                player.Debris -= 4
                player.Argent -= 5
                fmt.Println("Vous équipez l'Épée en métal.")
            } else {
                fmt.Println("Pas assez de ressources.")
            }
            pause()

        case 3:
            if player.Debris >= 2 && player.Argent >= 4 {
                player.Arme = "Hache"
                player.Degats = 50
                player.Debris -= 2
                player.Argent -= 4
                fmt.Println("Vous équipez la Hache.")
            } else {
                fmt.Println("Pas assez de ressources.")
            }
            pause()

        case 0:
            return

        default:
            fmt.Println("Choix invalide.")
            pause()
        }
    }
}
