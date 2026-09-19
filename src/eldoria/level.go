package main

type Level struct {
	Number  int
	Name    string
	HasBoss bool
	GroundY float64
}

func CreateLevels() []Level {
	return []Level{
		{
			Number:  1,
			Name:    "Foret",
			HasBoss: false,
			GroundY: 156,
		},
		{
			Number:  2,
			Name:    "Ruines",
			HasBoss: false,
			GroundY: 148,
		},
		{
			Number:  3,
			Name:    "Chateau",
			HasBoss: false,
			GroundY: 150,
		},
		{
			Number:  4,
			Name:    "Boss",
			HasBoss: true,
			GroundY: 151,
		},
	}
}
