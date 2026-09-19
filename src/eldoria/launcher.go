package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Eldoria")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
