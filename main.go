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
	gameInstance.GenerateRooms()
	// Set up Ebiten window properties using constants from config package
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Your Game Title")

	// Start the game
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}




// func main() {
//     game := &Game{
//         rooms: make(map[string]*Room),
//         player: Player{
//             x: float64(screenWidth) / 2 - playerWidth / 2,
//             y: float64(screenHeight) / 2 - playerHeight / 2,
//             inventory: make(map[string]int),
//             coordinates: [2]int{gridSize / 2, gridSize / 2}, // Start in the middle of the grid
//         },
// 		roomsVisited: [gridSize][gridSize]bool{},
//     }

//     generateRooms() // Pre-generate all rooms

//     game.currentRoom = roomGrid[gridSize/2][gridSize/2] // Set initial room

//     ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
//     ebiten.SetWindowSize(screenWidth, screenHeight)

//     if err := ebiten.RunGame(game); err != nil {
//         log.Fatal(err)
//     }
// }