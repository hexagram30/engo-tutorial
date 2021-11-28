package tiles

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func NewLevel(tmxAsset string) ([]*Tile, error) {
	resource, err := engo.Files.Resource(tmxAsset)
	if err != nil {
		return nil, err
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level
	tiles := make([]*Tile, 0)
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement.Image,
					Scale:    engo.Point{1, 1},
				}
				// tile.RenderComponent.SetZIndex(1)
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
				}
				tiles = append(tiles, tile)
			}
		}
	}
	common.CameraBounds = levelData.Bounds()
	return tiles, nil
}
