package room

import (
	"go-game/config"
	"image"
	// other necessary imports
)

// RoomType is an enum for different types of rooms.
type RoomType int


type Obstacle struct {
    Rect   image.Rectangle
  
}
// Room struct represents a room in the game.
type Room struct {
    Obstacles []Obstacle
    RoomType  RoomType
	StartingRoom bool
    // other room details
}
const (
    RegularRoom RoomType = iota
    ItemRoom
    BossRoom
)
// Obstacle represents an obstacle within a room.



func GenerateRoom(x, y int, roomType RoomType, roomGrid [config.GridSize][config.GridSize]*Room) *Room {
    if x < 0 || x >= config.GridSize || y < 0 || y >= config.GridSize || roomGrid[x][y] != nil {
        return nil // Bounds check and room existence check
    }

    newRoom := &Room{
        RoomType: roomType,
    }

    roomGrid[x][y] = newRoom
    return newRoom
}





