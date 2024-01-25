package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
    x, y float64
	inventory map[string]int // Inventory with item name and quantity
	coordinates [2]int // New field to track player's room coordinates
}


type Game struct {
    player  Player
    currentRoom *Room
    rooms   map[string]*Room // A map to hold the rooms
	roomsVisited [gridSize][gridSize]bool
}

type RoomType int

const (
    RegularRoom RoomType = iota
    ItemRoom
    BossRoom
)
type Obstacle struct {
    Rect  image.Rectangle
    IsDoor bool
}
type Room struct {
    obstacles  []Obstacle
    roomType   RoomType
    // other room details
}

const (
	screenWidth  = 1920
	screenHeight = 1080
	playerWidth  = 24
	playerHeight = 24
	MaxRooms = 7
	minimapRoomSize = 30 // Size of each room in the minimap
	minimapPadding  = 15  // Padding around the minimap
	
	
)


var (
    itemRoomGenerated = false
    bossRoomGenerated = false
)

const gridSize = 5 // Assuming a 5x5 grid
var roomGrid [gridSize][gridSize]*Room

func (g *Game) Update() error {
	if g.currentRoom == nil {
        g.currentRoom = g.getNextRoom("initial")
    }
	proposedX, proposedY := g.player.x, g.player.y

    speed := 8.0 // Speed of the player

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		proposedX += speed
		if proposedX > screenWidth - playerWidth {
			g.currentRoom = g.getNextRoom("right")
			proposedX = 0
			g.player.y = float64(screenHeight) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		proposedX -= speed
		if proposedX < 0 {
			g.currentRoom = g.getNextRoom("left")
			proposedX = screenWidth - playerWidth
			g.player.y = float64(screenHeight) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		proposedY -= speed
		if proposedY < 0 {
			g.currentRoom = g.getNextRoom("up")
			proposedY = screenHeight - playerHeight
			g.player.x = float64(screenWidth) / 2 // Align with the door
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		proposedY += speed
		if proposedY > screenHeight - playerHeight {
			g.currentRoom = g.getNextRoom("down")
			proposedY = 0
			g.player.x = float64(screenWidth) / 2 // Align with the door
		}
	}
  
	g.roomsVisited[g.player.coordinates[0]][g.player.coordinates[1]] = true
    // Calculate proposed new position and check for collisions...

    playerRect := image.Rect(int(proposedX), int(proposedY), int(proposedX)+playerWidth, int(proposedY)+playerHeight)

    collision := false
    for _, obstacle := range g.currentRoom.obstacles {
        if !obstacle.IsDoor && playerRect.Overlaps(obstacle.Rect) {
            collision = true
            break
        }
    }

    if !collision {
        g.player.x, g.player.y = proposedX, proposedY // Update position if no collision
    }

    return nil
}



func generateObstacles(withDoors map[string]bool) []Obstacle {
    var obstacles []Obstacle
    doorSize := 60

    // Create obstacles for each side of the room
    for _, side := range []string{"top", "bottom", "left", "right"} {
        isDoor := withDoors[side]
        switch side {
        case "top", "bottom":
            for x := 0; x < screenWidth; x += 30 {
                if side == "top" && (!isDoor || x < screenWidth/2-doorSize/2 || x > screenWidth/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(x, 0, x+30, 30), IsDoor: false})
                }
                if side == "bottom" && (!isDoor || x < screenWidth/2-doorSize/2 || x > screenWidth/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(x, screenHeight-30, x+30, screenHeight), IsDoor: false})
                }
            }
        case "left", "right":
            for y := 30; y < screenHeight-30; y += 30 {
                if side == "left" && (!isDoor || y < screenHeight/2-doorSize/2 || y > screenHeight/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(0, y, 30, y+30), IsDoor: false})
                }
                if side == "right" && (!isDoor || y < screenHeight/2-doorSize/2 || y > screenHeight/2+doorSize/2) {
                    obstacles = append(obstacles, Obstacle{Rect: image.Rect(screenWidth-30, y, screenWidth, y+30), IsDoor: false})
                }
            }
        }
    }

    return obstacles
}

func oppositeDirection(direction string) string {
    switch direction {
    case "up":
        return "bottom"
    case "down":
        return "top"
    case "left":
        return "right"
    case "right":
        return "left"
    }
    return ""
}


func generateRooms() {
    startX, startY := gridSize / 2, gridSize / 2
    generateRoom(startX, startY, nil, "initial", RegularRoom)

    bossRoomPos, itemRoomPos := getRandomSpecialRoomPositions(gridSize)

    for x := 0; x < gridSize; x++ {
        for y := 0; y < gridSize; y++ {
            if roomGrid[x][y] == nil {
                var roomType RoomType = RegularRoom
                if x == bossRoomPos[0] && y == bossRoomPos[1] {
                    roomType = BossRoom
                } else if x == itemRoomPos[0] && y == itemRoomPos[1] {
                    roomType = ItemRoom
                }

                direction := determineDirectionForNewRoom(x, y)
                generateRoom(x, y, nil, direction, roomType)
            }
        }
    }
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

    // Shuffle and pick one position from edge positions for Boss Room
    rand.Shuffle(len(edgePositions), func(i, j int) {
        edgePositions[i], edgePositions[j] = edgePositions[j], edgePositions[i]
    })
    bossRoomPos := edgePositions[0]

    // Shuffle and pick one position from other positions for Item Room
    rand.Shuffle(len(otherPositions), func(i, j int) {
        otherPositions[i], otherPositions[j] = otherPositions[j], otherPositions[i]
    })
    itemRoomPos := otherPositions[0]

    return bossRoomPos, itemRoomPos
}

func determineDirectionForNewRoom(x, y int) string {
    // Check the adjacent rooms and return the direction of the first found room
    if x > 0 && roomGrid[x-1][y] != nil {
        return "right" // Room to the left, so new room is accessed from right
    }
    if x < gridSize-1 && roomGrid[x+1][y] != nil {
        return "left" // Room to the right, so new room is accessed from left
    }
    if y > 0 && roomGrid[x][y-1] != nil {
        return "down" // Room above, so new room is accessed from below
    }
    if y < gridSize-1 && roomGrid[x][y+1] != nil {
        return "up" // Room below, so new room is accessed from above
    }
    
    return "" // No adjacent rooms, return an empty string
}


func alignDoorsWithAdjacentRooms(x, y int, withDoors map[string]bool) {
    if y > 0 && roomGrid[x][y-1] != nil { // Top
        withDoors["top"] = true
    }
    if y < gridSize-1 && roomGrid[x][y+1] != nil { // Bottom
        withDoors["bottom"] = true
    }
    if x > 0 && roomGrid[x-1][y] != nil { // Left
        withDoors["left"] = true
    }
    if x < gridSize-1 && roomGrid[x+1][y] != nil { // Right
        withDoors["right"] = true
    }
}

func hasDoor(room *Room, side string) bool {
    for _, obstacle := range room.obstacles {
        if obstacle.IsDoor {
            rect := obstacle.Rect

            // Check if the door is on the top side
            if side == "top" && rect.Min.Y == 0 {
                return true
            }

            // Check if the door is on the bottom side
            if side == "bottom" && rect.Max.Y == screenHeight {
                return true
            }

            // Check if the door is on the left side
            if side == "left" && rect.Min.X == 0 {
                return true
            }

            // Check if the door is on the right side
            if side == "right" && rect.Max.X == screenWidth {
                return true
            }
        }
    }
    return false
}
func generateRoom(x, y int, fromRoom *Room, direction string, roomType RoomType) {
    if x < 0 || x >= gridSize || y < 0 || y >= gridSize || roomGrid[x][y] != nil {
        return // Bounds check and room existence check
    }

    withDoors := map[string]bool{"top": false, "bottom": false, "left": false, "right": false}

    // Set the door for the direction we are coming from
    if direction != "initial" {
        withDoors[oppositeDirection(direction)] = true
    }

    // Check adjacent spaces for potential doors
    if y > 0 && roomGrid[x][y-1] != nil { // Check above
        withDoors["top"] = true
    }
    if y < gridSize-1 && roomGrid[x][y+1] != nil { // Check below
        withDoors["bottom"] = true
    }
    if x > 0 && roomGrid[x-1][y] != nil { // Check left
        withDoors["left"] = true
    }
    if x < gridSize-1 && roomGrid[x+1][y] != nil { // Check right
        withDoors["right"] = true
    }

    // Add doors for adjacent ungenerated rooms
    if y > 0 && roomGrid[x][y-1] == nil {
        withDoors["top"] = true
    }
    if y < gridSize-1 && roomGrid[x][y+1] == nil {
        withDoors["bottom"] = true
    }
    if x > 0 && roomGrid[x-1][y] == nil {
        withDoors["left"] = true
    }
    if x < gridSize-1 && roomGrid[x+1][y] == nil {
        withDoors["right"] = true
    }

	roomGrid[x][y] = &Room{
        obstacles: generateObstacles(withDoors),
        roomType:  roomType, // Set the room type
    }
}

func any(m map[string]bool) bool {
	for _, v := range m {
		if v {
			return true
		}
	}
	return false
}

func ensureDoorAlignment(x, y int, direction string) {
    room := roomGrid[x][y]
    oppositeDir := oppositeDirection(direction)

    // Check for an existing door in the opposite direction
    hasOppositeDoor := false
    for _, obstacle := range room.obstacles {
        if obstacle.IsDoor {
            switch oppositeDir {
            case "top":
                if obstacle.Rect.Min.Y == 0 {
                    hasOppositeDoor = true
                }
            case "bottom":
                if obstacle.Rect.Max.Y == screenHeight {
                    hasOppositeDoor = true
                }
            case "left":
                if obstacle.Rect.Min.X == 0 {
                    hasOppositeDoor = true
                }
            case "right":
                if obstacle.Rect.Max.X == screenWidth {
                    hasOppositeDoor = true
                }
            }
        }
    }

    if !hasOppositeDoor {
        // Add a door in the opposite direction
        doorSize := 30
        var doorRect image.Rectangle
        switch oppositeDir {
        case "top":
            doorRect = image.Rect(screenWidth/2-doorSize/2, 0, screenWidth/2+doorSize/2, 30)
        case "bottom":
            doorRect = image.Rect(screenWidth/2-doorSize/2, screenHeight-30, screenWidth/2+doorSize/2, screenHeight)
        case "left":
            doorRect = image.Rect(0, screenHeight/2-doorSize/2, 30, screenHeight/2+doorSize/2)
        case "right":
            doorRect = image.Rect(screenWidth-30, screenHeight/2-doorSize/2, screenWidth, screenHeight/2+doorSize/2)
        }

        room.obstacles = append(room.obstacles, Obstacle{Rect: doorRect, IsDoor: true})
        roomGrid[x][y] = room // Update the room in the grid
    }
}


func getRandomRoomPositions(gridSize, numBossRooms, numItemRooms int) ([][2]int, [][2]int) {
    var allPositions [][2]int
    for x := 0; x < gridSize; x++ {
        for y := 0; y < gridSize; y++ {
            allPositions = append(allPositions, [2]int{x, y})
        }
    }

    rand.Shuffle(len(allPositions), func(i, j int) {
        allPositions[i], allPositions[j] = allPositions[j], allPositions[i]
    })

    return allPositions[:numBossRooms], allPositions[numBossRooms : numBossRooms+numItemRooms]
}

func contains(slice [][2]int, value [2]int) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}


func (g *Game) getNextRoom(direction string) *Room {
    nextX, nextY := g.player.coordinates[0], g.player.coordinates[1]
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

    if nextX < 0 || nextX >= gridSize || nextY < 0 || nextY >= gridSize {
        return g.currentRoom // Prevent leaving the grid
    }

    if roomGrid[nextX][nextY] == nil {
        generateRoom(nextX, nextY, g.currentRoom, direction, RegularRoom) // Generate with door on correct side
    }

    g.player.coordinates[0], g.player.coordinates[1] = nextX, nextY
    return roomGrid[nextX][nextY]
}


func determineRoomKey(coordinates [2]int) string {
    return fmt.Sprintf("room_%d_%d", coordinates[0], coordinates[1])
}

func createCircleImage(radius int, clr color.Color) *ebiten.Image {
    // Create an image with enough size to hold the circle
    size := radius * 2
    img := ebiten.NewImage(size, size)

    // Draw a circle onto the image
    for y := -radius; y < radius; y++ {
        for x := -radius; x < radius; x++ {
            if x*x+y*y <= radius*radius {
                img.Set(x+radius, y+radius, clr)
            }
        }
    }
    return img
}


func (g *Game) Draw(screen *ebiten.Image) {
	centerX := float32(screenWidth) / 2
    centerY := float32(screenHeight) / 2
    const objectSize = 30 // Size of the object (item or boss)


	switch g.currentRoom.roomType {
    case BossRoom:
        // Draw a red circle for Boss Room
        bossCircle := createCircleImage(objectSize/2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(float64(centerX)-float64(objectSize)/2, float64(centerY)-float64(objectSize)/2)
        screen.DrawImage(bossCircle, opts)

    case ItemRoom:
        // Draw a yellow circle for Item Room
        itemCircle := createCircleImage(objectSize/2, color.RGBA{R: 255, G: 255, B: 0, A: 255})
        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(float64(centerX)-float64(objectSize)/2, float64(centerY)-float64(objectSize)/2)
        screen.DrawImage(itemCircle, opts)
    }

    // Draw the player
	vector.DrawFilledRect(screen, float32(g.player.x), float32(g.player.y), playerWidth, playerHeight, color.White, false)

	    // Draw obstacles
		for _, obstacle := range g.currentRoom.obstacles {
			vector.DrawFilledRect(screen, float32(obstacle.Rect.Min.X), float32(obstacle.Rect.Min.Y), float32(obstacle.Rect.Dx()), float32(obstacle.Rect.Dy()), color.Gray{Y: 0x80}, false)
		}

		for x := 0; x < gridSize; x++ {
			for y := 0; y < gridSize; y++ {
				var roomColor color.Color
	
				if g.roomsVisited[x][y] {
					roomColor = color.Gray{Y: 150} // Visited room color
				} else {
					roomColor = color.Gray{Y: 50} // Unvisited room color
				}
	
				if x == g.player.coordinates[0] && y == g.player.coordinates[1] {
					roomColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Current room color
				}
	
				minimapX := screenWidth - gridSize*minimapRoomSize - minimapPadding + x*minimapRoomSize
				minimapY := screenHeight - gridSize*minimapRoomSize - minimapPadding + y*minimapRoomSize
	
				ebitenutil.DrawRect(screen, float64(minimapX), float64(minimapY), float64(minimapRoomSize), float64(minimapRoomSize), roomColor)
			}
		}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Layout logic...
	return outsideWidth, outsideHeight
}








func main() {
    game := &Game{
        rooms: make(map[string]*Room),
        player: Player{
            x: float64(screenWidth) / 2 - playerWidth / 2,
            y: float64(screenHeight) / 2 - playerHeight / 2,
            inventory: make(map[string]int),
            coordinates: [2]int{gridSize / 2, gridSize / 2}, // Start in the middle of the grid
        },
		roomsVisited: [gridSize][gridSize]bool{},
    }

    generateRooms() // Pre-generate all rooms

    game.currentRoom = roomGrid[gridSize/2][gridSize/2] // Set initial room

    ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
    ebiten.SetWindowSize(screenWidth, screenHeight)

    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}