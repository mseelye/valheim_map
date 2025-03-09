package main

import (
	"ValheimTileToImage/valheim_map"
	"flag"
	"fmt"
	"os"
)

var versionString = "v1.0.0"

func banner() {
	fmt.Printf("ValheimTileToImage - %s - mseelye@yahoo.com - 2025-03-01\n", versionString)
}

func help() {
	fmt.Println(`
This is a tool for converting exported Valheim map data.
	
Map data can be exported from the (great) online Valheim map tool at https://valheim-map.world/

You either specify the path to the /data directory for the exported data, or run the tool from the data directory.
The config.json file has a somewhat current set of prefab names, icon names, and colors.

Icons:
All icons come from (the great) site: https://game-icons.net.
To get a list of what is embedded use -icons
To set a path to get local icons from use -iconsPath=path/to/icons

ValheimTileToImage [-flags] [path/to/exported/data]`)
	flag.Usage()
}

func handleMainFlags(config *valheim_map.AppConfig) {
	if config.Verbose {
		valheim_map.Verbosity = 1
	}
	if config.HelpFlag {
		help()
		os.Exit(0)
	}
	if config.IconNamesFlag {
		valheim_map.IconNames()
		os.Exit(0)
	}
}

func main() {
	banner()
	config, err := valheim_map.InitConfig()
	handleMainFlags(config)
	if err == nil {
		if config.OutputOrphanPrefabNames {
			valheim_map.OutputOrphanPrefabNames(config)
		} else {
			valheim_map.Process(config)
		}
	} else {
		panic(fmt.Sprintf("Fatal error initializing application: %s\n", err.Error()))
	}
}

// go run . c:\Users\Mark\Downloads\MapData_bork2\data
// bash shell: rm /c/Users/Mark/Downloads/MapData_bork2/data/tiles/*GO.png
