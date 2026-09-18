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

func accessInventory(perso structure.Personnage) {
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

	inventaire := []string{
	}

	perso := initCharacter(
		"DE Bergerac",
		"Dartagnan",
		26,
		"Elfe",
		1,
		100,
		40,
		inventaire,
	)

	displayInfo(perso)

	fmt.Println()

	accessInventory(perso)
}