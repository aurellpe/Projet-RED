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

func initCharacter() Character {
    return Character{
        Nom:                string
        Prenom				string
		PointsDeVieActuels: int
        PointsDeVieMaximum: int
        Attaque:            int
        Or:                 int
        Exp:               	int
        ExpMax:             int
        Niveau:             int
        Fragments:          int
    }
}
