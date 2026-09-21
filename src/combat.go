package main

import (
    "fmt"
    "math/rand"
    "time"
    "main/structure"
)


func GenerateRewards() (int, int) {
    rand.Seed(time.Now().UnixNano())

    gold := rand.Intn(7-3+1) + 3
    fragments := rand.Intn(10) + 1

    for fragments == gold {
        fragments = rand.Intn(10) + 1
    }

    return gold, fragments
}

func PlayerTurn(p *structure.Character, m *structure.Monster) {
    fmt.Println("\n=== Tour du joueur ===")
    fmt.Println("1 - Attaquer")
    fmt.Println("2 - Passer le tour")

    var choice int
    fmt.Scanln(&choice)

    if choice == 1 {
        fmt.Printf("%s attaque et inflige %d dégâts à %s !\n", p.Nom, p.Attaque, m.Nom)
        m.PointsDeVieActuels -= p.Attaque
        if m.PointsDeVieActuels < 0 {
            m.PointsDeVieActuels = 0
        }
        fmt.Printf("PV du monstre : %d / %d\n", m.PointsDeVieActuels, m.PointsDeVieMaximum)
    } else {
        fmt.Println("Vous observez le monstre...")
    }
}

func MonsterTurn(p *structure.Character, m *structure.Monster) {
    if m.PointsDeVieActuels <= 0 {
        return
    }

    fmt.Printf("%s attaque et inflige %d dégâts à %s !\n", m.Nom, m.PointsDattaque, p.Nom)
    p.PointsDeVieActuels -= m.PointsDattaque
    if p.PointsDeVieActuels < 0 {
        p.PointsDeVieActuels = 0
    }
    fmt.Printf("PV du joueur : %d / %d\n", p.PointsDeVieActuels, p.PointsDeVieMaximum)
}

func GiveRewards(p *structure.Character, m structure.Monster) {
    gold, fragments := GenerateRewards()

    fmt.Println("\n🎉 Victoire !")
    fmt.Printf("Vous gagnez %d XP, %d pièces d'or et %d fragments.\n",
        m.ExpGain, gold, fragments)

    p.Exp += m.ExpGain
    p.Or += gold
    p.Fragments += fragments

    CheckLevel(p)
}

func EndOfRoundBonus(p *structure.Character) {
    fmt.Println("\n✨ Bonus de fin de round : +50 XP")
    p.Exp += 50
    CheckLevel(p)
}

func CheckLevel(p *structure.Character) {
    for p.Exp >= p.ExpMax {
        p.Exp -= p.ExpMax
        p.Niveau++
        p.ExpMax += 50
        p.PointsDeVieMaximum += 20
        p.PointsDeVieActuels = p.PointsDeVieMaximum

        fmt.Printf("⬆️ Niveau %d atteint ! PV max : %d\n", p.Niveau, p.PointsDeVieMaximum)
    }
}

func TrainingFight(p *structure.Character) {
    round := 1

    for p.PointsDeVieActuels > 0 {
        fmt.Printf("\n=== ROUND %d ===\n", round)
        monster := structure.InitMonster()

        for p.PointsDeVieActuels > 0 && monster.PointsDeVieActuels > 0 {
            PlayerTurn(p, &monster)
            if monster.PointsDeVieActuels <= 0 {
                GiveRewards(p, monster)
                break
            }
            MonsterTurn(p, &monster)
        }

        if p.PointsDeVieActuels <= 0 {
            fmt.Println("💀 Vous êtes mort... Réanimation à 50% PV.")
            p.PointsDeVieActuels = p.PointsDeVieMaximum / 2
        }

        EndOfRoundBonus(p)
        round++

        fmt.Println("\nContinuer ? (1 = Oui / 0 = Non)")
        var keepGoing int
        fmt.Scanln(&keepGoing)
        if keepGoing == 0 {
            break
        }
    }

    fmt.Println("\nFin de l'entraînement.")
    fmt.Printf("Stats finales : Niveau %d | XP : %d/%d | Or : %d | Fragments : %d\n",
        p.Niveau, p.Exp, p.ExpMax, p.Or, p.Fragments)
}
