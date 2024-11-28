package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/exp/rand"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {

	var (
		elements         []*pixel.Sprite
		currentX         []float64
		matrices         []pixel.Matrix
		elementsToRemove []int
		backSpeedFactor  float64 = 100
	)

	cfg := opengl.WindowConfig{
		Title:  "Gopher Hunter",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	pic, err := loadPicture("../images/scene.png")
	if err != nil {
		panic(err)
	}

	sceneSprite := pixel.NewSprite(pic, pic.Bounds())

	backSpriteSheet, err := loadPicture("../images/back_spritesheet.png")
	if err != nil {
		panic(err)
	}

	var backSprites []pixel.Rect
	for x := backSpriteSheet.Bounds().Min.X; x < backSpriteSheet.Bounds().Max.X; x += 300 {
		for y := backSpriteSheet.Bounds().Min.Y; y < backSpriteSheet.Bounds().Max.Y; y += 300 {
			backSprites = append(backSprites, pixel.R(x, y, x+300, y+300))
		}
	}

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		sceneSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if win.JustPressed(pixel.MouseButtonLeft) {
			element := pixel.NewSprite(backSpriteSheet, backSprites[rand.Intn(len(backSprites))])
			elements = append(elements, element)
			mouseX := win.MousePosition().X
			currentX = append(currentX, mouseX)
			matrices = append(matrices, pixel.IM.Moved(pixel.V(mouseX, 350)))
		}

		elementsToRemove = []int{}
		for i, element := range elements {
			element.Draw(win, matrices[i])
			currentX[i] = currentX[i] - (backSpeedFactor * dt)
			if currentX[i] < -150 {
				elementsToRemove = append(elementsToRemove, i)
			} else {
				matrices[i] = matrices[i].Moved(pixel.V(-backSpeedFactor*dt, 0))
			}
		}
		for _, i := range elementsToRemove {
			elements = append(elements[:i], elements[i+1:]...)
			matrices = append(matrices[:i], matrices[i+1:]...)
			currentX = append(currentX[:i], currentX[i+1:]...)
		}
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
