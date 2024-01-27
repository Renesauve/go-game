package items

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// In a new file, items.go

// ItemType is an enum for different types of items.
type ItemType int

const (
	HealthPotion ItemType = iota
	Key
	CatFood
	// Add more items as needed.
)

// Item represents an item in the game.
type Item struct {
	Type        ItemType
	Name        string
	Position    image.Point
	Image       *ebiten.Image
	Width       int
	Height      int
	Collected   bool
	IsShootable bool
}

type Projectile struct {
	X, Y      float64
	VelocityX float64
	VelocityY float64
	// ... other fields like image, direction, etc.
}

var Projectiles []Projectile

// InitializeItem is a function that initializes an item and returns it.
func InitializeItem(itemType ItemType, name string, imagePath string, roomWidth, roomHeight int, IsShootable bool) Item {

	itemImage, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to load item image: %v", err)
	}

	// Calculate the center position
	centerX := roomWidth/2 - itemImage.Bounds().Dx()/2
	centerY := roomHeight/2 - itemImage.Bounds().Dy()/2

	return Item{
		Type:     itemType,
		Name:     name,
		Position: image.Pt(centerX, centerY),
		Image:    itemImage,
		Width:    itemImage.Bounds().Dx(),
		Height:   itemImage.Bounds().Dy(),
	}
}

var catFoodImage *ebiten.Image // Declare a variable to hold the cat food image

func init() {
	// Load the cat food image from a file
	var err error
	catFoodImage, _, err = ebitenutil.NewImageFromFile("assets/catfood.png")
	if err != nil {
		log.Fatalf("Failed to load cat food image: %v", err)
	}
}

func (p *Projectile) Draw(screen *ebiten.Image) {

	// Assuming you have a sprite or image for the projectile
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(catFoodImage, opts) // Replace 'projectileImage' with your image variable

}

func (p *Projectile) Update() {
	// Example: Move the projectile
	p.X += p.VelocityX
	p.Y += p.VelocityY

	// Add any other update logic here (like collision detection)
}
func NewProjectile(x, y, velocityX, velocityY float64) Projectile {
	return Projectile{
		X:         x,
		Y:         y,
		VelocityX: velocityX,
		VelocityY: velocityY,
		// ... initialize other fields
	}
}
