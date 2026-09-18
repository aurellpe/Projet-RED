package structure

import "fmt"

type Personnage struct {
	Nom                 string
	Prenom              string
	Age                 int
	Classe              string
	Niveau              int
	PointsDeVieMaximum  int
	PointsDeVieActuels  int
	Inventaire          []string
}
