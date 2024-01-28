package config

const (
	GridSize = 5

	//PLAYER
	PlayerWidth  = 24
	PlayerHeight = 24
	Speed        = 15.0

	//MINIMAP

	MinimapPadding = 15

	//ITEMS
	ItemWidth        = 1040
	ItemHeight       = 1040
	MaxMinimapWidth  = 150 // Maximum width for the minimap
	MaxMinimapHeight = 150 // Maximum height for the minimap
	//INVENTORY
	InventoryBackgroundWidth  = 400
	InventoryBackgroundHeight = 200
	InventoryItemSize         = 48 // Size of each item icon
	InventoryItemPadding      = 50 // Space between items
	InventoryStartX           = 0  // Starting X position for the inventory
	InventoryStartY           = 0  // Starting Y position for the inventory
)

var (
	MinimapRoomSize = calculateMinimapRoomSize(GridSize)
	MinimapWidth    = GridSize * MinimapRoomSize
	MinimapHeight   = GridSize * MinimapRoomSize
)

type GameViewportConfig struct {
	ScreenWidth  int
	ScreenHeight int
	// other viewport-related fields
}

// need to set a maximum screen size and have all viewport in the center of the screen
func (c *GameViewportConfig) UpdateScreenSize(width, height int) {
	c.ScreenWidth = width
	c.ScreenHeight = height
}

func calculateMinimapRoomSize(gridSize int) int {
	roomSizeWidth := MaxMinimapWidth / gridSize
	roomSizeHeight := MaxMinimapHeight / gridSize
	return min(roomSizeWidth, roomSizeHeight) // Use the smaller of the two to ensure it fits
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Usage:
