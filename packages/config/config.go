package config

const (
	GridSize        = 2
	MaxScreenWidth  = 1920
	MaxScreenHeight = 1080
	//PLAYER
	PlayerWidth  = 24
	PlayerHeight = 24
	Speed        = 15.0

	//MINIMAP

	MinimapPadding = 15

	//ITEMS

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
	OffsetX      int // Horizontal offset for centering
	OffsetY      int // Vertical offset for centering
	// other viewport-related fields
}

// need to set a maximum screen size and have all viewport in the center of the screen
func (c *GameViewportConfig) UpdateScreenSize(width, height int) {
	// Ensure the width and height do not exceed the maximum
	if width > MaxScreenWidth {
		width = MaxScreenWidth
	}
	if height > MaxScreenHeight {
		height = MaxScreenHeight
	}

	// Calculate the offsets for centering
	c.OffsetX = (MaxScreenWidth - width) / 2
	c.OffsetY = (MaxScreenHeight - height) / 2

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
