package player

import (
	"go-game/packages/config"
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

func (p *Player) StartingPositionInNewRoom(direction Direction) (float64, float64) {
	var newX, newY float64

	switch direction {
	case DirectionRight:
		// Start at the left edge of the new room
		newX = 0
		newY = p.Y // Keep the vertical position the same
	case DirectionLeft:
		// Start at the right edge of the new room
		newX = config.ScreenWidth - float64(config.PlayerWidth)
		newY = p.Y // Keep the vertical position the same
	case DirectionUp:
		// Start at the bottom of the new room
		newX = p.X // Keep the horizontal position the same
		newY = config.ScreenHeight - float64(config.PlayerHeight)
	case DirectionDown:
		// Start at the top of the new room
		newX = p.X // Keep the horizontal position the same
		newY = 0
	}

	return newX, newY
}
