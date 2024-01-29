package main

import (
	"fmt"
	"go-game/packages/game"
	"go-game/packages/items"
	"go-game/packages/xmlparser"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

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
	fmt.Println(allItems)
	gameInstance := game.NewGame(allItems)
	// Set up Ebiten window properties using constants from config package
	ebiten.SetWindowSize(gameInstance.ViewportConfig.ScreenWidth, gameInstance.ViewportConfig.ScreenHeight)
	ebiten.SetWindowTitle("Cool New Game")

	// Start the game
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
