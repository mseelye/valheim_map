package valheim_map

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type MapDescriptor struct {
	tileSideCount int
	tileRowCount  int
	tileSize      int
	iconSize      int
	basePath      string
	mapObject     map[string]interface{}
	iconsPath     string
}

func NewMapDescriptor(basePath string) *MapDescriptor {
	return &MapDescriptor{
		tileSideCount: -1,
		tileRowCount:  -1,
		tileSize:      -1,
		iconSize:      32,
		basePath:      basePath,
		mapObject:     map[string]interface{}{},
		iconsPath:     "",
	}
}

func (md *MapDescriptor) ReadProperties() {
	md.tileSideCount = int(md.mapObject["TileSideCount"].(float64))
	md.tileRowCount = int(md.mapObject["TileRowCount"].(float64))
	md.tileSize = int(md.mapObject["TileSize"].(float64))
}

func (md *MapDescriptor) Path() string {
	return path.Join(md.basePath, "map.json")
}

func (md *MapDescriptor) MapOutputPath() string {
	return path.Join(md.basePath, "mapGO.png")
}

func (md *MapDescriptor) MapWithIconsOutputPath() string {
	return path.Join(md.basePath, "map-with-iconsGO.png")
}

func (md *MapDescriptor) TileInputPath(x int, z int) string {
	return path.Join(md.basePath, "tiles", fmt.Sprintf("%02d-%02d.bin.gz", x, z))
}

func (md *MapDescriptor) TileOutputPath(x int, z int) string {
	return path.Join(md.basePath, "tiles", fmt.Sprintf("%02d-%02dGO.png", x, z))
}

func (md *MapDescriptor) IconPath(iconPng string) string {
	if md.iconsPath != "" {
		return path.Join(md.iconsPath, iconPng)
	} else {
		return path.Join(md.basePath, "icons", iconPng)
	}
}

func (md *MapDescriptor) ScaleFactor() float64 {
	return float64(md.tileSize) / float64(md.tileRowCount)
}
func (md *MapDescriptor) ImageWidth() float64 {
	return float64(md.tileRowCount * md.tileSideCount)
}

func ReadMapDescriptor(basePath string, iconsPath string) *MapDescriptor {
	mapDescriptor := NewMapDescriptor(basePath)
	mapDescriptor.iconsPath = iconsPath
	content, err := os.ReadFile(mapDescriptor.Path())
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("You must specify the data folder as an argument, or run this from the data folder.")
		return nil
	} else {
		err := json.Unmarshal(content, &mapDescriptor.mapObject)
		if err != nil {
			fmt.Printf("Cannot parse json file, '%s'. Error: %v\n", mapDescriptor.Path(), err)
			return nil
		} else {
			mapDescriptor.ReadProperties()
		}
	}
	return mapDescriptor
}
