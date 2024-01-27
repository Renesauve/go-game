package room

import (
	"go-game/packages/config"
	"go-game/packages/items"
	"go-game/packages/player"
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
	// Assuming the center of the grid is always a regular room
	startX, startY := config.GridSize/2, config.GridSize/2
	rm.RoomGrid[startX][startY] = &Room{RoomType: RegularRoom}

	// Iterate through the grid and generate rooms
	for x := 0; x < config.GridSize; x++ {
		for y := 0; y < config.GridSize; y++ {
			// Skip if room is already initialized
			if rm.RoomGrid[x][y] != nil {
				continue
			}

			// Determine the room type based on the positions
			roomType := RegularRoom

			// Create the new room and assign it to the grid
			rm.RoomGrid[x][y] = &Room{
				RoomType: roomType,
			}
		}
	}

}
