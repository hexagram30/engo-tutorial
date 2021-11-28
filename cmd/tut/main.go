package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/hexagram30/engo-tutorial/pkg/systems"
	"github.com/hexagram30/engo-tutorial/pkg/tiles"
)

const (
	edgeScrollSpeed = 400
	edgeRadius      = 20 // units = pixels
	gameName        = "myGame"
	gameTitle       = "Hello World"
	hudHeight       = 200
	hudWidth        = 200
	keyboardSpeed   = edgeScrollSpeed
	worldHeight     = 400
	worldWidth      = 400
	zoomSpeed       = -0.125 // negative means "scrolling down = zooming out"
)

type myGame struct{}

// Type uniquely defines your game type
func (*myGame) Type() string { return gameName }

// Preload is called before loading any assets from the disk, to allow you to
// register / queue them
func (*myGame) Preload() {
	engo.Files.Load(
		"textures/city.png",
		"tilemap/TrafficMap.tmx",
	)
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
		keyboardSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis,
	))
	world.AddSystem(&common.EdgeScroller{edgeScrollSpeed, edgeRadius})
	world.AddSystem(&common.MouseZoomer{zoomSpeed})
	world.AddSystem(&systems.CityBuildingSystem{})
	hud := systems.NewHUD(hudWidth, hudHeight)
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
	level, err := tiles.NewLevel("tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, tile := range level {
				sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
			}
		}
	}
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
