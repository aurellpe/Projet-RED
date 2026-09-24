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
	RecompenseDebris   int
}

func NewBoss() Boss {
	return Boss{
		Nom:                "Yael",
		Niveau:             67,
		PointsDeVieMaximum: 500,
		PointsDeVieActuels: 100,
		Degats:             15,
		Initiative:         15,
		RecompenseXP:       150,
		RecompenseArgent:   3,
		RecompenseDebris:   3,
	}
}
