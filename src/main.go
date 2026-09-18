package main

import (
	"fmt"
	"main/structure"
)



func initCharacter(nom string, prenom string, age int, classe string, niveau int, pointsDeVieMaximum int, pointsDeVieActuels int, inventaire []structure.Item) structure.Character {
	return structure.Character{
		Nom:                 nom,
		Prenom:              prenom,
		Age:                 age,
		Classe:              classe,
		Niveau:              niveau,
		PointsDeVieMaximum:  pointsDeVieMaximum,
		PointsDeVieActuels:  pointsDeVieActuels,
		Inventaire:          inventaire,
	}
}

func displayInfo(perso structure.Character) {
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", perso.Nom)
	fmt.Println("Prénom :", perso.Prenom)
	fmt.Println("Âge :", perso.Age)
	fmt.Println("Classe :", perso.Classe)
	fmt.Println("Niveau :", perso.Niveau)
	fmt.Println("Points de vie :", perso.PointsDeVieActuels, "/", perso.PointsDeVieMaximum)
}

func accessInventory(perso structure.Character) {
	fmt.Println("===== INVENTAIRE =====")

	if len(perso.Inventaire) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	for i, objet := range perso.Inventaire {
		fmt.Println(i+1, "-", objet)
	}
}

func main() {

	

	perso := initCharacter(
		"DE Bergerac",
		"Dartagnan",
		26,
		"Légionnaire",
		1,
		10000,
		150,
		[]structure.Item{{Nom : "potion", Quantity : 0}, {Nom : "epee", Quantity : 1}},
	)

	displayInfo(perso)

	fmt.Println()

	accessInventory(perso)
}

	

func takePot( p structure.Character){
	for i, it := range p.Inventaire {
		if it.Nom == "potion" && it.Quantity > 0 {
			p.PointsDeVieActuels +=50
			if p.PointsDeVieActuels > p.PointsDeVieMaximum {
				p.PointsDeVieActuels = p.PointsDeVieMaximum
			}
			p.Inventaire[i].Quantity--
			
			if p.Inventaire[i].Quantity == 0 {
				p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
			}
			return

	

		}
	}
}