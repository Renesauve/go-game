package game

import (
	"go-game/config"
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
    
}




func NewGame() *Game {
    gridMiddleX, gridMiddleY := 2, 2 

    // Calculate screen position from grid position
    screenMiddleX := float64(gridMiddleX) * (float64(config.ScreenWidth) / float64(config.GridSize)) + (float64(config.ScreenWidth) / float64(config.GridSize) / 2.0) - float64(config.PlayerWidth)/2.0
    screenMiddleY := float64(gridMiddleY) * (float64(config.ScreenHeight) / float64(config.GridSize)) + (float64(config.ScreenHeight) / float64(config.GridSize) / 2.0) - float64(config.PlayerHeight)/2.0

    p := player.NewPlayer(screenMiddleX, screenMiddleY, [2]int{config.GridSize / 2, config.GridSize / 2})
    roomGrid := [config.GridSize][config.GridSize]*room.Room{}

  
    return &Game{
        Player: p,
        Rooms:  make(map[string]*room.Room),
        RoomsVisited: [config.GridSize][config.GridSize]bool{},
        RoomGrid: roomGrid,
        State: MainMenu, // Start at the main menu
    }
   
}


func (g *Game) GenerateRooms() {
   
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
}

func (g *Game) Update() error {
    
    if g.CurrentRoom == nil {
     
        g.GenerateRooms()
    }
    
    proposedX, proposedY := g.Player.X, g.Player.Y
    


    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        proposedX += config.Speed
        if proposedX > config.ScreenWidth - float64(config.PlayerWidth) {
            if g.Player.Coordinates[0] < config.GridSize-1 {
                g.CurrentRoom = g.getNextRoom("right")
                proposedX = 0 // Reset to left edge of the new room
            } else {
                proposedX = config.ScreenWidth - float64(config.PlayerWidth) // Prevent leaving the grid
            }
        }
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        proposedX -= config.Speed
        if proposedX < 0 { // Check for left boundary
            if g.Player.Coordinates[0] > 0 {
                g.CurrentRoom = g.getNextRoom("left")
         
                proposedX = config.ScreenWidth - float64(config.PlayerWidth) // Start at the right edge of the new room
            } else {
                proposedX = 0 // Stay within current room
            }
        }
    }
    
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        proposedY -= config.Speed
        if proposedY < 0 { // Check for top boundary
            if g.Player.Coordinates[1] > 0 {
                g.CurrentRoom = g.getNextRoom("up")
           
                proposedY = config.ScreenHeight - float64(config.PlayerHeight) // Start at the bottom of the new room
            } else {
                proposedY = 0 // Stay within current room
            }
        }
    }
    
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        proposedY += config.Speed
        if proposedY > config.ScreenHeight - float64(config.PlayerHeight) { // Check for bottom boundary
            if g.Player.Coordinates[1] < config.GridSize-1 {
                g.CurrentRoom = g.getNextRoom("down")
              
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


    return nil
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    // For simplicity, let's return the outside size; you can adjust as needed for your game
    return outsideWidth, outsideHeight
}


func (g *Game) Draw(screen *ebiten.Image) {

    if g.CurrentRoom != nil && g.CurrentRoom.RoomType == room.BossRoom && g.BossImage != nil {
        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(float64(g.BambiPosition.X), float64(g.BambiPosition.Y))
        screen.DrawImage(g.BossImage, opts)
    }

    if g.CurrentRoom != nil && g.CurrentRoom.RoomType == room.ItemRoom && g.ItemImage != nil {
        opts := &ebiten.DrawImageOptions{}
  
        screen.DrawImage(g.ItemImage, opts)
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




func (g *Game) getNextRoom(direction string) *room.Room {
    


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

    // Check if the next room is a BossRoom and initialize it
    if nextRoom != nil && nextRoom.RoomType == room.BossRoom && g.CurrentRoom != nextRoom {
        g.CurrentRoom = nextRoom
        g.initializeBossRoom() // Initialize the boss room only once
    }
    if nextRoom != nil && nextRoom.RoomType == room.ItemRoom && g.CurrentRoom != nextRoom {
        g.CurrentRoom = nextRoom
        g.initializeItemRoom() // Initialize the boss room only once
    }

    

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
    if g.ItemImage == nil {
        var err error
        g.ItemImage, _, err = ebitenutil.NewImageFromFile(`assets/catfood.png`)
        if err != nil {
            log.Fatal(err)
        }
    }

    centerX, centerY := config.ScreenWidth/2, config.ScreenHeight/2
    g.BambiPosition = image.Point{X: centerX - g.ItemImage.Bounds().Dx()/2, Y: centerY - g.ItemImage.Bounds().Dy()/2}



    bambiObstacle := room.Obstacle{
        Rect: image.Rect(
            g.BambiPosition.X,
            g.BambiPosition.Y,
            g.BambiPosition.X + g.ItemImage.Bounds().Dx(),
            g.BambiPosition.Y + g.ItemImage.Bounds().Dy(),
        ),
      
    }

    g.CurrentRoom.Obstacles = append(g.CurrentRoom.Obstacles, bambiObstacle)
}