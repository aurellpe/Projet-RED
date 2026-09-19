package main

import "fmt"

type Enemy struct {
	Name      string
	MaxHP     int
	CurrentHP int
	Distance  int
	State     string
}

func EnemyAI(enemy *Enemy) {
	if enemy.CurrentHP <= 20 {
		enemy.State = "Fuite"
	} else if enemy.Distance <= 2 {
		enemy.State = "Attaque"
	} else if enemy.Distance <= 5 {
		enemy.State = "Poursuite"
	} else {
		enemy.State = "Patrouille"
	}
}

func main() {
	enemy := Enemy{
		Name:      "silverhand",
		MaxHP:     100,
		CurrentHP: 80,
		Distance:  4,
	}

	EnemyAI(&enemy)

	fmt.Println("Ennemi :", enemy.Name)
	fmt.Println("Etat :", enemy.State)
}