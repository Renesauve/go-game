package player

import (
	"go-game/config"
	"go-game/items"

	"image"
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
    X, Y         float64
    Coordinates  [2]int
    Inventory   []items.Item// Use the imported Item type
    Facing       Direction // Add this field to track the facing direction
    // ... other fields
}

func NewPlayer(x, y float64, coordinates [2]int) Player {
    return Player{
        X: x,
        Y: y,
        Coordinates: [2]int{coordinates[0], coordinates[1]},
    }
}


// And a method to remove items from the inventory (if needed)


func (p *Player) CheckCollisionWithItem(item *items.Item) bool {


    playerRect := image.Rect(
        int(p.X),
        int(p.Y),
        int(p.X)+config.PlayerWidth,
        int(p.Y)+config.PlayerHeight,
    )
    itemRect := image.Rect(
        item.Position.X,
        item.Position.Y,
        item.Position.X+item.Width,
        item.Position.Y+item.Height,
    )

    // Debug output for bounding boxes
  

    return playerRect.Overlaps(itemRect)
}

func (p *Player) ThrowCatFood() {
  
    for i, item := range p.Inventory {
    
        if item.Name == "CatFood" && item.IsShootable {
            // Determine the velocity based on the facing direction
            var velocityX, velocityY float64
            switch p.Facing {
            case DirectionUp:
                velocityY = -5.0 // Example velocity going up
            case DirectionDown:
                velocityY = 5.0 // Example velocity going down
            case DirectionLeft:
                velocityX = -5.0 // Example velocity going left
            case DirectionRight:
                velocityX = 5.0 // Example velocity going right
            }

            // Create and throw the projectile
            projectile := items.NewProjectile(p.X, p.Y, velocityX, velocityY)
            items.Projectiles = append(items.Projectiles, projectile)
            p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
            break
        }
    }
}