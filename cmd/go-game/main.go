package main

import (
	"go-game/packages/game"
	"go-game/packages/items"
	"go-game/packages/socket"
	"go-game/packages/xmlparser"
	"log" // This is the line you need to add

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	go socket.StartWebSocketServer()

	itemsList, err := xmlparser.ParseItemsXML("assets/gfx/items.xml")
	if err != nil {
		log.Fatalf("Failed to parse items: %v", err)
	}

	// Convert parsed items to a slice of Itemizable
	var allItems []items.Itemizable
	for _, weapon := range itemsList.Weapons {
		allItems = append(allItems, weapon)
	}
	for _, armor := range itemsList.Armors {
		allItems = append(allItems, armor)
	}

	gameInstance := game.NewGame(allItems) // Pass the first connection from the slice
	// Set up Ebiten window properties using constants from config package
	ebiten.SetWindowSize(gameInstance.ViewportConfig.ScreenWidth, gameInstance.ViewportConfig.ScreenHeight)
	ebiten.SetWindowTitle("Cool New Game")

	// Start the game
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}

}
