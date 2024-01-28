package player

import (
	"go-game/packages/items"
)

// Item represents an item in the game.

// Import the package that defines the Item type

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

type Player struct {
	X, Y        float64
	Coordinates [2]int
	Facing      Direction
	Armor       items.Armor
	Weapon      items.Weapon
	Inventory   items.Inventory // Add inventory to the player
}

func NewPlayer(x, y float64, coordinates [2]int) Player {
	return Player{
		X:           x,
		Y:           y,
		Coordinates: [2]int{coordinates[0], coordinates[1]},
	}
}
