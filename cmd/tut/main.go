package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/hexagram30/engo-tutorial/pkg/systems"
)

const (
	cameraSpeed = 400
	edgeRadius  = 20 // units = pixels
	gameName    = "myGame"
	gameTitle   = "Hello World"
	worldHeight = 400
	worldWidth  = 400
	zoomSpeed   = -0.125 // negative means "scrolling down = zooming out"
)

type myGame struct{}

// Type uniquely defines your game type
func (*myGame) Type() string { return gameName }

// Preload is called before loading any assets from the disk, to allow you to
// register / queue them
func (*myGame) Preload() {
	engo.Files.Load("textures/city.png")
}

// Setup is called before the main loop starts. It allows you to add entities
// and systems to your Scene.
func (*myGame) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	world.AddSystem(common.NewKeyboardScroller(
		cameraSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis,
	))
	world.AddSystem(&common.EdgeScroller{cameraSpeed, edgeRadius})
	world.AddSystem(&common.MouseZoomer{zoomSpeed})
	world.AddSystem(&systems.CityBuildingSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          gameTitle,
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &myGame{})
}
