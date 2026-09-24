package main

import (
	"fmt"
 	"main/structure"
)

type Monster struct {
	Nom        string
	PVMax      int
	PV         int
	Attaque    int
	Initiative int
	XP         int
	Or         int
	Debris     int
}

func initGoblin() Monster {
	return Monster{Nom: "mirco la menace", PVMax: 150, PV: 150, Attaque: 10, Initiative: 10, XP: 40, Or: 10, Debris: 6}
}

func isDead(p *structure.Character) bool {
	if p.PointsDeVieActuels > 0 {
		return false
	}
	p.PointsDeVieActuels = p.PointsDeVieMaximum / 2
	if p.PointsDeVieActuels < 1 {
		p.PointsDeVieActuels = 1
	}
	fmt.Printf("\n*** %s est mort ! ***\n", p.Nom)
	fmt.Printf("Réanimation avec %d / %d PV.\n", p.PointsDeVieActuels, p.PointsDeVieMaximum)
	return true
}

func goblinPattern(m *Monster, p *structure.Character, tour int) {
	degats := m.Attaque
	if tour%3 == 0 {
		degats *= 2
	}
	p.PointsDeVieActuels -= degats
	if p.PointsDeVieActuels < 0 {
		p.PointsDeVieActuels = 0
	}
	fmt.Printf("%s inflige à %s %d dégâts\n", m.Nom, p.Nom, degats)
	fmt.Printf("%s PV : %d / %d\n", p.Nom, p.PointsDeVieActuels, p.PointsDeVieMaximum)
}

func hitMonster(p *structure.Character, m *Monster, attaque string, degats int) {
	m.PV -= degats
	if m.PV < 0 {
		m.PV = 0
	}
	fmt.Printf("%s utilise %s et inflige %d dégâts à %s\n", p.Nom, attaque, degats, m.Nom)
	fmt.Printf("%s PV : %d / %d\n", m.Nom, m.PV, m.PVMax)
}

func attackMenu(p *structure.Character, m *Monster) bool {
	for {
		fmt.Println("\n--- Attaquer ---")
		fmt.Printf("1. Attaque basique (%d dégâts, gratuit)\n", p.Degats)
		for i, nom := range p.Skills {
			s := grimoire[nom]
			fmt.Printf("%d. %s (%d dégâts, %d mana)\n", i+2, s.Nom, s.Degats, s.Mana)
		}
		fmt.Println("0. Retour")
		switch choix := lireEntier(); {
		case choix == 0:
			return false
		case choix == 1:
			hitMonster(p, m, "Attaque basique", p.Degats)
			return true
		case choix >= 2 && choix < len(p.Skills)+2:
			s := grimoire[p.Skills[choix-2]]
			if p.Mana < s.Mana {
				fmt.Printf("Pas assez de mana pour %s.\n", s.Nom)
				continue
			}
			p.Mana -= s.Mana
			hitMonster(p, m, s.Nom, s.Degats)
			return true
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func characterTurn(p *structure.Character, m *Monster, protege *bool) bool {
	for {
		fmt.Printf("\n--- Tour de %s ---\n1 - Attaquer\n2 - Inventaire\n3 - Défendre\n0 - Fuir\n", p.Nom)
		switch lireEntier() {
		case 1:
			if attackMenu(p, m) { return false }
		case 2:
			if inventoryMenu(p, m) { return false }
		case 3:
			*protege = true
			fmt.Println("Vous vous mettez en position défensive !")
			return false
		case 0:
			return true
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func trainingFightCommand(player *structure.Character) {
	m := initGoblin()
	tour := 1
	protege := false
	fuite := false
	joueurCommence := player.Initiative >= m.Initiative
	for {
		fmt.Printf("\n========== TOUR %d ==========\n", tour)
		fini := false
		if joueurCommence {
			fini = characterTurn(player, &m, &protege)
			if fini {
				fuite = true
			}
			if !fini {
				if protege { protege = false } else { goblinPattern(&m, player, tour) }
			}
		} else {
			if protege { protege = false } else { goblinPattern(&m, player, tour) }
			if player.PointsDeVieActuels <= 0 { fini = true } else { fini = characterTurn(player, &m, &protege) }
		}
		if fini || m.PV <= 0 || player.PointsDeVieActuels <= 0 { break }
		tour++
	}
	if fuite {
		fmt.Println("Vous fuyez le combat.")
	} else if m.PV <= 0 {
		fmt.Printf("Vous avez vaincu %s !\n", m.Nom)
		gainExperience(player, m.XP)
		player.Argent += m.Or
		player.Debris += m.Debris
	} else {
		isDead(player)
		fmt.Println("Vous avez perdu !")
	}
}
