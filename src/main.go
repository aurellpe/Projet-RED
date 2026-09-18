package main

import (
	"fmt"
	"main/structure"
)

func initCharacter(nom string, prenom string, age int, classe string, niveau int, pointsDeVieMaximum int, pointsDeVieActuels int, inventaire []string) structure.Personnage {
	return structure.Personnage{
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

func displayInfo(perso structure.Personnage) {
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", perso.Nom)
	fmt.Println("Prénom :", perso.Prenom)
	fmt.Println("Âge :", perso.Age)
	fmt.Println("Classe :", perso.Classe)
	fmt.Println("Niveau :", perso.Niveau)
	fmt.Println("Points de vie :", perso.PointsDeVieActuels, "/", perso.PointsDeVieMaximum)
}

func main() {
	perso := initCharacter(
		"DE Bergerac",
		"Dartagnan",
		26,
		"Legionnaire",
		1,
		10000,
		150,
		inventaire,
	)

	displayInfo(perso)
}
