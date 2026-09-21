package structure

type Monster struct {
	Nom                string
	PointsDeVieMaximum int
	PointsDeVieActuels int
	Attaque            int
	ExpGain            int
}

func InitMonster() Monster {
	return Monster{
		Nom:                "Python",
		PointsDeVieMaximum: 10000,
		PointsDeVieActuels: 200,
		Attaque:            35,
		ExpGain:            50,
	}
}