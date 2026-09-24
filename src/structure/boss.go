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
        Nom:                "Gobelin",
        Niveau:             1,
        PointsDeVieMaximum: 500,
        PointsDeVieActuels: 100,
        Degats:             25,
        Initiative:         15,
        RecompenseXP:       150,
        RecompenseArgent:   3,
    }
}