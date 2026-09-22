package main

import (
    "fmt"
    "main/structure"
)

// =====================================================
// OBJETS DE LA FORGE
// =====================================================

type ForgeItem struct {
    Name   string
    Damage int
    Price  int
    Debris int
    Free   bool
}

var forgeItems = []ForgeItem{
    {"Épée en bois", 35, 0, 0, true},
    {"Épée en métal", 65, 5, 4, false},
    {"Hache", 50, 4, 2, false},
}

// =====================================================
// ÉTAT DE LA FORGE
// =====================================================

type ForgeState struct {
    Argent int
    Debris int
    Arme   string
}

var forgeStates = make(map[*structure.Character]*ForgeState)

func getForgeState(player *structure.Character) *ForgeState {
    if forgeStates[player] == nil {
        forgeStates[player] = &ForgeState{
            Argent: 10,
            Debris: 10,
            Arme:   "Épée en bois",
        }
    }
    return forgeStates[player]
}

// =====================================================
// DÉGÂTS DES ARMES
// =====================================================

func weaponDamage(name string) int {
    switch name {
    case "Épée en bois":
        return 35
    case "Épée en métal":
        return 65
    case "Hache":
        return 50
    default:
        return 5
    }
}

// =====================================================
// FORGE
// =====================================================

func AccessForge(player *structure.Character) {

    state := getForgeState(player)

    fmt.Println("\n================================")
    fmt.Println("             FORGE")
    fmt.Println("================================")
    fmt.Printf("Vous avez : %d or | %d débris\n\n", state.Argent, state.Debris)
    fmt.Printf("Arme équipée : %s (%d dégâts)\n\n", state.Arme, weaponDamage(state.Arme))

    fmt.Println("--------- ARMES ---------")
    fmt.Println("1 - Épée en bois (35 dégâts) - GRATUIT")
    fmt.Println("2 - Épée en métal (65 dégâts) - 4 débris + 5 or")
    fmt.Println("3 - Hache (50 dégâts) - 2 débris + 4 or")
    fmt.Println("0 - Retour\n")

    choice := lireEntier()

    switch choice {

    case 0:
        return

    case 1:
        if state.Arme == "Épée en bois" {
            fmt.Println("Vous avez déjà l'épée en bois équipée.")
            return
        }

        if !player.AddInventory(structure.Item{"Épée en bois", 1}) {
            fmt.Println("Inventaire plein (max 10).")
            return
        }

        state.Arme = "Épée en bois"
        fmt.Println("\nVous récupérez gratuitement une épée en bois !")
        fmt.Println("Ajoutée à l'inventaire.")
        return

    case 2:
        if state.Arme == "Épée en métal" {
            fmt.Println("Vous avez déjà l'épée en métal équipée.")
            return
        }

        if state.Argent < 5 {
            fmt.Printf("Pas assez d'or ! Il manque %d.\n", 5-state.Argent)
            return
        }

        if state.Debris < 4 {
            fmt.Printf("Pas assez de débris ! Il manque %d.\n", 4-state.Debris)
            return
        }

        if !player.AddInventory(structure.Item{"Épée en métal", 1}) {
            fmt.Println("Inventaire plein (max 10).")
            return
        }

        state.Argent -= 5
        state.Debris -= 4
        state.Arme = "Épée en métal"

        fmt.Println("\nVous fabriquez et équipez : Épée en métal !")
        fmt.Println("Dégâts : 65")
        fmt.Printf("Or restant : %d | Débris restants : %d\n", state.Argent, state.Debris)
        return

    case 3:
        if state.Arme == "Hache" {
            fmt.Println("Vous avez déjà la hache équipée.")
            return
        }

        if state.Argent < 4 {
            fmt.Printf("Pas assez d'or ! Il manque %d.\n", 4-state.Argent)
            return
        }

        if state.Debris < 2 {
            fmt.Printf("Pas assez de débris ! Il manque %d.\n", 2-state.Debris)
            return
        }

        if !player.AddInventory(structure.Item{"Hache", 1}) {
            fmt.Println("Inventaire plein (max 10).")
            return
        }

        state.Argent -= 4
        state.Debris -= 2
        state.Arme = "Hache"

        fmt.Println("\nVous fabriquez et équipez : Hache !")
        fmt.Println("Dégâts : 50")
        fmt.Printf("Or restant : %d | Débris restants : %d\n", state.Argent, state.Debris)
        return

    default:
        fmt.Println("Choix invalide.")
    }
}
