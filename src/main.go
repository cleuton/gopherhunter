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

// Common variables

var (
	snakeSpriteSheet pixel.Picture
	snakeSprites     []pixel.Rect
	elements         []*pixel.Sprite
	currentX         []float64
	matrices         []pixel.Matrix
	elementsToRemove []int
	backSpeedFactor  float64 = 50.0
	npcs                     = []Npc{}
	player           *Player
	lastTimeNpcAdded = time.Now()
	minNpcLaunchTime = 5 // seconds
	//crabSpeed                  = 100.0
	snakeSpeed = 100.0
	//mugSpeed                   = 100.0
	//crabJumpSpeed              = 100.0
	//crabJumpMaxHeight          = 100.0
	//crabHorizontalWay          = -1.0
	//mugHorizontalWay           = 1.0
	snakeHorizontalWay        = -1.0
	playerJumpLimit           = 500.0
	lastScenarioIndex         = 0
	secondsLastScenarioLaunch = 6.0
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

// NPC constructors

func NewSnake() *Snake {
	if snakeSpriteSheet == nil {
		err := error(nil)
		snakeSpriteSheet, err = loadPicture("../images/snakesSpriteSheet.png")
		if err != nil {
			panic(err)
		}
	}
	if len(snakeSprites) == 0 {
		for x := snakeSpriteSheet.Bounds().Min.X; x < snakeSpriteSheet.Bounds().Max.X; x += 128 {
			for y := snakeSpriteSheet.Bounds().Min.Y; y < snakeSpriteSheet.Bounds().Max.Y; y += 31 {
				snakeSprites = append(snakeSprites, pixel.R(x, y, x+128, y+31))
			}
		}
	}
	return &Snake{
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
}

type Mug struct {
	CommonNpcProperties
	flyingHeight float64
}

type Player struct {
	CommonNpcProperties
	jumpLimit                  float64
	isJumping                  bool
	isFalling                  bool
	isLoweingSpeed             bool
	isComingBack               bool
	currentJumpHeight          float64
	currentBackPosition        float64
	originalVerticalPosition   float64
	originalHorizontalPosition float64
}

func (c *Player) move(dt float64) bool {
	// Player is not actually moving it's a fake movement. But it can be jumping or lowering speed
	if c.isJumping {
		if c.isFalling {
			c.currentJumpHeight = c.currentJumpHeight - c.speed*dt
			if c.currentJumpHeight <= 0 {
				c.currentJumpHeight = 0
				c.isFalling = false
				c.isJumping = false
			}
		} else {
			c.currentJumpHeight = c.currentJumpHeight + c.speed*dt
			if c.currentJumpHeight >= c.jumpLimit {
				c.currentJumpHeight = c.currentJumpHeight - c.speed*dt
				c.isFalling = true
			}
		}
	} else if c.isLoweingSpeed {
		if c.isComingBack {
			c.currentBackPosition = c.currentBackPosition - c.speed*dt
			if c.currentBackPosition <= 0 {
				c.currentBackPosition = 0
				c.isLoweingSpeed = false
				c.isComingBack = false
			}
		} else {
			c.currentBackPosition = c.currentBackPosition + c.speed*dt
			limit := c.originalHorizontalPosition - c.width/2
			if c.currentBackPosition >= limit {
				c.currentBackPosition = c.currentBackPosition - c.speed*dt
				c.isComingBack = true
			}
		}
	}
	c.position.X = c.originalHorizontalPosition - c.currentBackPosition
	c.position.Y = c.originalVerticalPosition + c.currentJumpHeight
	if !c.isJumping && !c.isLoweingSpeed {
		c.secondsToFlip = c.secondsToFlip + dt
	}
	if c.secondsToFlip > 0.5 {
		c.inverted = !c.inverted
		c.secondsToFlip = 0
	}
	return false
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

	gopherSpriteSheet, err := loadPicture("../images/gopherSpriteSheet.png")
	if err != nil {
		panic(err)
	}

	var gopherSprites []pixel.Rect
	for x := gopherSpriteSheet.Bounds().Min.X; x < gopherSpriteSheet.Bounds().Max.X; x += 60 {
		for y := gopherSpriteSheet.Bounds().Min.Y; y < gopherSpriteSheet.Bounds().Max.Y; y += 49 {
			gopherSprites = append(gopherSprites, pixel.R(x, y, x+60, y+49))
		}
	}

	// Create player
	player = &Player{
		CommonNpcProperties{
			sprite1:       pixel.NewSprite(gopherSpriteSheet, gopherSprites[0]),
			sprite2:       pixel.NewSprite(gopherSpriteSheet, gopherSprites[1]),
			position:      pixel.V(200, 200+49/2),
			height:        49,
			width:         60,
			secondsToFlip: 0,
			speed:         500,
			horizontalWay: 0,
			inverted:      false,
		},
		playerJumpLimit,
		false,
		false,
		false,
		false,
		0.0,
		0.0,
		200.0 + 49/2,
		200.0,
	}

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		sceneSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if secondsLastScenarioLaunch > 5 {
			secondsLastScenarioLaunch = 0
			x := 0
			for {
				x = rand.Intn(len(backSprites))
				if x != lastScenarioIndex {
					break
				}
			}
			lastScenarioIndex = x
			element := pixel.NewSprite(backSpriteSheet, backSprites[x])
			elements = append(elements, element)
			currentX = append(currentX, 1174)
			matrices = append(matrices, pixel.IM.Moved(pixel.V(1174, 350)))
		}
		secondsLastScenarioLaunch += 1 * dt

		// Back scenario

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

		// Player

		player.move(dt)
		player.draw(win)

		// If User jumps

		if win.JustPressed(pixel.KeyUp) {
			// If he already pressed the key, do nothing
			if !player.isJumping && !player.isLoweingSpeed {
				player.isJumping = true
			}
		}

		// If User lowers speed

		if win.JustPressed(pixel.KeyLeft) {
			if !player.isJumping && !player.isLoweingSpeed {
				player.isLoweingSpeed = true
			}
		}

		if time.Since(lastTimeNpcAdded).Seconds() >= float64(minNpcLaunchTime) {
			launchNpcFactor := rand.Intn(2)
			doubleNpcLaunchFactor := rand.Intn(2)
			fmt.Printf("Launch factor: %d\n", launchNpcFactor)
			if launchNpcFactor == 1 {
				fmt.Println("Adding npc:")
				// Add an Npc
				snake := NewSnake()
				npcs = append(npcs, snake)
				if doubleNpcLaunchFactor == 1 {
					// add another Npc close to the first one
					fmt.Println("Adding another npc:")
					snake2 := NewSnake()
					snake2.position.X = snake.position.X + 600
					npcs = append(npcs, snake2)
				}
			}
			lastTimeNpcAdded = time.Now()
		}

		npcsToRemove := []int{}
		for i, npc := range npcs {
			if npc.move(dt) {
				npcsToRemove = append(npcsToRemove, i)
			}
			npc.draw(win)
		}

		// Remove things out of scene:

		for _, i := range elementsToRemove {
			elements = append(elements[:i], elements[i+1:]...)
			matrices = append(matrices[:i], matrices[i+1:]...)
			currentX = append(currentX[:i], currentX[i+1:]...)
			fmt.Println("Removing element")
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
