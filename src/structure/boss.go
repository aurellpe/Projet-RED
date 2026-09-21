package structure

type Boss struct {
    Nom                string
    Niveau             int
    PointsDeVieMaximum int
    PointsDeVieActuels int
    Degats             int
    Initiative         int
    RecompenseXP       int
    RecompenseArgent   int
}

func NewBoss() Boss {
    return Boss{
        Nom:                "Azrakar, Seigneur des Ombres",
        Niveau:             5,
        PointsDeVieMaximum: 500,
        PointsDeVieActuels: 500,
        Degats:             35,
        Initiative:         15,
        RecompenseXP:       200,
        RecompenseArgent:   100,
    }
}
package structure

type Boss struct {
	Nom                string
	Niveau             int
	PointsDeVieMaximum int
	PointsDeVieActuels int
	Degats             int
	Initiative         int
	RecompenseXP       int
	RecompenseArgent   int
}

func NewBoss() Boss {
	return Boss{
		Nom:                "Azrakar, Seigneur des Ombres",
		Niveau:             5,
		PointsDeVieMaximum: 500,
		PointsDeVieActuels: 500,
		Degats:             35,
		Initiative:         15,
		RecompenseXP:       200,
		RecompenseArgent:   100,
	}
}
