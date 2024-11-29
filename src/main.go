package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/kbinani/screenshot"
	"golang.org/x/exp/rand"
)

type CommonNpcProperties struct {
	sprite1       *pixel.Sprite
	sprite2       *pixel.Sprite
	position      pixel.Vec
	height        float64
	width         float64
	secondsToFlip float64
	speed         float64
	horizontalWay float64
	inverted      bool
}

type Npc interface {
	move(dt float64) bool
	draw(pixel.Target)
	colide(pixel.Rect) bool
}

func (c *CommonNpcProperties) move(dt float64) bool {
	c.position.X = c.position.X + c.horizontalWay*(c.speed*dt)
	c.secondsToFlip = c.secondsToFlip + dt
	if c.secondsToFlip > 0.5 {
		c.inverted = !c.inverted
		c.secondsToFlip = 0
	}
	return c.position.X < (-c.width / 2)
}
func (c *CommonNpcProperties) draw(target pixel.Target) {
	if c.inverted {
		c.sprite2.Draw(target, pixel.IM.Moved(c.position))
	} else {
		c.sprite1.Draw(target, pixel.IM.Moved(c.position))
	}
}

func (c CommonNpcProperties) colide(rect pixel.Rect) bool {
	return rect.Contains(c.position)
}

type Crab struct {
	CommonNpcProperties
	currentJumpHeight float64
}

type Snake struct {
	CommonNpcProperties
}

type Mug struct {
	CommonNpcProperties
	flyingHeight float64
}

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
		backSpeedFactor  float64 = 50.0
		npcs                     = []Npc{}
		//lastTimeNpcAdded           = time.Now()
		//minNpcLaunchTime           = 5 // seconds
		//crabSpeed                  = 100.0
		snakeSpeed = 100.0
		//mugSpeed                   = 100.0
		//crabJumpSpeed              = 100.0
		//crabJumpMaxHeight          = 100.0
		//crabHorizontalWay          = -1.0
		//mugHorizontalWay           = 1.0
		snakeHorizontalWay = -1.0
	)
	// Window width and height
	windowWidth := 1024.0
	windowHeight := 768.0

	// Primary display
	bounds := screenshot.GetDisplayBounds(0)
	screenWidth := float64(bounds.Dx())
	screenHeight := float64(bounds.Dy())

	// Calcula a posição para centralizar a janela
	posX := (float64(screenWidth) - windowWidth) / 2
	posY := (float64(screenHeight) - windowHeight) / 2
	cfg := opengl.WindowConfig{
		Title:    "Gopher Hunter",
		Bounds:   pixel.R(0, 0, 1024, 768),
		Position: pixel.V(posX, posY),
		VSync:    true,
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

	snakeSpriteSheet, err := loadPicture("../images/snakesSpriteSheet.png")
	if err != nil {
		panic(err)
	}

	var snakeSprites []pixel.Rect
	for x := snakeSpriteSheet.Bounds().Min.X; x < snakeSpriteSheet.Bounds().Max.X; x += 128 {
		for y := snakeSpriteSheet.Bounds().Min.Y; y < snakeSpriteSheet.Bounds().Max.Y; y += 31 {
			snakeSprites = append(snakeSprites, pixel.R(x, y, x+128, y+31))
		}
	}

	last := time.Now()
	runOneTime := false

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

		//if time.Since(lastTimeNpcAdded).Seconds() > minNpcLaunchTime {

		if !runOneTime {
			runOneTime = true
			fmt.Println("Adding npc")
			// Add an Npc
			snake := &Snake{
				CommonNpcProperties{
					sprite1:       pixel.NewSprite(snakeSpriteSheet, snakeSprites[0]),
					sprite2:       pixel.NewSprite(snakeSpriteSheet, snakeSprites[1]),
					position:      pixel.V(1024, 200+31/2),
					height:        31,
					width:         120,
					secondsToFlip: 0,
					speed:         snakeSpeed,
					horizontalWay: snakeHorizontalWay,
					inverted:      false,
				},
			}
			npcs = append(npcs, snake)
		}

		npcsToRemove := []int{}
		for i, npc := range npcs {
			if npc.move(dt) {
				npcsToRemove = append(npcsToRemove, i)
			}
			npc.draw(win)
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

		for _, i := range npcsToRemove {
			fmt.Println("Removing npc")
			npcs = append(npcs[:i], npcs[i+1:]...)
		}

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
