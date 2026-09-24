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
	Inventaire         []Item
	InventaireMax      int
	Ameliorations      int
	Skills             []string
	Argent             int
	XPmax              int
	XPactu             int
	Mana               int
	ManaMax            int
	Initiative         int
	Equipement         Equipment
	Degats             int
	Debris             int
	Arme               string
}

type Equipment struct {
	Tete  string
	Torse string
	Pieds string
}

func (p *Character) AddInventory(item Item) bool {

	for i := range p.Inventaire {
		if p.Inventaire[i].Nom == item.Nom {
			p.Inventaire[i].Quantity += item.Quantity
			return true
		}
	}

	max := p.InventaireMax
	if max <= 0 {
		max = 10
	}
	if len(p.Inventaire) >= max {
		return false
	}

	p.Inventaire = append(p.Inventaire, item)

	return true
}
