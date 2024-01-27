package game

import (
	"go-game/packages/config"
	"go-game/packages/items"
	"go-game/packages/player"
	"go-game/packages/utils"
	"image"
	"log"

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
	Player           player.Player
	CurrentRoom      *room.Room
	RoomsVisited     [config.GridSize][config.GridSize]bool
	Rooms            map[string]*room.Room // A map to hold the rooms
	RoomGrid         [config.GridSize][config.GridSize]*room.Room
	State            GameState // Start at the main menu
	RoomManager      *room.RoomManager
	InventoryOpen    bool // Add this field
	previousIPressed bool // Add this field
}

func NewGame() *Game {
	gridMiddleX, gridMiddleY := 2, 2

	// Calculate screen position from grid position
	screenMiddleX := float64(gridMiddleX)*(float64(config.ScreenWidth)/float64(config.GridSize)) + (float64(config.ScreenWidth) / float64(config.GridSize) / 2.0) - float64(config.PlayerWidth)/2.0
	screenMiddleY := float64(gridMiddleY)*(float64(config.ScreenHeight)/float64(config.GridSize)) + (float64(config.ScreenHeight) / float64(config.GridSize) / 2.0) - float64(config.PlayerHeight)/2.0

	p := player.NewPlayer(screenMiddleX, screenMiddleY, [2]int{config.GridSize / 2, config.GridSize / 2})
	roomGrid := [config.GridSize][config.GridSize]*room.Room{}

	g := &Game{
		Player:       p,
		Rooms:        make(map[string]*room.Room),
		RoomsVisited: [config.GridSize][config.GridSize]bool{},
		RoomGrid:     roomGrid,
		State:        MainMenu, // Start at the main menu
		RoomManager:  room.NewRoomManager(),
		CurrentRoom:  nil,
	}

	initialRoom := g.RoomManager.RoomGrid[gridMiddleX][gridMiddleY]
	battleAxe := items.NewBattleAxe(100, 200)
	initialRoom.Items = append(initialRoom.Items, battleAxe)
	if initialRoom != nil {
		initialRoom.Items = append(initialRoom.Items, battleAxe)
	}

	g.CurrentRoom = initialRoom

	return g
}

func (g *Game) Update() error {

	g.handlePlayerMovement()

	g.RoomsVisited[g.Player.Coordinates[0]][g.Player.Coordinates[1]] = true

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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	minimapStartX := config.ScreenWidth - config.GridSize*config.MinimapRoomSize - config.MinimapPadding
	minimapStartY := config.ScreenHeight - config.GridSize*config.MinimapRoomSize - config.MinimapPadding
	if g.InventoryOpen {
		g.DrawInventory(screen)
	}
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
			minimapX := minimapStartX + x*config.MinimapRoomSize
			minimapY := minimapStartY + y*config.MinimapRoomSize

			// Draw the room on the minimap
			vector.DrawFilledRect(screen, float32(minimapX), float32(minimapY), float32(config.MinimapRoomSize), float32(config.MinimapRoomSize), roomColor, true)
		}
	}

	vector.DrawFilledRect(screen, float32(g.Player.X), float32(g.Player.Y), config.PlayerWidth, config.PlayerHeight, color.White, false)

	if g.CurrentRoom != nil {
		for _, item := range g.CurrentRoom.Items {
			gfxPath := item.GetGFX()
			itemImage, err := utils.LoadImage(gfxPath)
			if err != nil {
				log.Printf("Failed to load item image: %v", err)
				continue
			}

			// Use item's coordinates directly
			x, y := item.GetX(), item.GetY()

			// scale

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(itemImage, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// For simplicity, let's return the outside size; you can adjust as needed for your game
	return outsideWidth, outsideHeight
}

func (g *Game) handlePlayerMovement() (float64, float64) {
	proposedX, proposedY := g.Player.X, g.Player.Y
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		proposedX += config.Speed
		if proposedX > config.ScreenWidth-float64(config.PlayerWidth) && g.Player.Coordinates[0] < config.GridSize-1 {
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
				proposedX = config.ScreenWidth - float64(config.PlayerWidth) // Start at the right edge of the new room
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
				proposedY = config.ScreenHeight - float64(config.PlayerHeight) // Start at the bottom of the new room
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		proposedY += config.Speed
		if proposedY > config.ScreenHeight-float64(config.PlayerHeight) && g.Player.Coordinates[1] < config.GridSize-1 {
			newRoom, newX, newY := g.RoomManager.GetRoomInDirection(g.Player.Coordinates[0], g.Player.Coordinates[1], player.DirectionDown)
			if newRoom != nil {
				g.CurrentRoom = newRoom
				g.Player.Coordinates[0], g.Player.Coordinates[1] = newX, newY
				proposedY = 0 // Start at the top of the new room
			}
		}
	}

	if proposedX >= 0 && proposedX <= config.ScreenWidth-float64(config.PlayerWidth) &&
		proposedY >= 0 && proposedY <= config.ScreenHeight-float64(config.PlayerHeight) {
		g.Player.X, g.Player.Y = proposedX, proposedY
	}
	return proposedX, proposedY

}

func (g *Game) playerIsOverItem() (items.Itemizable, int) {
	for index, item := range g.CurrentRoom.Items {

		if g.Player.X < item.GetX()+config.ItemWidth &&
			g.Player.X+config.PlayerWidth > item.GetX() &&
			g.Player.Y < item.GetY()+config.ItemHeight &&
			g.Player.Y+config.PlayerHeight > item.GetY() {

			return item, index // Return the item and its index
		}
	}
	return nil, -1 // No item is under the player, return -1 as the index
}

func (g *Game) DrawInventory(screen *ebiten.Image) {

	invBackground := image.Rect(0, 0, config.InventoryBackgroundWidth, config.InventoryBackgroundHeight)
	invBackgroundImg := ebiten.NewImageFromImage(image.NewRGBA(invBackground))
	invBackgroundImg.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180}) // Semi-transparent black background

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(config.InventoryStartX, config.InventoryStartY)
	screen.DrawImage(invBackgroundImg, op)

	// Draw each item in the inventory
	for i, item := range g.Player.Inventory.Items {
		// Calculate item position
		x := float64(i%4)*(config.InventoryItemSize+config.InventoryItemPadding) + config.InventoryStartX + config.InventoryItemPadding
		y := float64(i/4)*(config.InventoryItemSize+config.InventoryItemPadding) + config.InventoryStartY + config.InventoryItemPadding

		// Load the item image
		itemImage, err := utils.LoadImage(item.GetGFX())
		if err != nil {
			log.Printf("Failed to load item image: %v", err)
			continue
		}

		// Draw the item
		itemOp := &ebiten.DrawImageOptions{}
		itemOp.GeoM.Translate(x, y)
		itemOp.GeoM.Scale(config.InventoryItemSize/float64(itemImage.Bounds().Dx()), config.InventoryItemSize/float64(itemImage.Bounds().Dy())) // Scale the image to fit the item size
		screen.DrawImage(itemImage, itemOp)
	}
}

// GetRoomStartPosition returns the top-left position of the room on the screen
