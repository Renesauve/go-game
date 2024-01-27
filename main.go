package main

import (
	"go-game/config"
	"go-game/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Initialize the game structure
	gameInstance := game.NewGame()

	// Set up Ebiten window properties using constants from config package
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Cool New Game")

	// Start the game
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}




