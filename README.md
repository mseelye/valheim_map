# ValheimTileToImage v1.0.0
Mark Seelye - mseelye@yahoo.com  
This is a tool for converting exported Valheim map data to png maps locally.
Map data can be exported from the (great) online Valheim map tool at https://valheim-map.world/

## Data

You either specify the path to the /data directory for the exported data, or run the tool from the data directory.  
The config.json file has a somewhat current set of prefab names, icon names, and colors.

Currently, I do not have map generation like https://valheim-map.world/, but I may add that or create my own tool at some point.

So for now use https://valheim-map.world/ and download your map data.  
Use the *"All Data"* and *"Ultra"* quality options. I have not tested any of this with any other resolution levels.

## Icons

All icons come from (the great) site: https://game-icons.net.  
To get a list of what is embedded use `-icons`  
To set a path to get local icons from use `-iconsPath=path/to/icons`  
Icons are referenced in the confg.json file with *file name with extension*, no path. Ex: `ogre.png`

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
With the compiled executable:  
`ValheimTileToImage.exe c:\Users\You\Downloads\MapData_myawesomeseed\data`

With go:
`go run . c:\Users\You\Downloads\MapData_myawesomeseed\data`

## Config

You can output a default config file and modify it to configure whatshould be marked on the map.

The config file currently only lets you define colors, and an array of prefabs to plot.

### Colors

Color set in the colors section can be used from the various "plot" sections.  
You can use any name, and the color values are simple #RRGGBBAA values.  
You can also have a color use the name of another color.  
```json
{
  "colors": {
    "White": "#ffffffff",
    "Red": "#ff0000ff",
    "Yellow": "#ffff00ff",
    "Green": "#008000ff",
    "colorBosses": "Red",
...
```

Be careful as I did not put in any protection against infinite loops.  
Do *NOT* do this: 
```json
"colorQuests": "Yellow",
"Yellow": "colorQuests",
```  
You can have multple things reference the same thing though:  
```json
"colorQuests": "Yellow",
"colorBosses": "Yellow",
"colorSpawner": "Yellow",
```
There is nothing special about any of the names, however, there are some embedded colors that the program will try to use if it doesn't find the name in your config.  
```go
var (
	NilColor    = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgb(0, 0, 0, 0)
	Black       = color.RGBA{0x00, 0x00, 0x00, 0xff} // rgb(0, 0, 0)
	Green       = color.RGBA{0x00, 0x80, 0x00, 0xff} // rgb(0, 128, 0)
	Seagreen    = color.RGBA{0x2e, 0x8b, 0x57, 0xff} // rgb(46, 139, 87)
	Gray        = color.RGBA{0x33, 0x33, 0x33, 0xff} // rgb(51, 51, 51)
	Lightgray   = color.RGBA{0xd3, 0xd3, 0xd3, 0xff} // rgb(211, 211, 211)
	Darkgreen   = color.RGBA{0x00, 0x64, 0x00, 0xff} // rgb(0, 100, 0)
	Yellow      = color.RGBA{0xff, 0xff, 0x00, 0xff} // rgb(255, 255, 0)
	Red         = color.RGBA{0xff, 0x00, 0x00, 0xff} // rgb(255, 0, 0)
	White       = color.RGBA{0xff, 0xff, 0xff, 0xff} // rgb(255, 255, 255)
	Cyan        = color.RGBA{0x00, 0xff, 0xff, 0xff} // rgb(0, 255, 255)
	Darkblue    = color.RGBA{0x00, 0x00, 0x8b, 0xff} // rgb(0, 0, 139)
	Purple      = color.RGBA{0x80, 0x00, 0x80, 0xff} // rgb(128, 0, 128)
	Blue        = color.RGBA{0x00, 0x00, 0xff, 0xff} // rgb(0, 0, 255)
	Lime        = color.RGBA{0x00, 0xff, 0x00, 0xff} // rgb(0, 255, 0)
	Darkgrey    = color.RGBA{0xa9, 0xa9, 0xa9, 0xff} // rgb(169, 169, 169)
	Gold        = color.RGBA{0xff, 0xd7, 0x00, 0xff} // rgb(255, 215, 0)
	Darkbrown   = color.RGBA{0x2b, 0x1d, 0x0e, 0xff} // rgb(43, 29, 14)
	Brown       = color.RGBA{0xa5, 0x2a, 0x2a, 0xff} // rgb(165, 42, 42)
	Saddlebrown = color.RGBA{0x8b, 0x45, 0x13, 0xff} // rgb(139, 69, 19)
	Lightblue   = color.RGBA{0xad, 0xd8, 0xe6, 0xff} // rgb(173, 216, 230)
)
```
If color is in your config it will use that first.

## Plots

Plots are an array of the prefab-locations you want shown on your map with icons.  

```json
"plots": [
    {
      "PrefabName": "StartTemple",
      "IconPng": "impact-point.png",
      "IconColor": "White",
      "HighlightColor": "Gold",
      "Text": "Spawn"
    },
...
```  
`PrefabName` is string to look for for a location. It will attempt to match anything that starts with the text given.  
Meaning if you have "Swamp" it would match all the locations with a prefabname starting with "Swamp": `"SwampHut1", "SwampHut2", "SwampHut3", "SwampHut4", "SwampHut5", "SwampRuin1", "SwampRuin2", "SwampWell1",`  

`IconPng` is a string with a filename of a png to use for the icon. The icon png will be resized and tinted (the white portion of an icon) to the `IconColor` value.  
The filename is whatever the filename is in your `/icons` directory, you do not specify a path here.  
The system has some built-in icons that if your icon directory does not have it will attempt to use the embedded one.  

`IconColor` is a NAME of a color defined in the colors section, do not put #RRGGBBAA values here.  The color will be applied to the icon's "white" color.  

`HighlightColor` if provided the program will draw a shadowed ellipse under the icon and color it to this color.  

`Text` if provided the program will plot this text under/near the plotted icon.  No support for changing color, size, or font for this yet.

You can have as many or as few plots entries as you like. Currently there is no support for filtering, but you can always have different config files for different set of things you need plotted separately.

## Compiling

Pull the source from here.  
Do any usual golang dance you need to install modules (`go mod tidy` etc.) and then
```bash
go build .
```

## Thanks

Big thanks to https://valheim-map.world/ for providing the excellent export feature and code samples.  
Huge thanks also to https://game-icons.net/ for providing awesome free, high-quality icons!



