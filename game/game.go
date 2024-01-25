package game

import (
	"go-game/config"
	"go-game/player"
	"go-game/room"
	"go-game/util"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	// other imports as necessary
)





type Game struct {
	Player       player.Player
	CurrentRoom  *room.Room
	RoomsVisited [config.GridSize][config.GridSize]bool
	Rooms        map[string]*room.Room // A map to hold the rooms
	RoomGrid     [config.GridSize][config.GridSize]*room.Room
}




func NewGame() *Game {
    // Initialization logic
		// run generateRooms function
		
	

		p := player.NewPlayer(config.ScreenWidth/2-config.PlayerWidth/2, config.ScreenHeight/2-config.PlayerHeight/2)
		roomGrid := [config.GridSize][config.GridSize]*room.Room{}
	
    return &Game{
        Player: p,
		Rooms:  make(map[string]*room.Room),
		RoomsVisited: [config.GridSize][config.GridSize]bool{},
		RoomGrid: [config.GridSize][config.GridSize]*room.Room{},
		CurrentRoom: roomGrid[config.GridSize/2][config.GridSize/2],

        // ... other initialization
    }
	
}


func (g *Game) Update() error {
	
  

	if g.CurrentRoom == nil {
	
        g.CurrentRoom = g.getNextRoom("initial")
    }
	proposedX, proposedY := g.Player.X, g.Player.Y

    speed := 8.0 // Speed of the player

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		proposedX += speed
		if proposedX > config.ScreenWidth - config.PlayerWidth {
			g.CurrentRoom = g.getNextRoom("right")
			proposedX = 0
			g.Player.X = float64(config.ScreenHeight) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		proposedX -= speed
		if proposedX < 0 {
			g.CurrentRoom = g.getNextRoom("left")
			proposedX = config.ScreenWidth - config.PlayerWidth
			g.Player.Y = float64(config.ScreenHeight) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		proposedY -= speed
		if proposedY < 0 {
			g.CurrentRoom = g.getNextRoom("up")
			proposedY = config.ScreenHeight - config.PlayerHeight
			g.Player.X = float64(config.ScreenWidth) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		proposedY += speed
		if proposedY > config.ScreenHeight - config.PlayerHeight {
			g.CurrentRoom = g.getNextRoom("down")
			proposedY = 0
			g.Player.X = float64(config.ScreenWidth) / 2 // Align with the door
		}
	}
  
	g.RoomsVisited[g.Player.Coordinates[0]][g.Player.Coordinates[1]] = true
    // Calculate proposed new position and check for collisions...

    playerRect := image.Rect(int(proposedX), int(proposedY), int(proposedX)+config.PlayerWidth, int(proposedY)+config.PlayerHeight)

    collision := false
	if(g.CurrentRoom == nil) {
		return nil
	}
    for _, obstacle := range g.CurrentRoom.Obstacles {
        if !obstacle.IsDoor && playerRect.Overlaps(obstacle.Rect) {
            collision = true
            break
        }
    }

    if !collision {
        g.Player.X, g.Player.Y = proposedX, proposedY // Update position if no collision
    }

    return nil
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

    if g.RoomGrid[nextX][nextY] == nil {
        room.GenerateRoom(nextX, nextY,  direction, room.RegularRoom, g.RoomGrid) // Generate with door on correct side
    }

    g.Player.Coordinates[0], g.Player.Coordinates[1] = nextX, nextY
    return g.RoomGrid[nextX][nextY]
}



// func (g *Game) ensureDoorAlignment(x, y int, direction string) {
//     // room := g.RoomGrid[x][y]
// 	// fmt.Println(room)
//     oppositeDir := util.OppositeDirection(direction)

//     // Check for an existing door in the opposite direction
//     hasOppositeDoor := false
//     for _, obstacle := range room.Obstacles {
//         if obstacle.IsDoor {
//             switch oppositeDir {
//             case "top":
//                 if obstacle.Rect.Min.Y == 0 {
//                     hasOppositeDoor = true
//                 }
//             case "bottom":
//                 if obstacle.Rect.Max.Y == config.ScreenHeight {
//                     hasOppositeDoor = true
//                 }
//             case "left":
//                 if obstacle.Rect.Min.X == 0 {
//                     hasOppositeDoor = true
//                 }
//             case "right":
//                 if obstacle.Rect.Max.X == config.ScreenWidth {
//                     hasOppositeDoor = true
//                 }
//             }
//         }
//     }

	
//     if !hasOppositeDoor {
//         // Add a door in the opposite direction
//         doorSize := 30
//         var doorRect image.Rectangle
//         switch oppositeDir {
//         case "top":
//             doorRect = image.Rect(config.ScreenWidth/2-doorSize/2, 0, config.ScreenWidth/2+doorSize/2, 30)
//         case "bottom":
//             doorRect = image.Rect(config.ScreenWidth/2-doorSize/2, config.ScreenHeight-30, config.ScreenWidth/2+doorSize/2, config.ScreenHeight)
//         case "left":
//             doorRect = image.Rect(0, config.ScreenHeight/2-doorSize/2, 30, config.ScreenHeight/2+doorSize/2)
//         case "right":
//             doorRect = image.Rect(config.ScreenWidth-30, config.ScreenHeight/2-doorSize/2, config.ScreenWidth, config.ScreenHeight/2+doorSize/2)
//         }

//         room.Obstacles = append(room.Obstacles, Obstacle{Rect: doorRect, IsDoor: true})
//         g.RoomGrid[x][y] = room // Update the room in the grid
//     }
// }


func (g *Game) determineDirectionForNewRoom(x, y int) string {
    // Check the adjacent rooms and return the direction of the first found room
    if x > 0 && g.RoomGrid[x-1][y] != nil {
        return "right" // Room to the left, so new room is accessed from right
    }
    if x < config.GridSize-1 && g.RoomGrid[x+1][y] != nil {
        return "left" // Room to the right, so new room is accessed from left
    }
    if y > 0 && g.RoomGrid[x][y-1] != nil {
        return "down" // Room above, so new room is accessed from below
    }
    if y < config.GridSize-1 && g.RoomGrid[x][y+1] != nil {
        return "up" // Room below, so new room is accessed from above
    }
    
    return "" // No adjacent rooms, return an empty string
}





func getRandomSpecialRoomPositions(gridSize int) ([2]int, [2]int) {
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

func (g *Game) GenerateRooms() {
	
    startX, startY := config.GridSize / 2, config.GridSize / 2
    g.RoomGrid[startX][startY] = room.GenerateRoom(startX, startY, "initial", room.RegularRoom, g.RoomGrid)

    bossRoomPos, itemRoomPos := getRandomSpecialRoomPositions(config.GridSize)

    for x := 0; x < config.GridSize; x++ {
        for y := 0; y < config.GridSize; y++ {
            if g.RoomGrid[x][y] == nil {
                var roomType room.RoomType = room.RegularRoom
                if x == bossRoomPos[0] && y == bossRoomPos[1] {
                    roomType = room.BossRoom
                } else if x == itemRoomPos[0] && y == itemRoomPos[1] {
                    roomType = room.ItemRoom
                }

                direction := g.determineDirectionForNewRoom(x, y)
                g.RoomGrid[x][y] = room.GenerateRoom(x, y, direction, roomType, g.RoomGrid)
            }
        }
    }
}




func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    // For simplicity, let's return the outside size; you can adjust as needed for your game
    return outsideWidth, outsideHeight
}


func (g *Game) Draw(screen *ebiten.Image) {
	centerX := float32(config.ScreenWidth) / 2
    centerY := float32(config.ScreenHeight) / 2
    const objectSize = 30 // Size of the object (item or boss)

	
   
    
		switch g.CurrentRoom.RoomType {
		case room.BossRoom:
			// Draw a red circle for Boss Room
			bossCircle := util.CreateCircleImage(objectSize/2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(centerX)-float64(objectSize)/2, float64(centerY)-float64(objectSize)/2)
			screen.DrawImage(bossCircle, opts)

		case room.ItemRoom:
			// Draw a yellow circle for Item Room
			itemCircle := util.CreateCircleImage(objectSize/2, color.RGBA{R: 255, G: 255, B: 0, A: 255})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(centerX)-float64(objectSize)/2, float64(centerY)-float64(objectSize)/2)
			screen.DrawImage(itemCircle, opts)
		}

    // Draw the player
	vector.DrawFilledRect(screen, float32(g.Player.X), float32(g.Player.Y), config.PlayerWidth, config.PlayerHeight, color.White, false)

	    // Draw obstacles
		for _, obstacle := range g.CurrentRoom.Obstacles {
			vector.DrawFilledRect(screen, float32(obstacle.Rect.Min.X), float32(obstacle.Rect.Min.Y), float32(obstacle.Rect.Dx()), float32(obstacle.Rect.Dy()), color.Gray{Y: 0x80}, false)
		}

		for x := 0; x < config.GridSize; x++ {
			for y := 0; y < config.GridSize; y++ {
				var roomColor color.Color
	
				if g.RoomsVisited[x][y] {
					roomColor = color.Gray{Y: 150} // Visited room color
				} else {
					roomColor = color.Gray{Y: 50} // Unvisited room color
				}
	
				if x == g.Player.Coordinates[0] && y == g.Player.Coordinates[1] {
					roomColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Current room color
				}
	
				minimapX := config.ScreenWidth - config.GridSize*config.MinimapRoomSize - config.MinimapRoomSize + x*config.MinimapRoomSize
				minimapY := config.ScreenHeight - config.GridSize*config.MinimapRoomSize - config.MinimapRoomSize + y*config.MinimapRoomSize
	
				ebitenutil.DrawRect(screen, float64(minimapX), float64(minimapY), float64(config.MinimapRoomSize), float64(config.MinimapRoomSize), roomColor)
			}
		}
}
