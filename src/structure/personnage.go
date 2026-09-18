package structure

type Item struct {
	Nom string
	Quantity int 
}

type Character struct {
	Nom                 string
	Prenom              string
	Age                 int
	Classe              string
	Niveau              int
	PointsDeVieMaximum  int
	PointsDeVieActuels  int
	Inventaire          []Item
}
