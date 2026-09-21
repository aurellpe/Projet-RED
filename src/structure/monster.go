package structure

type Monster struct {
	Nom                 string
	PointsDeVieMaximum  int
	PointsDeVieActuels  int
	PointsDattaque		int
}

func InitMonster() Monster {
    return Monster{
        Nom:                "Python",
        PointsDeVieActuels: 100,
        PointsDeVieMaximum: 10000,
        Attaque:            25,
        ExpGain:            50,
    }
}