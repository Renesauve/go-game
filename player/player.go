package player

import (
	"fmt"
	"go-game/config"
	"go-game/items"

	"image"
)

// Item represents an item in the game.

// Import the package that defines the Item type

type Player struct {
    X, Y         float64
    Coordinates  [2]int
    Inventory   []items.Item// Use the imported Item type
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
        fmt.Println(item)
        if item.Name == "CatFood" && item.IsShootable {
            velocityX, velocityY := 5.0, 0.0 // Example velocities, adjust as needed
            projectile := items.NewProjectile(p.X, p.Y, velocityX, velocityY)
            items.Projectiles = append(items.Projectiles, projectile)
            p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
            break
        }
    }
}