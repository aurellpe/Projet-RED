package main

import (
    "fmt"
    "strconv"
)

func lireEntier() int {
    var saisie string

    fmt.Print("Votre choix : ")
    fmt.Scan(&saisie)

    n, err := strconv.Atoi(saisie)
    if err != nil {
        return -1
    }

    return n
}
