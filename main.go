package main

import (
	"go-game/packages/game"
	"go-game/packages/xmlparser"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Initialize the game structure
	xmlparser.ParseItemsXML("assets/gfx/items.xml")

	gameInstance := game.NewGame()
	// Set up Ebiten window properties using constants from config package
	ebiten.SetWindowSize(gameInstance.ViewportConfig.ScreenWidth, gameInstance.ViewportConfig.ScreenHeight)
	ebiten.SetWindowTitle("Cool New Game")

	// Start the game
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
