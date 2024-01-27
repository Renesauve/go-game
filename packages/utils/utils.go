package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

var images = make(map[string]*ebiten.Image)

func LoadImage(filename string) (*ebiten.Image, error) {
	if img, ok := images[filename]; ok {
		return img, nil
	}

	// Load the image file. The path should be adjusted based on your project structure.
	img, _, err := ebitenutil.NewImageFromFile("assets/gfx/" + filename)
	if err != nil {
		return nil, err
	}

	images[filename] = img
	return img, nil
}
