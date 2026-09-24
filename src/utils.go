package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func lireEntier() int {
	return lireChoix(lecteur)
}

func lireChoix(reader *bufio.Reader) int {
	fmt.Print("Choix : ")

	ligne, err := reader.ReadString('\n')
	if err != nil && len(ligne) == 0 {
		fmt.Println("\nFin de l'entrée. À bientôt !")
		os.Exit(0)
	}

	choix, err := strconv.Atoi(strings.TrimSpace(ligne))
	if err != nil {
		return -1
	}

	return choix
}
