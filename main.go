package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
