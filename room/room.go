package room

import (
	"go-game/config"
	"go-game/util"
	"image"
	// other necessary imports
)

// RoomType is an enum for different types of rooms.
type RoomType int


type Obstacle struct {
    Rect   image.Rectangle
    IsDoor bool
}
// Room struct represents a room in the game.
type Room struct {
    Obstacles []Obstacle
    RoomType  RoomType
	
    // other room details
}
const (
    RegularRoom RoomType = iota
    ItemRoom
    BossRoom
)
// Obstacle represents an obstacle within a room.



func GenerateRoom(x, y int, direction string, roomType RoomType, roomGrid [config.GridSize][config.GridSize]*Room) *Room {
    if x < 0 || x >= config.GridSize || y < 0 || y >= config.GridSize || roomGrid[x][y] != nil {
        return nil // Bounds check and room existence check
    }
	withDoors := map[string]bool{"top": false, "bottom": false, "left": false, "right": false}


    // Set the door for the direction we are coming from
    if direction != "initial" {
        withDoors[util.OppositeDirection(direction)] = true
    }

    // Check adjacent spaces for potential doors
    if y > 0 && roomGrid[x][y-1] != nil { // Check above
        withDoors["top"] = true
    }
    if y < config.GridSize-1 && roomGrid[x][y+1] != nil { // Check below
        withDoors["bottom"] = true
    }
    if x > 0 && roomGrid[x-1][y] != nil { // Check left
        withDoors["left"] = true
    }
    if x < config.GridSize-1 && roomGrid[x+1][y] != nil { // Check right
        withDoors["right"] = true
    }

    // Add doors for adjacent ungenerated rooms
    if y > 0 && roomGrid[x][y-1] == nil {
        withDoors["top"] = true
    }
    if y < config.GridSize-1 && roomGrid[x][y+1] == nil {
        withDoors["bottom"] = true
    }
    if x > 0 && roomGrid[x-1][y] == nil {
        withDoors["left"] = true
    }
    if x < config.GridSize-1 && roomGrid[x+1][y] == nil {
        withDoors["right"] = true
    }


	newRoom := &Room{
        Obstacles: generateObstacles(withDoors),
        RoomType:  roomType,
    }

    roomGrid[x][y] = newRoom
    return newRoom
}



func generateObstacles(withDoors map[string]bool) []Obstacle {
    var obstacles []Obstacle
    doorSize := 60

    // Create obstacles for each side of the room
    for _, side := range []string{"top", "bottom", "left", "right"} {
        isDoor := withDoors[side]
        switch side {
        case "top", "bottom":
            for x := 0; x < config.ScreenWidth; x += 30 {
                if side == "top" && (!isDoor || x < config.ScreenWidth/2-doorSize/2 || x > config.ScreenWidth/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(x, 0, x+30, 30), IsDoor: false})
                }
                if side == "bottom" && (!isDoor || x < config.ScreenWidth/2-doorSize/2 || x > config.ScreenWidth/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(x, config.ScreenHeight-30, x+30, config.ScreenHeight), IsDoor: false})
                }
            }
        case "left", "right":
            for y := 30; y < config.ScreenHeight-30; y += 30 {
                if side == "left" && (!isDoor || y < config.ScreenHeight/2-doorSize/2 || y > config.ScreenHeight/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(0, y, 30, y+30), IsDoor: false})
                }
                if side == "right" && (!isDoor || y < config.ScreenHeight/2-doorSize/2 || y > config.ScreenHeight/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(config.ScreenWidth-30, y, config.ScreenWidth, y+30), IsDoor: false})
                }
            }
        }
    }

    return obstacles
}


