package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/kbinani/screenshot"
	"golang.org/x/exp/rand"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// Common variables

var (
	cupSpriteSheet            pixel.Picture
	cupSprites                []pixel.Rect
	snakeSpriteSheet          pixel.Picture
	snakeSprites              []pixel.Rect
	crabSpriteSheet           pixel.Picture
	crabSprites               []pixel.Rect
	elements                  []*pixel.Sprite
	currentX                  []float64
	matrices                  []pixel.Matrix
	elementsToRemove          []int
	backSpeedFactor           float64 = 50.0
	npcs                              = []Npc{}
	player                    *Player
	lastTimeNpcAdded          = time.Now()
	minNpcLaunchTime          = 5 // seconds
	snakeSpeed                = 100.0
	crabSpeed                 = 120.0
	cupSpeed                  = 80.0
	crabJumpMaxHeight         = 250.0
	crabHorizontalWay         = -1.0
	cupHorizontalWay          = -1.0
	snakeHorizontalWay        = -1.0
	playerJumpLimit           = 500.0
	lastScenarioIndex         = 0
	secondsLastScenarioLaunch = 6.0
	npcStore                  = []func() Npc{}
	secondsText               = "Running for: %.2f seconds"
	secondsRunning            = 0.0
	last                      = time.Now()
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
	collide(pixel.Rect) bool
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

func (c CommonNpcProperties) collide(rect pixel.Rect) bool {
	lowerLeft := pixel.V(c.position.X-c.width/2, c.position.Y-c.height/2)
	upperRight := pixel.V(c.position.X+c.width/2, c.position.Y+c.height/2)
	elementRect := pixel.R(lowerLeft.X, lowerLeft.Y, upperRight.X, upperRight.Y)
	collision := elementRect.Intersect(rect)
	return collision.Area() > 0
}

type Crab struct {
	CommonNpcProperties
	currentJumpHeight        float64
	isJumping                bool
	isFalling                bool
	jumpLimit                float64
	originalVerticalPosition float64
}

func (c *Crab) move(dt float64) bool {
	// Crabs can jump and fall
	if c.isJumping {
		if c.isFalling {
			c.currentJumpHeight = c.currentJumpHeight - (c.speed * 2 * dt)
			if c.currentJumpHeight <= 0 {
				c.currentJumpHeight = 0
				c.isFalling = false
				c.isJumping = false
			}
		} else {
			c.currentJumpHeight = c.currentJumpHeight + (c.speed * 2 * dt)
			if c.currentJumpHeight >= c.jumpLimit {
				c.currentJumpHeight = c.currentJumpHeight - (c.speed * 2 * dt)
				c.isFalling = true
			}
		}
	}
	c.position.X = c.position.X + c.horizontalWay*(c.speed*dt)
	c.position.Y = c.originalVerticalPosition + c.currentJumpHeight
	if !c.isJumping {
		c.secondsToFlip = c.secondsToFlip + dt
		// Let's add a chance to jump
		if rand.Intn(90) == 1 {
			c.isJumping = true
		}
	}
	if c.secondsToFlip > 0.5 {
		c.inverted = !c.inverted
		c.secondsToFlip = 0
	}
	return c.position.X < (-c.width / 2)
}

type Snake struct {
	CommonNpcProperties
}

type Cup struct {
	CommonNpcProperties
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

func (c Player) rect() pixel.Rect {
	lowerLeft := pixel.V(c.position.X-c.width/2, c.position.Y-c.height/2)
	upperRight := pixel.V(c.position.X+c.width/2, c.position.Y+c.height/2)
	return pixel.R(lowerLeft.X, lowerLeft.Y, upperRight.X, upperRight.Y)
}

func (c *Player) move(dt float64) bool {
	// Player is not actually moving it's a fake movement. But it can be jumping or lowering speed
	if c.isJumping {
		if c.isFalling {
			c.currentJumpHeight = c.currentJumpHeight - c.speed*0.8*dt
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

// NPC constructor functions

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

func NewCrab() *Crab {
	if crabSpriteSheet == nil {
		err := error(nil)
		crabSpriteSheet, err = loadPicture("../images/crabSpriteSheet.png")
		if err != nil {
			panic(err)
		}
	}
	if len(crabSprites) == 0 {
		for x := crabSpriteSheet.Bounds().Min.X; x < crabSpriteSheet.Bounds().Max.X; x += 59 {
			for y := crabSpriteSheet.Bounds().Min.Y; y < crabSpriteSheet.Bounds().Max.Y; y += 41 {
				crabSprites = append(crabSprites, pixel.R(x, y, x+59, y+41))
			}
		}
	}
	return &Crab{
		CommonNpcProperties{
			sprite1:       pixel.NewSprite(crabSpriteSheet, crabSprites[0]),
			sprite2:       pixel.NewSprite(crabSpriteSheet, crabSprites[1]),
			position:      pixel.V(1024, 200+41/2),
			height:        41,
			width:         59,
			secondsToFlip: 0,
			speed:         crabSpeed,
			horizontalWay: crabHorizontalWay,
			inverted:      false,
		},
		0.0,
		false,
		false,
		crabJumpMaxHeight,
		200.0 + 41.0/2.0,
	}
}

func getCupYPosition() float64 {
	upperLimit := 768.0 - (60.0 / 2.0) - 50.0
	lowerLimit := 200.0 + (60 / 2.0) + 50.0
	cupPos := float64(rand.Intn(int(upperLimit-lowerLimit))) + 50.0 + lowerLimit
	return cupPos
}

func NewCup() *Cup {
	if cupSpriteSheet == nil {
		err := error(nil)
		cupSpriteSheet, err = loadPicture("../images/cupSpriteSheet.png")
		if err != nil {
			panic(err)
		}
	}
	if len(cupSprites) == 0 {
		for x := cupSpriteSheet.Bounds().Min.X; x < cupSpriteSheet.Bounds().Max.X; x += 42 {
			for y := cupSpriteSheet.Bounds().Min.Y; y < cupSpriteSheet.Bounds().Max.Y; y += 60 {
				cupSprites = append(cupSprites, pixel.R(x, y, x+42, y+60))
			}
		}
	}
	return &Cup{
		CommonNpcProperties{
			sprite1:       pixel.NewSprite(cupSpriteSheet, cupSprites[0]),
			sprite2:       pixel.NewSprite(cupSpriteSheet, cupSprites[1]),
			position:      pixel.V(1024, getCupYPosition()),
			height:        60,
			width:         42,
			secondsToFlip: 0,
			speed:         cupSpeed / float64(rand.Intn(3)+1),
			horizontalWay: cupHorizontalWay,
			inverted:      false,
		},
	}
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

func initGame() {
	elements = []*pixel.Sprite{}
	currentX = []float64{}
	matrices = []pixel.Matrix{}
	elementsToRemove = []int{}
	npcs = []Npc{}
	lastTimeNpcAdded = time.Now()
	secondsRunning = 0.0
	last = time.Now()
}

func run() {

	// Window width and height
	windowWidth := 1024.0
	windowHeight := 768.0

	// Primary display
	bounds := screenshot.GetDisplayBounds(0)
	screenWidth := float64(bounds.Dx())
	screenHeight := float64(bounds.Dy())

	// Center game window on screen
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

	// Text setup

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(10, 160), basicAtlas)
	basicTxt.LineHeight = basicAtlas.LineHeight() * 1.5
	basicTxt.Color = colornames.Yellow
	fmt.Fprintln(basicTxt, "You lost! Press Y to play again or ESC to exit")
	seconds := text.New(pixel.V(10, 15), basicAtlas)
	seconds.LineHeight = basicAtlas.LineHeight() * 1.5
	seconds.Color = colornames.Black
	fmt.Fprintf(seconds, secondsText, secondsRunning)
	helper := text.New(pixel.V(10, 40), basicAtlas)
	helper.LineHeight = basicAtlas.LineHeight() * 1.5
	helper.Color = colornames.Lightcyan
	fmt.Fprintln(helper, "Press UP to jump, LEFT to lower speed, ESC to exit")

	// Sprites loading

	sceneSprite := pixel.NewSprite(pic, pic.Bounds())

	backSpriteSheet, err := loadPicture("../images/back_spritesheet.png")
	if err != nil {
		panic(err)
	}

	boomPic, err := loadPicture("../images/boom.png")
	if err != nil {
		panic(err)
	}

	boomSprite := pixel.NewSprite(boomPic, boomPic.Bounds())

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

	// Init npc store

	npcStore = []func() Npc{
		func() Npc { return NewSnake() },
		func() Npc { return NewCrab() },
		func() Npc { return NewCup() },
	}

	last = time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		secondsRunning = secondsRunning + 1*dt
		last = time.Now()
		sceneSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		seconds.Clear()
		fmt.Fprintf(seconds, secondsText, secondsRunning)
		helper.Draw(win, pixel.IM.Scaled(helper.Orig, 2))
		seconds.Draw(win, pixel.IM.Scaled(seconds.Orig, 2))
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
			if !player.isJumping { //&& !player.isLoweingSpeed {
				player.isJumping = true
			}
		}

		// If User lowers speed

		if win.JustPressed(pixel.KeyLeft) {
			if !player.isJumping && !player.isLoweingSpeed {
				player.isLoweingSpeed = true
			}
		}

		// Launch Npcs?

		if time.Since(lastTimeNpcAdded).Seconds() >= float64(minNpcLaunchTime) {
			launchNpcFactor := rand.Intn(2)
			if launchNpcFactor == 1 {
				whichNpc := rand.Intn(len(npcStore))
				// Add an Npc
				npcs = append(npcs, npcStore[whichNpc]())
			}
			lastTimeNpcAdded = time.Now()
		}

		// Did we have a collision?
		for _, npc := range npcs {
			if npc.collide(player.rect()) {
				// Game Over
				playAgain := false
				for !win.Closed() {
					boomSprite.Draw(win, pixel.IM.Moved(player.position))
					basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 3))
					if win.JustPressed(pixel.KeyY) {
						playAgain = true
						break
					}
					if win.JustPressed(pixel.KeyEscape) {
						playAgain = false
						break
					}
					win.Update()
				}
				if !playAgain {
					return
				}
				initGame()
			}
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
		}

		for _, i := range npcsToRemove {
			npcs = append(npcs[:i], npcs[i+1:]...)
		}

		if win.JustPressed(pixel.KeyEscape) {
			return
		}

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
