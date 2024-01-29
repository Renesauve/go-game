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

func NewRoomManager(allItems []items.Itemizable) *RoomManager {
	rm := &RoomManager{
		Rooms: make([][config.GridSize]*Room, config.GridSize),
	}
	rm.GenerateRooms(allItems) // Method to generate rooms
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

func (rm *RoomManager) GenerateRooms(allItems []items.Itemizable) {
	shuffleItems(allItems)
	fmt.Println("allItems:", allItems)
	itemSpawnCount := make(map[items.Itemizable]int)

	bossRoomX, bossRoomY := randomEdgePosition(config.GridSize)
	rm.RoomGrid[bossRoomX][bossRoomY] = &Room{RoomType: BossRoom}

	fmt.Println("Boss Room generated at:", bossRoomX, bossRoomY)

	numItemRooms := calculateNumberOfItemRooms(config.GridSize)
	fmt.Println("Number of Item Rooms to generate:", numItemRooms)
	for i := 0; i < numItemRooms; i++ {
		itemRoomX, itemRoomY := randomPosition(config.GridSize, bossRoomX, bossRoomY, rm.RoomGrid)
		fmt.Println("itemRoomX, itemRoomY:", itemRoomX, itemRoomY)

		for _, item := range allItems {
			if itemSpawnCount[item] < 2 {
				rm.RoomGrid[itemRoomX][itemRoomY] = &Room{
					RoomType: ItemRoom,
					Items:    []items.Itemizable{item},
				}
				itemSpawnCount[item]++
				fmt.Printf("Item Room with item '%s' generated at: %d, %d\n", item.GetName(), itemRoomX, itemRoomY)
				break
			}
		}
	}

	// Fill the remaining grid with regular rooms
	for x := 0; x < config.GridSize; x++ {
		for y := 0; y < config.GridSize; y++ {
			if rm.RoomGrid[x][y] == nil {
				rm.RoomGrid[x][y] = &Room{RoomType: RegularRoom}
				fmt.Println("Regular Room generated at:", x, y)
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

func shuffleItems(items []items.Itemizable) {
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

func randomPosition(gridSize, excludeX, excludeY int, grid [config.GridSize][config.GridSize]*Room) (int, int) {
	var x, y int
	for {
		x, y = rand.Intn(gridSize), rand.Intn(gridSize)
		if (x != excludeX || y != excludeY) && grid[x][y] == nil {
			break
		}
	}
	return x, y
}
