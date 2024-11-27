package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
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

	treePic, err := loadPicture("../images/tree.png")
	if err != nil {
		panic(err)
	}

	treeSprite := pixel.NewSprite(treePic, treePic.Bounds())

	bushPic, err := loadPicture("../images/bush.png")
	if err != nil {
		panic(err)
	}

	bushSprite := pixel.NewSprite(bushPic, bushPic.Bounds())

	lampPic, err := loadPicture("../images/lamp.png")
	if err != nil {
		panic(err)
	}

	lampSprite := pixel.NewSprite(lampPic, lampPic.Bounds())

	for !win.Closed() {
		sceneSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		treeSprite.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Center().X, 350)))
		bushSprite.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Center().X+150, 350)))
		lampSprite.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Center().X-150, 350)))
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
