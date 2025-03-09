package valheim_map

import "image/color"

type Biome int

const (
	None        Biome = 0
	Meadows     Biome = 1
	Swamp       Biome = 2
	Mountain    Biome = 4
	BlackForest Biome = 8
	Plains      Biome = 16
	AshLands    Biome = 32
	DeepNorth   Biome = 64
	Ocean       Biome = 256
	Mistlands   Biome = 512
)

var biomeName = map[Biome]string{
	None:        "none",
	Meadows:     "Meadows",
	Swamp:       "Swamp",
	Mountain:    "Mountain",
	BlackForest: "Black Forest",
	Plains:      "Plains",
	AshLands:    "AshLands",
	DeepNorth:   "DeepNorth",
	Ocean:       "Ocean",
	Mistlands:   "Mistlands",
}

var biomeColor = map[Biome]color.RGBA{
	None:        Black,
	Meadows:     Green,
	Swamp:       Darkbrown, //Saddlebrown, // Brown, // Seagreen,
	Mountain:    Lightgray,
	BlackForest: Darkgreen,
	Plains:      Yellow,
	AshLands:    Red,
	DeepNorth:   White,
	Ocean:       Darkblue,
	Mistlands:   Purple,
}

func (biome Biome) String() string {
	return biomeName[biome]
}

func (biome Biome) Color() color.RGBA {
	return biomeColor[biome]
}
