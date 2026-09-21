package structure

type Item struct {
	Nom      string
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
    XPmax               int
	XPactu              int
	Argent              int
}

func(p *Character)AddInventory(item Item) {
	p.Inventaire = append(p.Inventaire,item)
}