package main

import (
	"fmt"
	"main/structure"
)


type ForgeItem struct {
	Name   string
	Damage int
	Price  int
	Debris int
	Free   bool
}

var forgeItems = []ForgeItem{
	{
		Name:   "Épée en bois",
		Damage: 35,
		Price:  0,
		Debris: 0,
		Free:   true,
	},
	{
		Name:   "Épée en métal",
		Damage: 65,
		Price:  5,
		Debris: 4,
		Free:   false,
	},
	{
		Name:   "Hache",
		Damage: 50,
		Price:  4,
		Debris: 2,
		Free:   false,
	},
}

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

func weaponDamage(name string) int {

	switch name {
	case "Épée en bois":
		return 35

	case "Épée en métal":
		return 65

	case "Hache":
		return 50
	}

	return 5
}

func AccessForge(player *structure.Character) {

	state := getForgeState(player)

	for {

		fmt.Println()
		fmt.Println("================================")
		fmt.Println("             FORGE")
		fmt.Println("================================")

		fmt.Printf(
			"Vous avez : %d pièces d'or | %d débris\n",
			state.Argent,
			state.Debris,
		)

		fmt.Println()

		fmt.Printf(
			"Arme équipée : %s (%d dégâts)\n",
			state.Arme,
			weaponDamage(state.Arme),
		)

		fmt.Println()
		fmt.Println("--------- ARMES ---------")

		fmt.Println("1 - Épée en bois (35 dégâts) - GRATUIT")

		fmt.Println("2 - Épée en métal (65 dégâts) - 4 débris + 5 or")

		fmt.Println("3 - Hache (50 dégâts) - 2 débris + 4 or")

		fmt.Println("0 - Retour")
		fmt.Println()

		choice := lireEntier()

		if choice == 0 {
			return
		}

		if choice < 1 || choice > len(forgeItems) {
			fmt.Println("Choix invalide.")
			continue
		}

		item := forgeItems[choice-1]

		if item.Free {

			if state.Arme == item.Name {
				fmt.Println("Vous avez déjà l'épée en bois équipée.")
				continue
			}

			state.Arme = item.Name

			fmt.Println()
			fmt.Println("Vous récupérez gratuitement une épée en bois !")
			fmt.Printf(
				"Arme équipée : %s (%d dégâts)\n",
				state.Arme,
				weaponDamage(state.Arme),
			)

			continue
		}

		if state.Argent < item.Price {

			fmt.Printf(
				"Pas assez d'or ! Il vous manque %d pièce(s).\n",
				item.Price-state.Argent,
			)

			continue
		}

		if state.Debris < item.Debris {

			fmt.Printf(
				"Pas assez de débris ! Il vous manque %d débris.\n",
				item.Debris-state.Debris,
			)

			continue
		}

		if state.Arme == item.Name {

			fmt.Printf(
				"Vous avez déjà équipé : %s.\n",
				item.Name,
			)

			continue
		}

		state.Argent -= item.Price
		state.Debris -= item.Debris

		state.Arme = item.Name

		fmt.Println()
		fmt.Printf(
			"Vous fabriquez et équipez : %s !\n",
			item.Name,
		)

		fmt.Printf(
			"Dégâts : %d\n",
			item.Damage,
		)

		fmt.Printf(
			"Or restant : %d\n",
			state.Argent,
		)

		fmt.Printf(
			"Débris restants : %d\n",
			state.Debris,
		)
	}
}