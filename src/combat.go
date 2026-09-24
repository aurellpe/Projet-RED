package main

import ("fmt"
		"bufio"
		"os"
)

type Monster struct {
	nom                string
	pointsDeVieMax     int
	pointsDeVieActuels int
	pointsAttaque      int
}

var reader = bufio.NewReader(os.Stdin)


func initGoblin() Monster {
	return Monster{
		nom:                "Azrakar, Seigneur des Ombres",
		pointsDeVieMax:     1000,
		pointsDeVieActuels: 100,
		pointsAttaque:      25,
	}
}


func isDead(c *Character) bool {
	if c.pointsDeVieActuels <= 0 {
		fmt.Println(c.nom, "est mort...")
		c.pointsDeVieActuels = c.pointsDeVieMax / 2
		fmt.Printf("%s est ressuscité avec %d/%d points de vie.\n",
			c.nom, c.pointsDeVieActuels, c.pointsDeVieMax)
		return true
	}
	return false
}

func afficherPV(nom string, pvActuels int, pvMax int) {
	fmt.Printf("%s -> PV : %d / %d\n", nom, pvActuels, pvMax)
}

func goblinPattern(tour int, monstre *Monster, joueur *Character) {
	degats := monstre.pointsAttaque

	if tour%3 == 0 {
		degats = monstre.pointsAttaque * 2
	}

	joueur.pointsDeVieActuels -= degats

	fmt.Printf("%s inflige à %s %d de dégâts\n", monstre.nom, joueur.nom, degats)
	afficherPV(joueur.nom, joueur.pointsDeVieActuels, joueur.pointsDeVieMax)
}


func characterTurn(joueur *Character, monstre *Monster) {
	fmt.Println("\n--- À vous de jouer ---")
	fmt.Println("1. Attaquer")
	fmt.Println("2. Inventaire")
	fmt.Print("Votre choix : ")

	var choix string
	fmt.Fscan(reader, &choix)

	switch choix {
	case "1":
		degats := 5
		monstre.pointsDeVieActuels -= degats
		fmt.Printf("%s utilise Attaque basique\n", joueur.nom)
		fmt.Printf("%s inflige %d dégâts à %s\n", joueur.nom, degats, monstre.nom)
		afficherPV(monstre.nom, monstre.pointsDeVieActuels, monstre.pointsDeVieMax)

	case "2":
		if len(joueur.inventaire) == 0 {
			fmt.Println("Votre inventaire est vide.")
			return
		}

		fmt.Println("Votre inventaire :")
		for i, item := range joueur.inventaire {
			fmt.Printf("%d. %s\n", i+1, item)
		}
		fmt.Print("Choisissez un objet à utiliser : ")

		var indexChoisi int
		fmt.Fscan(reader, &indexChoisi)

		if indexChoisi < 1 || indexChoisi > len(joueur.inventaire) {
			fmt.Println("Choix invalide.")
			return
		}

		item := joueur.inventaire[indexChoisi-1]
		fmt.Println("Vous utilisez", item)

		if item == "Potion de vie" {
			joueur.pointsDeVieActuels += 50
			if joueur.pointsDeVieActuels > joueur.pointsDeVieMax {
				joueur.pointsDeVieActuels = joueur.pointsDeVieMax
			}
			afficherPV(joueur.nom, joueur.pointsDeVieActuels, joueur.pointsDeVieMax)
		}

		joueur.inventaire = append(joueur.inventaire[:indexChoisi-1], joueur.inventaire[indexChoisi:]...)

	default:
		fmt.Println("Choix invalide.")
	}
}

func trainingFight(player *structure.Character) {

    fmt.Println("\nUn gobelin apparaît !")

    goblinPV := 50
    goblinAtk := 10

    for {
        fmt.Println("\n=== Combat ===")
        fmt.Printf("Gobelin : %d PV\n", goblinPV)
        fmt.Printf("%s : %d/%d PV\n", player.Nom, player.PointsDeVieActuels, player.PointsDeVieMaximum)

        fmt.Println("1 - Attaquer")
        fmt.Println("2 - Utiliser une potion")
        fmt.Println("0 - Fuir")

        choix := lireEntier()

        if choix == 1 {
            goblinPV -= 15
            fmt.Println("Vous infligez 15 dégâts au gobelin.")
        }

        if choix == 2 {
            fmt.Println("Potion non implémentée.")
        }

        if choix == 0 {
            fmt.Println("Vous fuyez le combat.")
            return
        }

        if goblinPV <= 0 {
            fmt.Println("Vous avez vaincu le gobelin !")
            return
        }

        player.PointsDeVieActuels -= goblinAtk
        fmt.Printf("Le gobelin vous inflige %d dégâts.\n", goblinAtk)

        if player.PointsDeVieActuels <= 0 {
            fmt.Println("Vous êtes mort... Réanimation à 50% PV.")
            player.PointsDeVieActuels = player.PointsDeVieMaximum / 2
            return
        }
    }
}