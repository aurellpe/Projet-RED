package main

import (
	"bufio"
	"fmt"
	"main/structure"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)


type Character struct {
	nom                string
	classe             string
	niveau             int
	pointsDeVieMax     int
	pointsDeVieActuels int
	inventaire         []string
}

func initCharacter(
	nom string,
	prenom string,
	age int,
	classe string,
	niveau int,
	pointsDeVieMaximum int,
	pointsDeVieActuels int,
	inventaire []structure.Item,
	argent int,
	xpmax int,
	xpactu int,
	degats int,
) structure.Character {

	return structure.Character{
		Nom:                nom,
		Prenom:             prenom,
		Age:                age,
		Classe:             classe,
		Niveau:             niveau,
		PointsDeVieMaximum: pointsDeVieMaximum,
		PointsDeVieActuels: pointsDeVieActuels,
		Inventaire:         inventaire,
		Argent:             argent,
		XPmax:              xpmax,
		XPactu:             xpactu,
		Degats:             degats,
	}
}

func displayInfo(player structure.Character) {

	fmt.Println()
	fmt.Println("===== PERSONNAGE =====")
	fmt.Println("Nom :", player.Nom)
	fmt.Println("Prénom :", player.Prenom)
	fmt.Println("Âge :", player.Age)
	fmt.Println("Classe :", player.Classe)
	fmt.Println("Niveau :", player.Niveau)
	fmt.Println("Points de vie :", player.PointsDeVieActuels, "/", player.PointsDeVieMaximum)
	fmt.Println("Arme :", player.Arme)
	fmt.Println("Dégâts :", player.Degats)
	fmt.Println("Argent :", player.Argent)
	fmt.Println("Débris :", player.Debris)
}


func pause() {

	fmt.Println("\nAppuyez sur Entrée pour continuer...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func accessInventory(player structure.Character) {

	fmt.Println()
	fmt.Println("===== INVENTAIRE =====")

	fmt.Printf("Emplacements : %d / 10\n", len(player.Inventaire))
	fmt.Println()

	if len(player.Inventaire) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	for i, objet := range player.Inventaire {

		fmt.Println(
			i+1,
			"-",
			objet.Nom,
			"(x",
			objet.Quantity,
			")",
		)
	}
}


func main() {

	player := createCharacter()

	var run bool = true

	for run {

		menu(&player)

	}
}

func takePot(p structure.Character) {

	for i, it := range p.Inventaire {

		if it.Nom == "potion" && it.Quantity > 0 {

			p.PointsDeVieActuels += 50

			if p.PointsDeVieActuels > p.PointsDeVieMaximum {
				p.PointsDeVieActuels = p.PointsDeVieMaximum
			}

			p.Inventaire[i].Quantity--

			if p.Inventaire[i].Quantity == 0 {
				p.Inventaire = append(
					p.Inventaire[:i],
					p.Inventaire[i+1:]...,
				)
			}

			return
		}
	}
}


func shopMenu(c *structure.Character, scanner *bufio.Scanner) {

	for {

		fmt.Println("\n===== Marchand =====")
		fmt.Println("1. Potion de vie")
		fmt.Println("2. Fourrure de Brice")
		fmt.Println("3. Casque de Guigui")
		fmt.Println("0. Retour")

		scanner.Scan()

		choix := scanner.Text()

		switch choix {

		case "1":
			acheterItem(c, "Potion de vie")

		case "2":
			acheterItem(c, "Fourrure de Brice")

		case "3":
			acheterItem(c, "Casque de Guigui")

		case "0":
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}


func acheterItem(c *structure.Character, nom string) {

	fmt.Println("Vous avez acheté :", nom)
}


type shopItem struct {
	Nom  string
	Prix int
}

var shopItems = []shopItem{

	{Nom: "Potion de vie", Prix: 0},
	{Nom: "Potion de poison", Prix: 6},
	{Nom: "Livre de Sort : Boule de Feu", Prix: 25},
	{Nom: "Fourrure de Loup", Prix: 4},
	{Nom: "Peau de Troll", Prix: 7},
	{Nom: "Cuir de vache", Prix: 1},
}


func shop(c *structure.Character, reader *bufio.Reader) {

	for {

		fmt.Println("\n===== Marchand =====")

		for i, si := range shopItems {

			fmt.Printf(
				"%d. %s (%d pièces d'or)\n",
				i+1,
				si.Nom,
				si.Prix,
			)
		}

		fmt.Println("0. Retour")
		fmt.Println("Votre argent :", c.Argent, "pièces d'or")

		choix := lireChoix(reader)

		if choix == 0 {
			return
		}

		if choix >= 1 && choix <= len(shopItems) {

			si := shopItems[choix-1]

			acheterItemPayant(
				c,
				si.Nom,
				si.Prix,
			)

		} else {

			fmt.Println("Choix invalide.")
		}
	}
}

func acheterItemPayant(
	c *structure.Character,
	nom string,
	prix int,
) {

	if c.Argent < prix {

		fmt.Println(
			"Vous n'avez pas assez d'argent pour acheter :",
			nom,
		)

		return
	}

	if !c.AddInventory(structure.Item{
		Nom:      nom,
		Quantity: 1,
	}) {

		fmt.Println("Votre inventaire est plein !")
		fmt.Println("Maximum : 10 emplacements.")

		return
	}

	c.Argent -= prix

	fmt.Println("Vous avez acheté :", nom)
	fmt.Println("Argent restant :", c.Argent)
}


func lireChoix(reader *bufio.Reader) int {

	fmt.Print("Choix : ")

	ligne, _ := reader.ReadString('\n')
	ligne = strings.TrimSpace(ligne)

	choix, err := strconv.Atoi(ligne)

	if err != nil {
		return -1
	}

	return choix
}


func gainExperience(c *structure.Character, xpGagne int) {

	fmt.Printf(
		"Vous gagnez %d points d'expérience.\n",
		xpGagne,
	)

	c.XPactu += xpGagne

	for c.XPactu >= c.XPmax {

		c.XPactu -= c.XPmax

		levelUp(c)
	}
}


func levelUp(c *structure.Character) {

	c.Niveau++

	bonusPV := 10

	c.PointsDeVieMaximum += bonusPV
	c.PointsDeVieActuels += bonusPV

	c.XPmax = int(float64(c.XPactu) * 1.2)

	fmt.Printf(
		"Niveau supérieur ! Vous êtes maintenant niveau %d.\n",
		c.Niveau,
	)

	fmt.Printf(
		"Bonus : +%d points de vie maximum.\n",
		bonusPV,
	)

	fmt.Printf(
		"PV : %d / %d\n",
		c.PointsDeVieActuels,
		c.PointsDeVieMaximum,
	)
}

func formatName(name string) string {
	words := strings.Fields(strings.ToLower(name))

	for i, word := range words {
		runes := []rune(word)

		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}

func chooseClass(reader *bufio.Reader) string {
	for {
		fmt.Println()
		fmt.Println("===== CHOIX DE LA CLASSE =====")
		fmt.Println("1. Humain")
		fmt.Println("2. Elfe")
		fmt.Println("3. Nain")

		fmt.Print("Choix : ")

		choix, _ := reader.ReadString('\n')
		choix = strings.TrimSpace(choix)

		switch choix {
		case "1":
			return "Humain"
		case "2":
			return "Elfe"
		case "3":
			return "Nain"
		default:
			fmt.Println("Choix invalide. Choisissez 1, 2 ou 3.")
		}
	}
}


func createCharacter() structure.Character {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("===== CRÉATION DU PERSONNAGE =====")

	fmt.Print("Nom : ")
	nom, _ := reader.ReadString('\n')
	nom = formatName(strings.TrimSpace(nom))

	fmt.Print("Prénom : ")
	prenom, _ := reader.ReadString('\n')
	prenom = formatName(strings.TrimSpace(prenom))

	fmt.Print("Âge : ")
	ageStr, _ := reader.ReadString('\n')
	ageStr = strings.TrimSpace(ageStr)

	age, _ := strconv.Atoi(ageStr)

	classe := chooseClass(reader)

	niveau := 1

	var pvMax int

	switch classe {
	case "Humain":
		pvMax = 100
	case "Elfe":
		pvMax = 80
	case "Nain":
		pvMax = 120
	}

	pvActu := pvMax / 2

	inventaire := []structure.Item{}

	argent := 100
	xpmax := 100
	xpactu := 0

	degats := 10

	return initCharacter(
		nom,
		prenom,
		age,
		classe,
		niveau,
		pvMax,
		pvActu,
		inventaire,
		argent,
		xpmax,
		xpactu,
		degats,
	)
}


func convertToCombatCharacter(
	p structure.Character,
) Character {

	return Character{
		nom:                p.Nom,
		classe:             p.Classe,
		niveau:             p.Niveau,
		pointsDeVieMax:     p.PointsDeVieMaximum,
		pointsDeVieActuels: p.PointsDeVieActuels,
		inventaire:         []string{"Potion de vie"},
	}
}


func trainingFightCommand(player *structure.Character) {

    fmt.Println("\nUn gobelin apparaît !")

    goblinPV := 100
    goblinAtk := 25

    for {

        fmt.Println("\n=== Combat ===")
        fmt.Printf("Gobelin : %d PV\n", goblinPV)
        fmt.Printf("%s : %d/%d PV\n", player.Nom, player.PointsDeVieActuels, player.PointsDeVieMaximum)

        fmt.Println("1 - Attaquer")
        fmt.Println("2 - Utiliser une potion")
        fmt.Println("3 - Défendre")
        fmt.Println("4 - Utiliser un sort")
        fmt.Println("0 - Fuir")

        choix := lireEntier()
        defended := false

        switch choix {

        case 1:
            goblinPV -= player.Degats
            fmt.Printf("Vous infligez %d dégâts au gobelin.\n", player.Degats)

        case 2:
            aPotionVie := false
            aPotionPoison := false

            // Vérifier les potions disponibles
            for _, it := range player.Inventaire {
                if it.Nom == "Potion de vie" && it.Quantity > 0 {
                    aPotionVie = true
                }
                if it.Nom == "Potion de poison" && it.Quantity > 0 {
                    aPotionPoison = true
                }
            }

            if !aPotionVie && !aPotionPoison {
                fmt.Println("Vous n'avez aucune potion à utiliser.")
                break
            }

            fmt.Println("\n=== Choisissez une potion ===")
            if aPotionVie {
                fmt.Println("1 - Potion de vie (+50 PV)")
            }
            if aPotionPoison {
                fmt.Println("2 - Potion de poison (20 dégâts au gobelin)")
            }
            fmt.Println("0 - Annuler")

            choixPotion := lireEntier()

            switch choixPotion {

            case 1:
                if !aPotionVie {
                    fmt.Println("Vous n'avez pas de potion de vie.")
                    break
                }

                for i, it := range player.Inventaire {
                    if it.Nom == "Potion de vie" && it.Quantity > 0 {

                        fmt.Println("Vous utilisez une Potion de vie ! +50 PV")
                        player.PointsDeVieActuels += 50

                        if player.PointsDeVieActuels > player.PointsDeVieMaximum {
                            player.PointsDeVieActuels = player.PointsDeVieMaximum
                        }

                        player.Inventaire[i].Quantity--
                        if player.Inventaire[i].Quantity == 0 {
                            player.Inventaire = append(player.Inventaire[:i], player.Inventaire[i+1:]...)
                        }

                        break
                    }
                }

            case 2:
                if !aPotionPoison {
                    fmt.Println("Vous n'avez pas de potion de poison.")
                    break
                }

                for i, it := range player.Inventaire {
                    if it.Nom == "Potion de poison" && it.Quantity > 0 {

                        fmt.Println("Vous utilisez une Potion de poison ! Le gobelin subit 20 dégâts.")
                        goblinPV -= 20

                        player.Inventaire[i].Quantity--
                        if player.Inventaire[i].Quantity == 0 {
                            player.Inventaire = append(player.Inventaire[:i], player.Inventaire[i+1:]...)
                        }

                        break
                    }
                }

            case 0:
                fmt.Println("Vous annulez.")
                break

            default:
                fmt.Println("Choix invalide.")
            }

        case 3:
            fmt.Println("Vous vous mettez en position défensive !")
            defended = true

        case 4:
            sortUtilise := false

            for i, it := range player.Inventaire {
                if it.Nom == "Livre de Sort : Boule de Feu" && it.Quantity > 0 {

                    fmt.Println("Vous lancez Boule de Feu ! Le gobelin subit 40 dégâts.")
                    goblinPV -= 40

                    player.Inventaire[i].Quantity--
                    if player.Inventaire[i].Quantity == 0 {
                        player.Inventaire = append(player.Inventaire[:i], player.Inventaire[i+1:]...)
                    }

                    sortUtilise = true
                    break
                }
            }

            if !sortUtilise {
                fmt.Println("Vous n'avez aucun sort à utiliser.")
            }

        case 0:
            fmt.Println("Vous fuyez le combat.")
            return

        default:
            fmt.Println("Choix invalide.")
        }

        if goblinPV <= 0 {
	fmt.Println("Vous avez vaincu le gobelin !")

	debrisGagnes := rand.Intn(6) + 3 
	orGagne := rand.Intn(11) + 5     

	player.Debris += debrisGagnes
	player.Argent += orGagne

	fmt.Printf("Vous récupérez %d débris !\n", debrisGagnes)
	fmt.Printf("Vous récupérez %d pièces d'or !\n", orGagne)
	fmt.Printf("Argent total : %d pièces d'or\n", player.Argent)
	fmt.Printf("Débris totaux : %d\n", player.Debris)

	return
}

        if defended {
            fmt.Println("Vous bloquez complètement l'attaque du gobelin ! 0 dégâts.")
        } else {
            player.PointsDeVieActuels -= goblinAtk
            fmt.Printf("Le gobelin vous inflige %d dégâts.\n", goblinAtk)
        }

        if player.PointsDeVieActuels <= 0 {
            fmt.Println("Vous êtes mort... Réanimation à 50% PV.")
            player.PointsDeVieActuels = player.PointsDeVieMaximum / 2
            return
        }
    }
}
