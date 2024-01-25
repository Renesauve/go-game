package util

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func CreateCircleImage(radius int, clr color.Color) *ebiten.Image {
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



func OppositeDirection(direction string) string {
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
