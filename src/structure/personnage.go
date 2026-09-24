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
	Degats 				int
	Debris				int
	Arme				string
}

func (p *Character) AddInventory(item Item) bool {

	for i := range p.Inventaire {
		if p.Inventaire[i].Nom == item.Nom {
			p.Inventaire[i].Quantity += item.Quantity
			return true
		}
	}

	if len(p.Inventaire) >= 10 {
		return false
	}

	p.Inventaire = append(p.Inventaire, item)

	return true
}