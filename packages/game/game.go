package game

import (
	"fmt"
	"go-game/packages/config"
	"go-game/packages/items"
	"go-game/packages/player"
	"go-game/packages/socket"
	"go-game/packages/utils"
	"image"
	"log"
	"math"

	"go-game/packages/room"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	// other imports as necessary
)

const (
	MainMenu GameState = iota
	InGame
	GameOver
	InventoryView
)

type GameState int
type Game struct {
	Player                   player.Player
	CurrentRoom              *room.Room
	RoomsVisited             [config.GridSize][config.GridSize]bool
	Rooms                    map[string]*room.Room // A map to hold the rooms
	RoomGrid                 [config.GridSize][config.GridSize]*room.Room
	State                    GameState // Start at the main menu
	RoomManager              *room.RoomManager
	ViewportConfig           config.GameViewportConfig
	previousIPressed         bool          // Add this field
	previousMPressed         bool          // Add this field
	InventoryOpen            bool          // Add this field
	InventoryBackgroundImage *ebiten.Image // Add this field
	MinimapOpen              bool          // Add this field
	MinimapImage             *ebiten.Image // Add this field
	WallSpriteImage          *ebiten.Image
	WallWidth                int
}

func NewGame(allItems []items.Itemizable) *Game {

	viewportConfig := config.GameViewportConfig{
		ScreenWidth:  1920, // Default width, or use ebiten.WindowSize() if window is already created
		ScreenHeight: 1080, // Default height
	}

	gridMiddleX, gridMiddleY := config.GridSize/2, config.GridSize/2

	// Calculate screen position from grid position
	screenMiddleX := float64(gridMiddleX)*(float64(viewportConfig.ScreenWidth)/float64(config.GridSize)) +
		(float64(viewportConfig.ScreenWidth) / float64(config.GridSize) / 2.0) -
		float64(config.PlayerWidth)/2.0

	screenMiddleY := float64(gridMiddleY)*(float64(viewportConfig.ScreenHeight)/float64(config.GridSize)) +
		(float64(viewportConfig.ScreenHeight) / float64(config.GridSize) / 2.0) -
		float64(config.PlayerHeight)/2.0
	p := player.NewPlayer(screenMiddleX, screenMiddleY, [2]int{gridMiddleX, gridMiddleY})
	roomGrid := [config.GridSize][config.GridSize]*room.Room{}

	g := &Game{
		ViewportConfig: viewportConfig,
		Player:         p,
		Rooms:          make(map[string]*room.Room),
		RoomsVisited:   [config.GridSize][config.GridSize]bool{},
		RoomGrid:       roomGrid,
		State:          MainMenu, // Start at the main menu
		RoomManager:    room.NewRoomManager(allItems),
		CurrentRoom:    nil,
		WallWidth:      100, // set the wall dimensions

	}
	wallSprite, err := utils.LoadImage("wallsprite.png")
	if err != nil {
		log.Fatalf("Failed to load wall sprite: %v", err)
	}
	g.WallSpriteImage = wallSprite

	initialRoom := g.RoomManager.RoomGrid[gridMiddleX][gridMiddleY]

	g.CurrentRoom = initialRoom
	fmt.Print(socket.Clients)
	return g
}

func (g *Game) Update() error {

	screenWidth, screenHeight := ebiten.WindowSize()

	g.ViewportConfig.UpdateScreenSize(screenWidth, screenHeight)
	g.handlePlayerMovement()
	g.processInput()
	g.RoomsVisited[g.Player.Coordinates[0]][g.Player.Coordinates[1]] = true
	// fmt.Println(g.CurrentRoom.Items)
	if g.CurrentRoom != nil {
		item, index := g.playerIsOverItem()

		if item != nil {

			// Add the item to the player's inventory
			g.Player.Inventory.AddItem(item)

			// Remove the item from the room
			g.CurrentRoom.Items = append(g.CurrentRoom.Items[:index], g.CurrentRoom.Items[index+1:]...)

		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		if !g.previousIPressed {
			g.InventoryOpen = !g.InventoryOpen
			g.previousIPressed = true
		}
	} else {
		g.previousIPressed = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		if !g.previousMPressed {
			g.MinimapOpen = !g.MinimapOpen
			g.previousMPressed = true
		}
	} else {
		g.previousMPressed = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(g.Player.X), float32(g.Player.Y), config.PlayerWidth, config.PlayerHeight, color.White, false)
	g.drawWalls(screen)

	if g.CurrentRoom != nil && g.CurrentRoom.RoomType == room.ItemRoom {
		for _, item := range g.CurrentRoom.Items {
			gfxPath := item.GetGFX()
			itemImage, err := utils.LoadImage(gfxPath)
			if err != nil {
				log.Printf("Failed to load item image: %v", err)
				continue
			}

			// Calculate the position to draw the item in the center of the screen
			screenWidth := float64(g.ViewportConfig.ScreenWidth)
			screenHeight := float64(g.ViewportConfig.ScreenHeight)
			itemX := (screenWidth - float64(itemImage.Bounds().Dx())) / 2
			itemY := (screenHeight - float64(itemImage.Bounds().Dy())) / 2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(itemX, itemY)
			screen.DrawImage(itemImage, op)
		}
	}

	if g.InventoryOpen {
		g.drawInventory(screen)
	}
	if g.MinimapOpen {
		g.drawMinimap(screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ViewportConfig.ScreenWidth, g.ViewportConfig.ScreenHeight
}

func (g *Game) processInput() {
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		g.toggleInventory()
	}
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.toggleMinimap()
	}
	g.handlePlayerMovement()
}

func (g *Game) toggleInventory() {
	if !g.previousIPressed {
		g.InventoryOpen = !g.InventoryOpen
	}
	g.previousIPressed = true
}

func (g *Game) toggleMinimap() {
	if !g.previousMPressed {
		g.MinimapOpen = !g.MinimapOpen
	}
	g.previousMPressed = true
}

func (g *Game) handlePlayerMovement() {
	proposedX, proposedY := g.Player.X, g.Player.Y

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		proposedX += config.Speed
		if proposedX > float64(g.ViewportConfig.ScreenWidth)-float64(config.PlayerWidth) && g.Player.Coordinates[0] < config.GridSize-1 {
			newRoom, newX, newY := g.RoomManager.GetRoomInDirection(g.Player.Coordinates[0], g.Player.Coordinates[1], player.DirectionRight)
			if newRoom != nil {
				g.CurrentRoom = newRoom
				g.Player.Coordinates[0], g.Player.Coordinates[1] = newX, newY
				proposedX = 0 // Reset to left edge of the new room
			}

		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		proposedX -= config.Speed
		if proposedX < 0 && g.Player.Coordinates[0] > 0 {
			newRoom, newX, newY := g.RoomManager.GetRoomInDirection(g.Player.Coordinates[0], g.Player.Coordinates[1], player.DirectionLeft)
			if newRoom != nil {
				g.CurrentRoom = newRoom
				g.Player.Coordinates[0], g.Player.Coordinates[1] = newX, newY
				proposedX = float64(g.ViewportConfig.ScreenWidth) - float64(config.PlayerWidth) // Start at the right edge of the new room
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		proposedY -= config.Speed
		if proposedY < 0 && g.Player.Coordinates[1] > 0 {
			newRoom, newX, newY := g.RoomManager.GetRoomInDirection(g.Player.Coordinates[0], g.Player.Coordinates[1], player.DirectionUp)
			if newRoom != nil {
				g.CurrentRoom = newRoom
				g.Player.Coordinates[0], g.Player.Coordinates[1] = newX, newY
				proposedY = float64(g.ViewportConfig.ScreenHeight) - float64(config.PlayerHeight) // Start at the bottom of the new room
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		proposedY += config.Speed
		if proposedY > float64(g.ViewportConfig.ScreenHeight)-float64(config.PlayerHeight) && g.Player.Coordinates[1] < config.GridSize-1 {
			newRoom, newX, newY := g.RoomManager.GetRoomInDirection(g.Player.Coordinates[0], g.Player.Coordinates[1], player.DirectionDown)
			if newRoom != nil {
				g.CurrentRoom = newRoom
				g.Player.Coordinates[0], g.Player.Coordinates[1] = newX, newY
				proposedY = 0 // Start at the top of the new room
			}
		}
	}

	minX, maxX, minY, maxY := 0.0, float64(g.ViewportConfig.ScreenWidth)-float64(config.PlayerWidth), 0.0, float64(g.ViewportConfig.ScreenHeight)-float64(config.PlayerHeight)

	// Check if the player is at the edge of the grid
	if g.Player.Coordinates[0] == 0 { // Left edge
		minX = float64(g.WallWidth)
	}
	if g.Player.Coordinates[0] == config.GridSize-1 { // Right edge
		maxX -= float64(g.WallWidth)
	}
	if g.Player.Coordinates[1] == 0 { // Top edge
		minY = float64(g.WallWidth)
	}
	if g.Player.Coordinates[1] == config.GridSize-1 { // Bottom edge
		maxY -= float64(g.WallWidth)
	}

	// Adjust proposedX and proposedY to ensure they don't cross the wall boundaries
	if proposedX < minX {
		proposedX = minX
	} else if proposedX > maxX {
		proposedX = maxX
	}
	if proposedY < minY {
		proposedY = minY
	} else if proposedY > maxY {
		proposedY = maxY
	}

	// Update the player's position
	g.Player.X, g.Player.Y = proposedX, proposedY
}

func (g *Game) playerIsOverItem() (items.Itemizable, int) {
	for index, item := range g.CurrentRoom.Items {
		// Central position of the item
		itemX := (float64(g.ViewportConfig.ScreenWidth) - float64(200)) / 2
		itemY := (float64(g.ViewportConfig.ScreenHeight) - float64(200)) / 2

		// Define interaction area (assuming item.GetWidth() and item.GetHeight() exist)
		itemRect := image.Rect(
			int(itemX), int(itemY),
			int(itemX)+200, int(itemY)+200,
		)

		// Player's current position and size
		playerRect := image.Rect(
			int(g.Player.X), int(g.Player.Y),
			int(g.Player.X)+config.PlayerWidth, int(g.Player.Y)+config.PlayerHeight,
		)

		// Check if player overlaps the item's interaction area
		if itemRect.Overlaps(playerRect) {
			return item, index
		}
	}
	return nil, -1
}

func (g *Game) drawInventory(screen *ebiten.Image) {
	// Instead of getting bounds every time, use the stored screen dimensions from the viewport config
	screenWidth, screenHeight := g.ViewportConfig.ScreenWidth, g.ViewportConfig.ScreenHeight

	// Calculate the starting X and Y coordinates for the inventory
	inventoryStartX := float64(screenWidth - config.InventoryBackgroundWidth)
	inventoryStartY := float64(screenHeight - config.InventoryBackgroundHeight)

	// Use a single draw image options for the entire inventory drawing for efficiency
	invOp := &ebiten.DrawImageOptions{}
	invOp.GeoM.Translate(inventoryStartX, inventoryStartY)

	// Draw the inventory background
	// It's more efficient to create a single background image and reuse it rather than creating a new one each frame
	if g.InventoryBackgroundImage == nil {
		g.InventoryBackgroundImage = ebiten.NewImage(config.InventoryBackgroundWidth, config.InventoryBackgroundHeight)
		g.InventoryBackgroundImage.Fill(color.RGBA{R: 210, G: 180, B: 140, A: 255}) // Adjust alpha for desired transparency
	}
	screen.DrawImage(g.InventoryBackgroundImage, invOp)

	// Draw each item in the inventory
	for i, item := range g.Player.Inventory.Items {
		itemImage, err := utils.LoadImage(item.GetGFX())
		if err != nil {
			log.Printf("Failed to load item image: %v", err)
			continue
		}

		// Calculate item position
		x := float64(i%4)*(config.InventoryItemSize+config.InventoryItemPadding) + config.InventoryItemPadding
		y := float64(i/4)*(config.InventoryItemSize+config.InventoryItemPadding) + config.InventoryItemPadding

		// Reuse the existing draw options for efficiency by resetting and translating for each item
		itemOp := &ebiten.DrawImageOptions{}
		itemOp.GeoM.Translate(inventoryStartX+x, inventoryStartY+y)
		itemOp.GeoM.Scale(config.InventoryItemSize/float64(itemImage.Bounds().Dx()), config.InventoryItemSize/float64(itemImage.Bounds().Dy())) // Scale the image to fit the item size
		screen.DrawImage(itemImage, itemOp)
	}
}

func (g *Game) drawMinimap(screen *ebiten.Image) {
	// Create the minimap image once if it doesn't exist
	if g.MinimapImage == nil {
		g.MinimapImage = ebiten.NewImage(config.MinimapWidth, config.MinimapHeight)
	}

	// Clear the minimap image to start fresh each frame
	g.MinimapImage.Clear()

	// Drawing logic for the minimap
	for x := 0; x < config.GridSize; x++ {
		for y := 0; y < config.GridSize; y++ {
			var roomColor color.Color

			// Determine room color
			if g.RoomsVisited[x][y] {
				roomColor = color.Gray{Y: 150} // Visited room color
			} else {
				roomColor = color.Gray{Y: 50} // Unvisited room color
			}

			// Highlight the current room
			if x == g.Player.Coordinates[0] && y == g.Player.Coordinates[1] {
				roomColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Current room color
			}

			// Calculate the position of each room in the minimap
			minimapX := x * config.MinimapRoomSize
			minimapY := y * config.MinimapRoomSize

			// Draw the room on the minimap image
			vector.DrawFilledRect(g.MinimapImage, float32(minimapX), float32(minimapY), float32(config.MinimapRoomSize), float32(config.MinimapRoomSize), roomColor, true)
		}
	}

	// Draw the minimap image onto the screen
	minimapOp := &ebiten.DrawImageOptions{}
	minimapStartX := g.ViewportConfig.ScreenWidth - config.MinimapWidth - config.MinimapPadding
	minimapStartY := config.MinimapPadding
	minimapOp.GeoM.Translate(float64(minimapStartX), float64(minimapStartY))
	screen.DrawImage(g.MinimapImage, minimapOp)
}

// GetRoomStartPosition returns the top-left position of the room on the screen

func (g *Game) drawWalls(screen *ebiten.Image) {
	// Constants for the dimensions of each room and the grid size

	// Variables for the width and height of the wall sprite
	wallSpriteWidth, wallSpriteHeight := 100, 163

	// Player's current room coordinates
	playerX, playerY := g.Player.Coordinates[0], g.Player.Coordinates[1]

	// Rotate the sprite by 90 degrees for top and bottom walls
	rotate90 := ebiten.GeoM{}
	rotate90.Rotate(math.Pi / 2) // Rotate by 90 degrees

	// Draw top wall if the player is in the first row
	if playerY == 0 {
		for x := 0; x < g.ViewportConfig.ScreenWidth; x += wallSpriteHeight { // Notice we use wallSpriteHeight after rotation
			op := &ebiten.DrawImageOptions{}
			op.GeoM = rotate90
			op.GeoM.Translate(float64(x+wallSpriteHeight), 0) // Adjust the position after rotation
			screen.DrawImage(g.WallSpriteImage, op)
		}
	}

	// Draw bottom wall if the player is in the last row
	if playerY == config.GridSize-1 {
		for x := 0; x < g.ViewportConfig.ScreenWidth; x += wallSpriteHeight { // Notice we use wallSpriteHeight after rotation
			op := &ebiten.DrawImageOptions{}
			op.GeoM = rotate90
			op.GeoM.Translate(float64(x+wallSpriteHeight), float64(g.ViewportConfig.ScreenHeight-wallSpriteWidth)) // Adjust for sprite width after rotation
			screen.DrawImage(g.WallSpriteImage, op)
		}
	}

	// Draw left wall if the player is in the first column
	if playerX == 0 {
		for y := 0; y < g.ViewportConfig.ScreenHeight; y += wallSpriteHeight {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(y))
			screen.DrawImage(g.WallSpriteImage, op)
		}
	}

	// Draw right wall if the player is in the last column
	if playerX == config.GridSize-1 {
		for y := 0; y < g.ViewportConfig.ScreenHeight; y += wallSpriteHeight {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.ViewportConfig.ScreenWidth-wallSpriteWidth), float64(y))
			screen.DrawImage(g.WallSpriteImage, op)
		}
	}
}
