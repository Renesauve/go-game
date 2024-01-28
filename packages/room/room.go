package room

import (
	"fmt"
	"go-game/packages/config"
	"go-game/packages/items"
	"go-game/packages/player"
	"math/rand"
	// other necessary imports
)

// RoomType is an enum for different types of rooms.
type RoomType int

const (
	RegularRoom RoomType = iota
	BossRoom
	ItemRoom
	// Add additional room types as needed
)

// Room struct represents a room in the game.
type Room struct {
	RoomType RoomType
	Items    []items.Itemizable
}
type RoomManager struct {
	Rooms       [][config.GridSize]*Room
	RoomGrid    [config.GridSize][config.GridSize]*Room
	CurrentRoom *Room
	// Other fields as needed
}

// Obstacle represents an obstacle within a room.

func NewRoomManager() *RoomManager {
	rm := &RoomManager{
		Rooms: make([][config.GridSize]*Room, config.GridSize),
	}
	rm.GenerateRooms() // Method to generate rooms
	return rm
}

func (rm *RoomManager) GetRoomInDirection(currentX, currentY int, direction player.Direction) (*Room, int, int) {
	nextX, nextY := currentX, currentY
	switch direction {
	case player.DirectionUp:
		nextY--
	case player.DirectionDown:
		nextY++
	case player.DirectionLeft:
		nextX--
	case player.DirectionRight:
		nextX++
	}

	if nextX < 0 || nextX >= config.GridSize || nextY < 0 || nextY >= config.GridSize {
		return nil, currentX, currentY // Prevent leaving the grid, return nil and current coordinates
	}
	if nextRoom := rm.RoomGrid[nextX][nextY]; nextRoom != nil {
		rm.CurrentRoom = nextRoom
		return nextRoom, nextX, nextY
	}

	return nil, currentX, currentY
}

func (rm *RoomManager) GenerateRooms() {

	// Place the boss room randomly on an edge
	bossRoomX, bossRoomY := randomEdgePosition(config.GridSize)
	rm.RoomGrid[bossRoomX][bossRoomY] = &Room{RoomType: BossRoom}

	// Determine the number of item rooms based on gridSize
	numItemRooms := calculateNumberOfItemRooms(config.GridSize)

	// Place item rooms
	for i := 0; i < numItemRooms; i++ {
		var itemRoomX, itemRoomY int
		for {

			itemRoomX, itemRoomY = rand.Intn(config.GridSize), rand.Intn(config.GridSize)
			// Ensure the item room is not in the center, not on the boss room, and not overlapping another item room
			if !(itemRoomX == config.GridSize/2 && itemRoomY == config.GridSize/2) &&
				!(itemRoomX == bossRoomX && itemRoomY == bossRoomY) &&
				rm.RoomGrid[itemRoomX][itemRoomY] == nil {
				fmt.Println("itemRoomX", itemRoomX, "itemRoomY", itemRoomY)
				fmt.Println("bossRoomX", bossRoomX, "bossRoomY", bossRoomY)

				break
			}
		}
		rm.RoomGrid[itemRoomX][itemRoomY] = &Room{RoomType: ItemRoom}
	}

	// Fill the remaining grid with regular rooms
	for x := 0; x < config.GridSize; x++ {
		for y := 0; y < config.GridSize; y++ {
			if rm.RoomGrid[x][y] == nil {
				rm.RoomGrid[x][y] = &Room{RoomType: RegularRoom}
			}
		}
	}
}

func randomEdgePosition(gridSize int) (int, int) {
	edge := rand.Intn(4)
	var x, y int
	switch edge {
	case 0: // Top edge
		x, y = rand.Intn(gridSize), 0
	case 1: // Bottom edge
		x, y = rand.Intn(gridSize), gridSize-1
	case 2: // Left edge
		x, y = 0, rand.Intn(gridSize)
	case 3: // Right edge
		x, y = gridSize-1, rand.Intn(gridSize)
	}
	return x, y
}

func calculateNumberOfItemRooms(gridSize int) int {
	// Example calculation, adjust as needed
	return gridSize / 2 // Half the gridSize, for example
}
