package structure

type Item struct {
	Nom      string
	Quantity int
}

type Character struct {
	Nom                string
	Prenom             string
	Age                int
	Classe             string
	Niveau             int
	PointsDeVieMaximum int
	PointsDeVieActuels int
	Attaque            int
	Or                 int
	Exp                int
	ExpMax             int
	Fragments          int
	Inventaire         []Item
}

func initCharacter() Character {
	return Character{
		Nom:                "De bergerac",
		Prenom:             "Dartagnan",
		Age:                26,
		Classe:             "Légionnaire",
		Niveau:             1,
		PointsDeVieMaximum: 10000,
		PointsDeVieActuels: 90,
		Attaque:            35,
		Or:                 0,
		Exp:                0,
		ExpMax:             100000,
		Fragments:          0,
		Inventaire:         []Item{},
	}
}
