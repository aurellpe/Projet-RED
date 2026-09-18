package main

import "main/structure"

func main() {
	perso := structure.Personnage{
		Nom:                "DE Bergerac",
		Prenom:             "Dartagnan",
		Age:                26,
		Classe:             "Legionnaire",
		Niveau:             1,
		PointsDeVieMaximum: 10000,
		PointsDeVieActuels: 150,
		Inventaire:         []string{},
	}

}
