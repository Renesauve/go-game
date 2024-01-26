package game

import (
	"fmt"
	"go-game/config"
	"go-game/items"
	"go-game/player"

	"go-game/room"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	// other imports as necessary
)



type GameState int

const (
    MainMenu GameState = iota
    InGame
    GameOver
)

type Game struct {
	Player       player.Player
	CurrentRoom  *room.Room
	RoomsVisited [config.GridSize][config.GridSize]bool
	Rooms        map[string]*room.Room // A map to hold the rooms
	RoomGrid     [config.GridSize][config.GridSize]*room.Room
    BossImage    *ebiten.Image  // Field to store the Bambi image
    ItemImage   *ebiten.Image  // Field to store the item image
    BambiPosition image.Point // Initialize to (0,0)
    BambiVelocity image.Point 
    State GameState // Start at the main menu
    CollectedItems map[string]bool
    LastKeyState map[ebiten.Key]bool
    
}




func NewGame() *Game {
    gridMiddleX, gridMiddleY := 2, 2 

    // Calculate screen position from grid position
    screenMiddleX := float64(gridMiddleX) * (float64(config.ScreenWidth) / float64(config.GridSize)) + (float64(config.ScreenWidth) / float64(config.GridSize) / 2.0) - float64(config.PlayerWidth)/2.0
    screenMiddleY := float64(gridMiddleY) * (float64(config.ScreenHeight) / float64(config.GridSize)) + (float64(config.ScreenHeight) / float64(config.GridSize) / 2.0) - float64(config.PlayerHeight)/2.0

    p := player.NewPlayer(screenMiddleX, screenMiddleY, [2]int{config.GridSize / 2, config.GridSize / 2})
    roomGrid := [config.GridSize][config.GridSize]*room.Room{}

  
    g := &Game{
        Player:       p,
        Rooms:        make(map[string]*room.Room),
        RoomsVisited: [config.GridSize][config.GridSize]bool{},
        RoomGrid:     roomGrid,
        State:        MainMenu, // Start at the main menu
        CollectedItems: make(map[string]bool),
        LastKeyState: make(map[ebiten.Key]bool),
        // Initialize other fields...
    }
    bossRoomPos, itemRoomPos := g.GenerateRooms() // This function now needs to return the boss and item room positions

    // Access the boss and item rooms using the positions provided by GenerateRooms
    bossRoom := g.RoomGrid[bossRoomPos[0]][bossRoomPos[1]]
    itemRoom := g.RoomGrid[itemRoomPos[0]][itemRoomPos[1]]
    // Save the current room
    currentRoom := g.CurrentRoom
    // Temporarily set the current room to the boss room to initialize it
    g.CurrentRoom = bossRoom
    g.initializeBossRoom()
    // Temporarily set the current room to the item room to initialize it
    g.CurrentRoom = itemRoom
    g.initializeItemRoom()
    // Restore the current room
    g.CurrentRoom = currentRoom

    return g
    
   
}


func (g *Game) GenerateRooms() ([2]int, [2]int) {
   
    startX, startY := config.GridSize / 2, config.GridSize / 2
    room.GenerateRoom(startX, startY, room.RegularRoom, g.RoomGrid)
   
    bossRoomPos, itemRoomPos := getRandomSpecialRoomPositions(config.GridSize)
    
    for x := 0; x < config.GridSize; x++ {
        for y := 0; y < config.GridSize; y++ {
            if g.RoomGrid[x][y] == nil {
                var roomType room.RoomType = room.RegularRoom
                if x == bossRoomPos[0] && y == bossRoomPos[1] {
                    roomType = room.BossRoom
                } else if x == itemRoomPos[0] && y == itemRoomPos[1] {
                    roomType = room.ItemRoom // Set the room type to ItemRoom
                }

              
                g.RoomGrid[x][y] = room.GenerateRoom(x, y, roomType,  g.RoomGrid)
            }
        }
    }
    return bossRoomPos, itemRoomPos
}

func (g *Game) Update() error {
    if g.CurrentRoom == nil {
        g.GenerateRooms()
    }
    for i := range items.Projectiles {
        items.Projectiles[i].Update()
    }
  
    
    
    proposedX, proposedY := g.Player.X, g.Player.Y
    
    g.PickupItems()

    if ebiten.IsKeyPressed(ebiten.KeyX) {
    
        g.Player.ThrowCatFood()
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        proposedX += config.Speed
        g.Player.Facing = player.DirectionRight
        if proposedX > config.ScreenWidth - float64(config.PlayerWidth) {
            if g.Player.Coordinates[0] < config.GridSize-1 {
                g.CurrentRoom = g.GetRoomInDirection("right")
                proposedX = 0 // Reset to left edge of the new room
            } else {
                proposedX = config.ScreenWidth - float64(config.PlayerWidth) // Prevent leaving the grid
            }
        }
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        proposedX -= config.Speed
        g.Player.Facing = player.DirectionLeft
        if proposedX < 0 { // Check for left boundary
            if g.Player.Coordinates[0] > 0 {
                g.CurrentRoom = g.GetRoomInDirection("left")
         
                proposedX = config.ScreenWidth - float64(config.PlayerWidth) // Start at the right edge of the new room
            } else {
                proposedX = 0 // Stay within current room
            }
        }
    }
    
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        proposedY -= config.Speed
        g.Player.Facing = player.DirectionUp
        if proposedY < 0 { // Check for top boundary
            if g.Player.Coordinates[1] > 0 {
                g.CurrentRoom = g.GetRoomInDirection("up")
                
                proposedY = config.ScreenHeight - float64(config.PlayerHeight) // Start at the bottom of the new room
            } else {
                proposedY = 0 // Stay within current room
            }
        }
    }
    
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        proposedY += config.Speed
        g.Player.Facing = player.DirectionDown
        if proposedY > config.ScreenHeight - float64(config.PlayerHeight) { // Check for bottom boundary
            if g.Player.Coordinates[1] < config.GridSize-1 {
                g.CurrentRoom = g.GetRoomInDirection("down")
              
                proposedY = 0 // Start at the top of the new room
            } else {
                proposedY = config.ScreenHeight - float64(config.PlayerHeight) // Prevent leaving the grid
            }
        }
    }

    // Update the player's position
    g.Player.X, g.Player.Y = proposedX, proposedY
  
	g.RoomsVisited[g.Player.Coordinates[0]][g.Player.Coordinates[1]] = true
    // Calculate proposed new position and check for collisions...


    collision := false

    if !collision {
        g.Player.X, g.Player.Y = proposedX, proposedY // Update position if no collision
    } 

    if g.CurrentRoom != nil && g.CurrentRoom.RoomType == room.BossRoom {
        // Randomize Bambi's movement
        if rand.Intn(10) == 0 { // This will only change direction approximately every 10 frames
            g.BambiVelocity = image.Point{X: rand.Intn(5) - 1, Y: rand.Intn(5) - 1}
        }

        // Update Bambi's position
        g.BambiPosition = g.BambiPosition.Add(g.BambiVelocity)

        // Keep Bambi within room bounds
        if g.BambiPosition.X < 0 || g.BambiPosition.X+g.BossImage.Bounds().Dx() > config.ScreenWidth {
            g.BambiVelocity.X = -g.BambiVelocity.X
            g.BambiPosition.X += g.BambiVelocity.X * 2 // Adjust position after reversing velocity
        }
        if g.BambiPosition.Y < 0 || g.BambiPosition.Y+g.BossImage.Bounds().Dy() > config.ScreenHeight {
            g.BambiVelocity.Y = -g.BambiVelocity.Y
            g.BambiPosition.Y += g.BambiVelocity.Y * 2 // Adjust position after reversing velocity
        }

           // Update Bambi's Rect to the new position
           if len(g.CurrentRoom.Obstacles) > 0 {
            g.CurrentRoom.Obstacles[0].Rect = image.Rect(
                g.BambiPosition.X,
                g.BambiPosition.Y,
                g.BambiPosition.X + g.BossImage.Bounds().Dx(),
                g.BambiPosition.Y + g.BossImage.Bounds().Dy(),
            )
        }

        // Check for collision between Bambi and the player
        playerRect := image.Rect(
            int(g.Player.X),
            int(g.Player.Y),
            int(g.Player.X) + config.PlayerWidth,
            int(g.Player.Y) + config.PlayerHeight,
        )

        if g.CurrentRoom.Obstacles[0].Rect.Overlaps(playerRect) {
            *g = *NewGame() // Restart the game
            return nil
        }
    }

    g.LastKeyState[ebiten.KeyX] = ebiten.IsKeyPressed(ebiten.KeyX)

    return nil
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    // For simplicity, let's return the outside size; you can adjust as needed for your game
    return outsideWidth, outsideHeight
}


func (g *Game) Draw(screen *ebiten.Image) {
    for _, projectile := range items.Projectiles {
        projectile.Draw(screen)
    }
    if g.CurrentRoom != nil && g.CurrentRoom.RoomType == room.BossRoom && g.BossImage != nil {
        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(float64(g.BambiPosition.X), float64(g.BambiPosition.Y))
        screen.DrawImage(g.BossImage, opts)
    }
    if g.CurrentRoom != nil {
        for _, item := range g.CurrentRoom.Items {
            if !item.Collected { // Only draw the item if it hasn't been collected
                opts := &ebiten.DrawImageOptions{}
                opts.GeoM.Translate(float64(item.Position.X), float64(item.Position.Y))
                screen.DrawImage(item.Image, opts)
            }
        }
    }

    for i, item := range g.Player.Inventory {
        opts := &ebiten.DrawImageOptions{}

        // Scale the image to the defined inventory item size
        scale := float64(config.InventoryItemSize) / float64(item.Image.Bounds().Dx())
        opts.GeoM.Scale(scale, scale)

        // Calculate the position for each item
        xPos := float64(i * (config.InventoryItemSize + 5)) // 5 pixels space between items
        yPos := float64(config.ScreenHeight - config.InventoryItemSize) // Positioned at bottom

        opts.GeoM.Translate(xPos, yPos)
        screen.DrawImage(item.Image, opts)
    }
    

    minimapStartX := config.ScreenWidth - config.GridSize*config.MinimapRoomSize - config.MinimapPadding
    minimapStartY := config.ScreenHeight - config.GridSize*config.MinimapRoomSize - config.MinimapPadding

    for x := 0; x < config.GridSize; x++ {
        for y := 0; y < config.GridSize; y++ {
            var roomColor color.Color

            // Determine room color
            if g.RoomsVisited[x][y] {
                roomColor = color.Gray{Y: 150} // Visited room color
            } else {
                roomColor = color.Gray{Y: 50} // Unvisited room color
            }

            // Highlight the current room
            if x == g.Player.Coordinates[0] && y == g.Player.Coordinates[1] {
                roomColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Current room color
            }

            // Calculate the position of each room in the minimap
            minimapX := minimapStartX + x*config.MinimapRoomSize
            minimapY := minimapStartY + y*config.MinimapRoomSize

            // Draw the room on the minimap
            vector.DrawFilledRect(screen, float32(minimapX), float32(minimapY), float32(config.MinimapRoomSize), float32(config.MinimapRoomSize), roomColor, true)
        }
    }

    vector.DrawFilledRect(screen, float32(g.Player.X), float32(g.Player.Y), config.PlayerWidth, config.PlayerHeight, color.White, false)


}


func getRandomSpecialRoomPositions(gridSize int, ) ([2]int, [2]int) {
    var edgePositions, otherPositions [][2]int
    centerX, centerY := gridSize/2, gridSize/2

    // Separate edge positions and other positions
    for x := 0; x < gridSize; x++ {
        for y := 0; y < gridSize; y++ {
            if x == 0 || x == gridSize-1 || y == 0 || y == gridSize-1 {
                edgePositions = append(edgePositions, [2]int{x, y})
            } else if x != centerX || y != centerY {
                otherPositions = append(otherPositions, [2]int{x, y})
            }
        }
    }

    // Shuffle the slices
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(edgePositions), func(i, j int) {
        edgePositions[i], edgePositions[j] = edgePositions[j], edgePositions[i]
    })
    rand.Shuffle(len(otherPositions), func(i, j int) {
        otherPositions[i], otherPositions[j] = otherPositions[j], otherPositions[i]
    })

    // Select positions for special rooms
    var bossRoomPos, itemRoomPos [2]int
    if len(edgePositions) > 0 {
        bossRoomPos = edgePositions[0]
    }
    if len(otherPositions) > 0 {
        itemRoomPos = otherPositions[0]
    }

    return bossRoomPos, itemRoomPos
}




func (g *Game) GetRoomInDirection(direction string) *room.Room {
    nextX, nextY := g.Player.Coordinates[0], g.Player.Coordinates[1]
    switch direction {
    case "up":
        nextY--
    case "down":
        nextY++
    case "left":
        nextX--
    case "right":
        nextX++
    }
    if nextX < 0 || nextX >= config.GridSize || nextY < 0 || nextY >= config.GridSize {
        return g.CurrentRoom // Prevent leaving the grid
    }
    g.Player.Coordinates[0], g.Player.Coordinates[1] = nextX, nextY
    nextRoom := g.RoomGrid[nextX][nextY]
    return nextRoom
}

func (g *Game) initializeBossRoom() {
    if g.BossImage == nil {
        var err error
        g.BossImage, _, err = ebitenutil.NewImageFromFile(`assets/bambi.jpg`)
        if err != nil {
            log.Fatal(err)
        }
    }

    centerX, centerY := config.ScreenWidth/2, config.ScreenHeight/2
    g.BambiPosition = image.Point{X: centerX - g.BossImage.Bounds().Dx()/2, Y: centerY - g.BossImage.Bounds().Dy()/2}
    g.BambiVelocity = image.Point{X: 1, Y: 1} // Example initial velocity

    bambiObstacle := room.Obstacle{
        Rect: image.Rect(
            g.BambiPosition.X,
            g.BambiPosition.Y,
            g.BambiPosition.X + g.BossImage.Bounds().Dx(),
            g.BambiPosition.Y + g.BossImage.Bounds().Dy(),
        ),
      
    }

    g.CurrentRoom.Obstacles = append(g.CurrentRoom.Obstacles, bambiObstacle)
}

func (g *Game) initializeItemRoom() {
    if !g.CollectedItems["CatFood"] { // Check if CatFood has been collected
        item := items.InitializeItem(items.CatFood, "CatFood", `assets/catfood.png`, config.ScreenWidth, config.ScreenHeight, true)
        g.CurrentRoom.Items = append(g.CurrentRoom.Items, item)
     
    }
}

func (g *Game) PickupItems() {
    if g.CurrentRoom != nil {
        for i := 0; i < len(g.CurrentRoom.Items); i++ {
            // Check if the item is not collected and there's a collision with the player
            if !g.CurrentRoom.Items[i].Collected && g.Player.CheckCollisionWithItem(&g.CurrentRoom.Items[i]) {
                // Mark the item as collected
                g.CurrentRoom.Items[i].Collected = true

                // Add the item to the player's inventory if it's not already added
          
                if _, exists := g.CollectedItems[g.CurrentRoom.Items[i].Name]; !exists {
                    g.Player.Inventory = append(g.Player.Inventory, g.CurrentRoom.Items[i])
                    g.Player.Inventory[i].IsShootable = true
                    g.CollectedItems[g.CurrentRoom.Items[i].Name] = true
                    fmt.Println(g.CurrentRoom.Items[i].IsShootable)
                }
            }
        }
    }
}


