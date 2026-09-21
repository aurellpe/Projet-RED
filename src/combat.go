package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Personnage struct {
    Nom                 string
    PointsDeVieActuels  int
    PointsDeVieMaximum  int
    Attaque             int
    Or                  int
    Exp                 int
    ExpMax              int
    Niveau              int
    Fragments           int
}

type Monstre struct {
    Nom                 string
    PointsDeVieActuels  int
    PointsDeVieMaximum  int
    Attaque             int
    ExpGain             int
}

func InitPersonnage() Personnage {
    return Personnage{
        Nom:                "Dartagnan",
        PointsDeVieActuels: 90,
        PointsDeVieMaximum: 10000,
        Attaque:            35,
        Or:                 0,
        Exp:                0,
        ExpMax:             1000000,
        Niveau:             1,
        Fragments:          0,
    }
}

func InitMonstre() Monstre {
    return Monstre{
        Nom:                "Python",
        PointsDeVieActuels: 100,
        PointsDeVieMaximum: 10000,
        Attaque:            25,
        ExpGain:            50,
    }
}

func GenererRecompenses() (int, int) {
    rand.Seed(time.Now().UnixNano())

    or := rand.Intn(7-3+1) + 3
    fragments := rand.Intn(10) + 1

    for fragments == or {
        fragments = rand.Intn(10) + 1
    }

    return or, fragments
}

func TourPersonnage(p *Personnage, m *Monstre) {
    fmt.Println("\n=== Tour du joueur ===")
    fmt.Println("1 - Attaquer")
    fmt.Println("2 - Passer le tour")

    var choix int
    fmt.Scanln(&choix)

    if choix == 1 {
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

func TourMonstre(p *Personnage, m *Monstre) {
    if m.PointsDeVieActuels <= 0 {
        return
    }

    fmt.Printf("%s attaque et inflige %d dégâts à %s !\n", m.Nom, m.Attaque, p.Nom)
    p.PointsDeVieActuels -= m.Attaque
    if p.PointsDeVieActuels < 0 {
        p.PointsDeVieActuels = 0
    }
    fmt.Printf("PV du joueur : %d / %d\n", p.PointsDeVieActuels, p.PointsDeVieMaximum)
}

func DonnerRecompenses(p *Personnage, m Monstre) {
    or, fragments := GenererRecompenses()

    fmt.Println("\n🎉 Victoire !")
    fmt.Printf("Vous gagnez %d XP, %d pièces d'or et %d fragments.\n",
        m.ExpGain, or, fragments)

    p.Exp += m.ExpGain
    p.Or += or
    p.Fragments += fragments

    VerifierNiveau(p)
}

func BonusFinDeRound(p *Personnage) {
    fmt.Println("\n✨ Bonus de fin de round : +50 XP")
    p.Exp += 50
    VerifierNiveau(p)
}

func VerifierNiveau(p *Personnage) {
    for p.Exp >= p.ExpMax {
        p.Exp -= p.ExpMax
        p.Niveau++
        p.ExpMax += 50
        p.PointsDeVieMaximum += 20
        p.PointsDeVieActuels = p.PointsDeVieMaximum

        fmt.Printf("⬆️ Niveau %d atteint ! PV max : %d\n", p.Niveau, p.PointsDeVieMaximum)
    }
}

func CombatEntrainement(p *Personnage) {
    round := 1

    for p.PointsDeVieActuels > 0 {
        fmt.Printf("\n=== ROUND %d ===\n", round)
        monstre := InitMonstre()

        for p.PointsDeVieActuels > 0 && monstre.PointsDeVieActuels > 0 {
            TourPersonnage(p, &monstre)
            if monstre.PointsDeVieActuels <= 0 {
                DonnerRecompenses(p, monstre)
                break
            }
            TourMonstre(p, &monstre)
        }

        if p.PointsDeVieActuels <= 0 {
            fmt.Println("💀 Vous êtes mort... Réanimation à 50% PV.")
            p.PointsDeVieActuels = p.PointsDeVieMaximum / 2
        }

        BonusFinDeRound(p)
        round++

        fmt.Println("\nContinuer ? (1 = Oui / 0 = Non)")
        var cont int
        fmt.Scanln(&cont)
        if cont == 0 {
            break
        }
    }

    fmt.Println("\nFin de l'entraînement.")
    fmt.Printf("Stats finales : Niveau %d | XP : %d/%d | Or : %d | Fragments : %d\n",
        p.Niveau, p.Exp, p.ExpMax, p.Or, p.Fragments)
}
