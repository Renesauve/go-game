package player

type Player struct {
    X, Y         float64
    Inventory    map[string]int
    Coordinates  [2]int
    // ... other fields
}

func NewPlayer(x, y float64, coordinates [2]int) Player {
    return Player{
        X: x,
        Y: y,
        Inventory: make(map[string]int),
        Coordinates: [2]int{coordinates[0], coordinates[1]},
    }
}

// Player methods