package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	spriteBorder = 1
	spriteHeight = 16
	spriteWidth  = 16
	spriteZIndex = 1
)

var (
	// Spritesheet is exported because it's needed by other systems
	Spritesheet *common.Spritesheet
	cities      = [...][12]int{
		{99, 100, 101,
			454, 269, 455,
			415, 195, 416,
			452, 306, 453,
		},
		{99, 100, 101,
			268, 269, 270,
			268, 269, 270,
			305, 306, 307,
		},
		{75, 76, 77,
			446, 261, 447,
			446, 261, 447,
			444, 298, 445,
		},
		{75, 76, 77,
			407, 187, 408,
			407, 187, 408,
			444, 298, 445,
		},
		{75, 76, 77,
			186, 150, 188,
			186, 150, 188,
			297, 191, 299,
		},
		{83, 84, 85,
			413, 228, 414,
			411, 191, 412,
			448, 302, 449,
		},
		{83, 84, 85,
			227, 228, 229,
			190, 191, 192,
			301, 302, 303,
		},
		{91, 92, 93,
			241, 242, 243,
			278, 279, 280,
			945, 946, 947,
		},
		{91, 92, 93,
			241, 242, 243,
			278, 279, 280,
			945, 803, 947,
		},
		{91, 92, 93,
			238, 239, 240,
			238, 239, 240,
			312, 313, 314,
		},
	}
)

type City struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type CityBuildingSystem struct {
	buildTime, elapsed float32
	built              int
	mouseTracker       MouseTracker
	usedTiles          []int
	world              *ecs.World
}

// Remove is called whenever an Entity is removed from the World, in order to
// remove it from this system as well
func (*CityBuildingSystem) Remove(ecs.BasicEntity) {}

// Update is ran every frame, with `dt` being the time in seconds since the
// last frame
func (cb *CityBuildingSystem) Update(dt float32) {
	cb.elapsed += dt
	if cb.elapsed >= cb.buildTime {
		cb.generateCity()
		cb.elapsed = 0
		cb.updateBuildTime()
		cb.built++
	}
}

// New is the initialisation of the System
func (cb *CityBuildingSystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	cb.world = w
	fmt.Println("CityBuildingSystem was added to the Scene")
	cb.mouseTracker.BasicEntity = ecs.NewBasic()
	cb.mouseTracker.MouseComponent = common.MouseComponent{Track: true}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(
				&cb.mouseTracker.BasicEntity,
				&cb.mouseTracker.MouseComponent,
				nil,
				nil,
			)
		}
	}
	Spritesheet = common.NewSpritesheetWithBorderFromFile(
		"sprites/CitySheet.png",
		spriteWidth,
		spriteHeight,
		spriteBorder,
		spriteBorder,
	)
	cb.updateBuildTime()
}

// generateCity randomly generates a city in a random location on the map
func (cb *CityBuildingSystem) generateCity() {
	x := rand.Intn(18)
	y := rand.Intn(18)
	t := x + y*18

	for cb.isTileUsed(t) {
		if len(cb.usedTiles) > 300 {
			break //to avoid infinite loop
		}
		x = rand.Intn(18)
		y = rand.Intn(18)
		t = x + y*18
	}
	cb.usedTiles = append(cb.usedTiles, t)
	city := rand.Intn(len(cities))
	cityTiles := make([]*City, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			tile := &City{BasicEntity: ecs.NewBasic()}
			tile.SpaceComponent.Position = engo.Point{
				X: float32(((x+1)*64)+8) + float32(i*16),
				Y: float32(((y + 1) * 64)) + float32(j*16),
			}
			tile.RenderComponent.Drawable = Spritesheet.Cell(cities[city][i+3*j])
			tile.RenderComponent.SetZIndex(spriteZIndex)
			cityTiles = append(cityTiles, tile)
		}
	}
	for _, system := range cb.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, tile := range cityTiles {
				sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
			}
		}
	}
	fmt.Println("Generated new city")
}

func (cb *CityBuildingSystem) isTileUsed(tile int) bool {
	for _, t := range cb.usedTiles {
		if tile == t {
			return true
		}
	}
	return false
}

func (cb *CityBuildingSystem) updateBuildTime() {
	switch {
	case cb.built < 2:
		// 10 to 15 seconds
		cb.buildTime = 5*rand.Float32() + 10
	case cb.built < 8:
		// 60 to 90 seconds
		cb.buildTime = 30*rand.Float32() + 60
	case cb.built < 18:
		// 30 to 90 seconds
		cb.buildTime = 60*rand.Float32() + 30
	case cb.built < 28:
		// 30 to 65 seconds
		cb.buildTime = 35*rand.Float32() + 30
	case cb.built < 33:
		// 30 to 60 seconds
		cb.buildTime = 30*rand.Float32() + 30
	default:
		// 20 to 40 seconds
		cb.buildTime = 20*rand.Float32() + 20
	}
}
