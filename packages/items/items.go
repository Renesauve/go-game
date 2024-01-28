package items

import (
	"encoding/xml"
	_ "image/png"
)

type Itemizable interface {
	GetID() int
	GetName() string
	GetDescription() string
	GetX() float64
	GetY() float64
	GetGFX() string
	// Other common methods for items
}

type Item struct {
	ID          int     `xml:"id,attr"`
	Name        string  `xml:"name,attr"`
	Description string  `xml:"description,attr"`
	GFX         string  `xml:"gfx,attr"`
	X, Y        float64 // Position of the item
}

func (i Item) GetID() int             { return i.ID }
func (i Item) GetName() string        { return i.Name }
func (i Item) GetDescription() string { return i.Description }
func (i Item) GetGFX() string         { return i.GFX }
func (i Item) GetX() float64          { return i.X }
func (i Item) GetY() float64          { return i.Y }

// Items struct to hold slices of Weapon and Armor

// Weapon struct extends Item and includes specific attributes for weapons
type Weapon struct {
	Item              // Embedding Item struct
	Damage     int    `xml:"damage,attr"`
	Speed      int    `xml:"speed,attr"`
	Handedness string // "one-handed" or "two-handed"
	Range      int    // Attack range of the weapon
}

// Armor represents an armor in the game.
type Armor struct {
	Item        // Embedding Item struct
	Defense int `xml:"defense,attr"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Weapons []Weapon `xml:"weapon"`
	Armors  []Armor  `xml:"armor"`
}

type Inventory struct {
	Items []Itemizable
}

func (inv *Inventory) AddItem(item Itemizable) {
	inv.Items = append(inv.Items, item)
}

// RemoveItem removes an item from the inventory by ID.
func (inv *Inventory) RemoveItem(itemID int) {
	for i, item := range inv.Items {
		if item.GetID() == itemID {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			break
		}
	}
}
