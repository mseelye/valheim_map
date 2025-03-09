# ValheimTileToImage v1.0.0
Mark Seelye - mseelye@yahoo.com  
This is a tool for converting exported Valheim map data to png maps locally.
Map data can be exported from the (great) online Valheim map tool at https://valheim-map.world/

## Data

You either specify the path to the /data directory for the exported data, or run the tool from the data directory.  
The config.json file has a somewhat current set of prefab names, icon names, and colors.

## Icons

All icons come from (the great) site: https://game-icons.net.  
To get a list of what is embedded use -icons  
To set a path to get local icons from use -iconsPath=path/to/icons
Icons are referenced in the confg.json file with file name with extension only, no path.

## Usage
```text
ValheimTileToImage [-flags] [path/to/exported/data]
Usage of ValheimTileToImage.exe:
  -configName string  Name of config file to use (omit json extension!) (default "config")
  -defaultConfig      Outputs the embedded default config file to default-config.json.
  -forcereload        Force the reload and generation of the tile bitmaps.
  -help               Show detailed help information.
  -icons              Show list of icons included with this project.
  -iconsPath string   Path to a directory containing png icons. (default "./icons")
  -prefabCheck        Outputs any PrefabName values NOT represented in the current config.
  -skipMapGeneration  If set existing map.png will be used, tile pngs will not be reprocessed.
  -verbose            Log more debug information.
```

### Example

Download your map data from  https://valheim-map.world/  
Use the "All Data" and "Ultra" quality options.  
(I have not tested with any other exported data yet.)

This will download a file called: MapData_[seedname].zip. Like, "MapData_myawesomeseed.zip"  

Extract this directory and make note of the location.

Run the tool: (pointing at the `\data` subdirectory!)
ValheimTileToImage.exe c:\Users\You\Downloads\MapData_myawesomeseed\data

## Config

You can output a default config file and modify it to configure whatshould be marked on the map.

TODO

## Compiling

TODO

## Thanks

Big thanks to https://valheim-map.world/ and https://game-icons.net/

