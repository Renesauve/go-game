package player

type Player struct {
    X, Y         float64
    Inventory    map[string]int
    Coordinates  [2]int
    // ... other fields
}

func NewPlayer(x, y float64) Player {
    return Player{
        X: x,
        Y: y,
        Inventory: make(map[string]int),


        // ... other initialization
    }
}

// Player methods