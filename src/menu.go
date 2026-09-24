package main

import (
    "bufio"
    "fmt"
    "main/structure"
    "os"
)

func menu(player *structure.Character) {

    for {

        fmt.Println()
        fmt.Println("*******************************")
        fmt.Println("             MENU")
        fmt.Println("*******************************")
        fmt.Println()
        fmt.Println("1 : Afficher les informations du personnage")
        fmt.Println("2 : Accéder à l'inventaire")
        fmt.Println("3 : Accéder au shop")
        fmt.Println("4 : Accéder à la forge")
        fmt.Println("5 : Combat")
        fmt.Println("6 : Quitter")
        fmt.Println()

        choix := lireEntier()

        switch choix {

        case 1:
            displayInfo(*player)
            pause()
            continue

        case 2:
            accessInventory(*player)
            pause()
            continue

        case 3:
            shop(player, bufio.NewReader(os.Stdin))
            pause()
            continue

        case 4:
            AccessForge(player)
            pause()
            continue

        case 5:
            trainingFightCommand(player)
            pause()
            continue

        case 6:
            fmt.Println("Fermeture du jeu...")
            return

        default:
            fmt.Println("Choix invalide.")
            pause()
            continue
        }
    }
}
